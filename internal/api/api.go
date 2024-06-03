package api

import (
  "context"
  "errors"
  "fmt"
  "github.com/benthosdev/benthos-builder/internal/store"
  "github.com/gorilla/mux"
  "github.com/nats-io/nats.go"
  "github.com/nats-io/nats.go/jetstream"
  "github.com/rs/zerolog/log"
  "io"
  "net/http"
  "time"
)

type Api struct {
  port      int
  nc        *nats.Conn
  artifacts jetstream.ObjectStore
}

func NewApi(nc *nats.Conn, js jetstream.JetStream, port int) (*Api, error) {
  artifacts, err := js.ObjectStore(context.Background(), store.JetstreamOSArtifacts)
  if err != nil {
    return nil, fmt.Errorf("failed to create object store: %w", err)
  }

  return &Api{
    port:      port,
    nc:        nc,
    artifacts: artifacts,
  }, nil
}

func (a *Api) Run(ctx context.Context) error {
  address := fmt.Sprintf(":%d", a.port)

  server := &http.Server{
    Addr:         address,
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout:  10 * time.Second,
  }

  router := mux.NewRouter().StrictSlash(true)

  buildRouter := router.PathPrefix("/builds").Subrouter().StrictSlash(true)
  buildRouter.Handle("/", createHandlerFunc(a.nc, "build.request")).Methods(http.MethodPost)
  buildRouter.Handle("/{bid}/artifact", createObjectReader(a.artifacts, "bid")).Methods(http.MethodPost)

  router.Use(func(inner http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      start := time.Now()
      defer func() {
        log.Info().
          Str("method", r.Method).
          Str("path", r.RequestURI).
          Msgf("processing time: %dms", time.Since(start).Milliseconds())
      }()

      inner.ServeHTTP(w, r)
    })
  })

  log.Info().Msgf("api running on port %d", a.port)
  return server.ListenAndServe()
}

func createHandlerFunc(nc *nats.Conn, endpoint string) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    b, err := io.ReadAll(r.Body)
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

func createObjectReader(obj jetstream.ObjectStore, idParam string) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id := params[idParam]

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

    if _, err := io.Copy(w, or); err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      w.Write([]byte(err.Error()))
    }
  }
}
