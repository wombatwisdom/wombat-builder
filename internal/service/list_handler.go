package service

import (
  "context"
  "encoding/json"
  "github.com/benthosdev/benthos-builder/internal/store"
  "github.com/benthosdev/benthos-builder/public/model"
  "github.com/nats-io/nats.go/micro"
  "github.com/rs/zerolog/log"
)

type (
  BuildListRequest struct {
    Query string `json:"query" jsonschema_description:"The query to search for"`
  }
  BuildListResponse struct {
    Builds []model.Build `json:"builds" jsonschema_description:"The list of builds"`
  }
)

func (r *BuildListRequest) Validate() error {
  return nil
}

func getBuildListHandler(s *store.Store) micro.HandlerFunc {
  return func(request micro.Request) {
    var req BuildListRequest
    if err := json.Unmarshal(request.Data(), &req); err != nil {
      _ = request.Error("BAD_REQUEST", "failed to parse request", []byte(err.Error()))
      return
    }

    if err := req.Validate(); err != nil {
      _ = request.Error("BAD_REQUEST", "invalid request", []byte(err.Error()))
      return
    }

    ids, err := s.BuildsIndex.Search(req.Query)
    if err != nil {
      _ = request.Error("BACKBONE_ERROR", "failed to search for builds", []byte(err.Error()))
      return
    }

    var result []model.Build
    for _, id := range ids {
      build, err := s.Builds.Get(context.Background(), id)
      if err != nil {
        log.Warn().Err(err).Msgf("failed to get build %s", id)
        continue
      }
      result = append(result, *build)
    }

    if err := request.RespondJSON(BuildListResponse{Builds: result}); err != nil {
      _ = request.Error("INTERNAL_ERROR", "failed to respond", []byte(err.Error()))
      return
    }
  }
}
