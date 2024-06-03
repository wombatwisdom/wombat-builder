package builder

import (
  "context"
  "encoding/json"
  "fmt"
  "github.com/benthosdev/benthos-builder/internal/store"
  "github.com/benthosdev/benthos-builder/public/model"
  "github.com/nats-io/nats.go/jetstream"
  "github.com/rs/xid"
  "github.com/rs/zerolog/log"
)

func NewBuilder(s *store.Store, workers int) (*Builder, error) {
  return &Builder{
    Id:    xid.New().String(),
    s:     s,
    queue: make(chan buildWithRevision, workers),
  }, nil
}

type Builder struct {
  Id string
  s  *store.Store

  queue chan buildWithRevision
}

func (b *Builder) Run(ctx context.Context) error {
  // -- start the workers
  for i := 0; i < cap(b.queue); i++ {
    go b.worker(ctx, i)
  }

  // -- start watching the builds for new builds
  kw, err := b.s.Builds.Watch(ctx)
  if err != nil {
    return fmt.Errorf("failed to watch builds: %w", err)
  }

  log.Info().Msgf("builder started with %d workers", cap(b.queue))

  for {
    select {
    case <-ctx.Done():
      return nil
    case update := <-kw.Updates():
      if update == nil {
        continue
      }

      if update.Operation() == jetstream.KeyValueDelete {
        continue
      }

      // -- get the build from the update
      var build model.Build
      if err := json.Unmarshal(update.Value(), &build); err != nil {
        log.Error().Err(err).Msg("failed to unmarshal build")
      }

      // -- this is where the race starts. We will update the build state and try to write it. If the
      // -- write succeeds, we are the first ones to claim the build and we can start building it
      // -- otherwise, we will ignore the build and let the other builder handle it
      build.Builder = b.Id
      if _, err := b.s.Builds.Update(ctx, update.Key(), &build, update.Revision()); err != nil {
        // TODO: we actually need to check if the error is a conflict error
        continue
      }

      // -- if we made it here, we can start the build
      b.queue <- buildWithRevision{build, update.Revision()}
    }
  }
}

func (b *Builder) worker(ctx context.Context, id int) {
  logger := log.With().Int("worker", id).Logger()

  for {
    select {
    case <-ctx.Done():
      return
    case build := <-b.queue:
      // -- update the build status to pending
      build.Status = model.BuildStatusBuilding

      rev, err := b.s.Builds.Update(ctx, build.Id(), &build.Build, build.revision)
      if err != nil {
        logger.Error().Err(err).Msg("failed to start build")
        continue
      }

      // -- start the build
      task := &BuildTask{Build: &build.Build}
      artifactPath, err := task.Run(ctx)
      if err == nil {
        // -- upload the artifact to the object store
        oi, err := b.s.Artifacts.WriteFile(ctx, build.Id(), artifactPath)
        if err == nil {
          build.Artifact = model.ArtifactReference(oi.Name)
          build.Status = model.BuildStatusSuccess
        } else {
          build.Status = model.BuildStatusFailed
          build.Error = err.Error()
        }
      } else {
        build.Status = model.BuildStatusFailed
        build.Error = err.Error()
      }

      // -- update the build
      _, err = b.s.Builds.Update(ctx, build.Id(), &build.Build, rev)
      if err != nil {
        logger.Error().Err(err).Msg("failed to update build")
      }
    }
  }
}

type buildWithRevision struct {
  model.Build
  revision uint64
}
