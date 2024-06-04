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

func NewStore(js jetstream.JetStream, withIndex bool) (*Store, error) {
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

  var bi *BuildIndex
  if withIndex {
    bi, err = NewBuildIndex(ctx, js)
    if err != nil {
      return nil, err
    }
  }

  return &Store{
    Artifacts:   &Artifacts{obj: artifacts},
    Builds:      &Builds{kv: builds},
    Repos:       &Repos{kv: repos},
    BuildsIndex: bi,
  }, nil
}

type Store struct {
  Artifacts   *Artifacts
  Builds      *Builds
  BuildsIndex *BuildIndex
  Repos       *Repos
}
