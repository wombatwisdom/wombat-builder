package store

import (
  "context"
  "encoding/json"
  "errors"
  "github.com/benthosdev/benthos-builder/public/model"
  "github.com/nats-io/nats.go/jetstream"
)

type Builds struct {
  kv jetstream.KeyValue
}

func (b *Builds) Watch(ctx context.Context) (jetstream.KeyWatcher, error) {
  return b.kv.Watch(ctx, "build.*")
}

func (b *Builds) Get(ctx context.Context, key string) (*model.Build, error) {
  entry, err := b.kv.Get(ctx, key)
  if err != nil {
    if errors.Is(err, jetstream.ErrKeyNotFound) {
      return nil, nil
    }

    return nil, err
  }

  var build model.Build
  if err := json.Unmarshal(entry.Value(), &build); err != nil {
    return nil, err
  }

  return &build, nil
}

func (b *Builds) Update(ctx context.Context, key string, build *model.Build, revision uint64) (uint64, error) {
  bb, err := json.Marshal(build)
  if err != nil {
    return 0, err
  }

  return b.kv.Update(ctx, key, bb, revision)
}

func (b *Builds) Set(ctx context.Context, build *model.Build) (uint64, error) {
  data, err := json.Marshal(build)
  if err != nil {
    return 0, err
  }

  return b.kv.Put(ctx, build.Id(), data)
}
