package store

import (
  "context"
  "github.com/nats-io/nats.go/jetstream"
)

const (
  JetstreamKVBuilds    = "builds"
  JetstreamKVRepos     = "repos"
  JetstreamOSArtifacts = "artifacts"
)

func NewStore(js jetstream.JetStream) (*Store, error) {
  ctx := context.Background()
  builds, err := js.KeyValue(ctx, JetstreamKVBuilds)
  if err != nil {
    return nil, err
  }

  repos, err := js.KeyValue(ctx, JetstreamKVRepos)
  if err != nil {
    return nil, err
  }

  artifacts, err := js.ObjectStore(ctx, JetstreamOSArtifacts)
  if err != nil {
    return nil, err
  }

  return &Store{
    Artifacts: &Artifacts{obj: artifacts},
    Builds:    &Builds{kv: builds},
    Repos:     &Repos{kv: repos},
  }, nil
}

type Store struct {
  Artifacts *Artifacts
  Builds    *Builds
  Repos     *Repos
}
