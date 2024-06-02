package store

import "github.com/nats-io/nats.go/jetstream"

type Repos struct {
  kv jetstream.KeyValue
}
