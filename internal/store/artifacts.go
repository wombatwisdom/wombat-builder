package store

import (
  "context"
  "github.com/nats-io/nats.go/jetstream"
  "io"
  "os"
)

type Artifacts struct {
  obj jetstream.ObjectStore
}

func (a *Artifacts) WriteFile(ctx context.Context, name string, path string) (*jetstream.ObjectInfo, error) {
  reader, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer reader.Close()

  om := jetstream.ObjectMeta{
    Name: name,
  }

  return a.obj.Put(ctx, om, reader)
}

func (a *Artifacts) Read(ctx context.Context, name string) (io.ReadCloser, error) {
  return a.obj.Get(ctx, name)
}
