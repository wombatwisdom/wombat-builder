package cmd

import (
  "github.com/nats-io/nats.go"
  "github.com/nats-io/nats.go/jetstream"
  "github.com/rs/zerolog/log"
  "github.com/urfave/cli/v2"
)

var NatsFlags = []cli.Flag{
  &cli.StringFlag{
    Name:    "nats-url",
    Usage:   "the nats url to connect to",
    Value:   "tls://connect.ngs.global",
    EnvVars: []string{"NATS_URL"},
  },
  &cli.StringFlag{
    Name:    "nats-user-jwt",
    Usage:   "the nats user jwt",
    EnvVars: []string{"NATS_USER_JWT"},
  },
  &cli.StringFlag{
    Name:    "nats-user-seed",
    Usage:   "the nats user seed",
    EnvVars: []string{"NATS_USER_SEED"},
  },
  &cli.StringFlag{
    Name:    "nats-user-creds",
    Usage:   "the nats user creds file",
    EnvVars: []string{"NATS_USER_CREDS"},
  },
}

func ConnectNats(cCtx *cli.Context) (*nats.Conn, jetstream.JetStream, error) {
  natsURL, opts, err := ParseNatsFlags(cCtx)
  if err != nil {
    return nil, nil, err
  }

  nc, err := nats.Connect(natsURL, opts...)
  if err != nil {
    return nil, nil, err
  }

  js, err := jetstream.New(nc)
  if err != nil {
    return nil, nil, err
  }

  return nc, js, nil
}

func ParseNatsFlags(cCtx *cli.Context) (string, []nats.Option, error) {
  natsURL := cCtx.String("nats-url")
  var opts []nats.Option

  natsUserJwt := cCtx.String("nats-user-jwt")
  natsUserSeed := cCtx.String("nats-user-seed")
  natsUserCreds := cCtx.String("nats-user-creds")
  if natsUserJwt != "" && natsUserSeed != "" {
    log.Info().Msgf("using user jwt and seed")
    opts = append(opts, nats.UserJWTAndSeed(natsUserJwt, natsUserSeed))
  } else if natsUserCreds != "" {
    log.Info().Msgf("using user creds file: %s", natsUserCreds)
    opts = append(opts, nats.UserCredentials(natsUserCreds))
  } else {
    log.Info().Msgf("not using auth")
  }

  return natsURL, opts, nil
}
