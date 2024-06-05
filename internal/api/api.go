package api

import (
  "context"
  "embed"
  "errors"
  "fmt"
  "github.com/benthosdev/benthos-builder/internal/store"
  "github.com/gorilla/mux"
  "github.com/nats-io/nats.go"
  "github.com/nats-io/nats.go/jetstream"
  "github.com/rs/zerolog/log"
  "io"
  "io/fs"
  "net/http"
  "net/url"
  "time"
)

//go:embed web
var web embed.FS

type Api struct {
  port      int
  nc        *nats.Conn
  artifacts jetstream.ObjectStore
  enableUi  bool
}

func NewApi(nc *nats.Conn, js jetstream.JetStream, port int, enableUi bool) (*Api, error) {
  artifacts, err := js.ObjectStore(context.Background(), store.JetstreamOSArtifacts)
  if err != nil {
    return nil, fmt.Errorf("failed to create object store: %w", err)
  }

  return &Api{
    port:      port,
    nc:        nc,
    artifacts: artifacts,
    enableUi:  enableUi,
  }, nil
}

func (a *Api) Run(ctx context.Context) error {
  address := fmt.Sprintf(":%d", a.port)

  server := &http.Server{
    Addr:         address,
    ReadTimeout:  60 * time.Second,
    WriteTimeout: 0,
    IdleTimeout:  10 * time.Second,
  }

  router := mux.NewRouter()

  ar := router.PathPrefix("/api").Subrouter().StrictSlash(true)
  ar.Use(func(inner http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      start := time.Now()
      defer func() {
        log.Info().
          Str("method", r.Method).
          Str("path", r.RequestURI).
          Msgf("processing time: %dms", time.Since(start).Milliseconds())
      }()

      log.Info().
        Str("method", r.Method).
        Str("path", r.RequestURI).
        Msg("requested")

      inner.ServeHTTP(w, r)
    })
  })
  ar.Use(func(inner http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      w.Header().Add("Access-Control-Allow-Origin", "*")
      inner.ServeHTTP(w, r)
    })
  })

  buildRouter := ar.PathPrefix("/builds").Subrouter()
  buildRouter.Handle("", createHandlerFunc(a.nc, "build.request")).Methods(http.MethodPost)
  buildRouter.Handle("", createHandlerFuncWithCallback(a.nc, "build.list", func(r *http.Request) ([]byte, error) {
    q, err := url.QueryUnescape(r.URL.Query().Get("q"))
    if err != nil {
      return nil, err
    }

    return []byte(fmt.Sprintf("{\"query\": \"%s\"}", q)), nil
  })).Methods(http.MethodGet)

  artifactRouter := router.PathPrefix("/artifacts").Subrouter()
  artifactRouter.Handle("/{arch}/{os}/{ver}/{hash}", createObjectReader(a.artifacts, func(r *http.Request) string {
    params := mux.Vars(r)
    return fmt.Sprintf("build.%s.%s.%s.%s", params["arch"], params["os"], params["ver"], params["hash"])
  })).Methods(http.MethodGet)

  if a.enableUi {
    dist, err := fs.Sub(web, "web/dist")
    if err != nil {
      return fmt.Errorf("failed to navigate web fs: %w", err)
    }

    router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.FS(dist))))
    //router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    //  http.ServeFile(w, r, "./build/index.html")
    //})

  }

  log.Info().Msgf("api running on port %d", a.port)
  router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
    url, _ := route.GetPathTemplate()
    met, err := route.GetMethods()
    if err != nil {
      met = []string{}
    }

    log.Info().Msgf("api %s %s", url, met)
    return nil
  })
  server.Handler = router

  return server.ListenAndServe()
}

func createHandlerFunc(nc *nats.Conn, endpoint string) http.HandlerFunc {
  return createHandlerFuncWithCallback(nc, endpoint, func(r *http.Request) ([]byte, error) {
    return io.ReadAll(r.Body)
  })
}

func createHandlerFuncWithCallback(nc *nats.Conn, endpoint string, cb func(r *http.Request) ([]byte, error)) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    b, err := cb(r)
    if err != nil {
      w.WriteHeader(http.StatusBadRequest)
      return
    }

    resp, err := nc.Request(endpoint, b, 10*time.Second)
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      _, _ = w.Write([]byte(err.Error()))
      return
    }
    w.WriteHeader(http.StatusOK)
    for h, v := range resp.Header {
      w.Header().Set(h, v[0])
    }
    _, _ = w.Write(resp.Data)
  }
}

func createObjectReader(obj jetstream.ObjectStore, idCb func(r *http.Request) string) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    id := idCb(r)

    or, err := obj.Get(r.Context(), id)
    if err != nil {
      if errors.Is(err, jetstream.ErrObjectNotFound) {
        w.WriteHeader(http.StatusNotFound)
      } else {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
      }

      return
    }
    defer or.Close()

    oi, err := or.Info()
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      w.Write([]byte(err.Error()))
      return
    }

    w.Header().Set("Content-Length", fmt.Sprintf("%d", oi.Size))
    w.Header().Set("Content-Type", "application/octet-stream")
    w.Header().Set("Content-Disposition", "attachment; filename=\"wombat\"")
    if _, err := io.Copy(w, or); err != nil {
      log.Err(err).Msgf("failed to write object %s", id)
      w.WriteHeader(http.StatusInternalServerError)
      w.Write([]byte(err.Error()))
    }
  }
}
