package cmd

import (
  "fmt"
  "github.com/benthosdev/benthos-builder/internal/cmd"
  "github.com/benthosdev/benthos-builder/internal/service"
  "github.com/benthosdev/benthos-builder/internal/store"
  "github.com/nats-io/nats.go"
  "github.com/rs/zerolog/log"
  "github.com/urfave/cli/v2"
)

var ServiceCommand = &cli.Command{
  Name:  "service",
  Usage: "run the builder service",
  Description: `
The service exposes a nats micro service that can be used to manage the process of building artifacts. 
`,
  Flags: append(cmd.NatsFlags, []cli.Flag{}...),
  Action: func(cCtx *cli.Context) error {
    nc, js, err := cmd.ConnectNats(cCtx)
    if err != nil {
      return fmt.Errorf("failed to connect to nats: %w", err)
    }
    defer nc.Close()

    s, err := store.NewStore(js)
    if err != nil {
      return fmt.Errorf("failed to create store: %w", err)
    }

    return runService(cCtx, nc, s)
  },
}

func runService(cCtx *cli.Context, nc *nats.Conn, s *store.Store) error {
  svc, err := service.NewService(nc, s)
  if err != nil {
    return err
  }

  if err := svc.Run(cCtx.Context); err != nil {
    return err
  }

  log.Info().Msg("service finished")
  return nil
}
