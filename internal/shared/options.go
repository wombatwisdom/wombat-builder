package shared

import "github.com/nats-io/nats.go"

type Option func(o *Options)

type Options struct {
  NatsUrl     string
  NatsOptions []nats.Option
}

func WithNatsUrl(url string) Option {
  return func(o *Options) {
    o.NatsUrl = url
  }
}

func WithNatsOptions(opts ...nats.Option) Option {
  return func(o *Options) {
    o.NatsOptions = append(o.NatsOptions, opts...)
  }
}
