package cmd

import (
  "fmt"
  "github.com/benthosdev/benthos-builder/internal/cmd"
  "github.com/benthosdev/benthos-builder/internal/store"
  "github.com/rs/zerolog/log"
  "github.com/urfave/cli/v2"
  "runtime"
  "sync"
)

var AllCommand = &cli.Command{
  Name:  "all",
  Usage: "run the builder, service and api within the same process",
  Flags: append(cmd.NatsFlags, []cli.Flag{
    &cli.IntFlag{
      Name:  "port",
      Usage: "the port to run the api on",
      Value: 4430,
    },
    &cli.IntFlag{
      Name:    "workers",
      Usage:   "the number of workers to run",
      Value:   runtime.NumCPU(),
      EnvVars: []string{"WORKERS"},
    },
    &cli.BoolFlag{
      Name:  "ui",
      Usage: "enable the ui",
      Value: false,
    },
  }...),
  Action: func(cCtx *cli.Context) error {
    nc, js, err := cmd.ConnectNats(cCtx)
    if err != nil {
      return fmt.Errorf("failed to connect to nats: %w", err)
    }
    defer nc.Close()

    s, err := store.NewStore(js, true)
    if err != nil {
      return fmt.Errorf("failed to create store: %w", err)
    }

    wg := &sync.WaitGroup{}
    wg.Add(3)

    go func() {
      if err := runBuilder(cCtx, s); err != nil {
        log.Panic().Err(err).Msg("builder failed")
        wg.Done()
      } else {
        wg.Add(-1)
      }
    }()

    go func() {
      if err := runService(cCtx, nc, s); err != nil {
        log.Panic().Err(err).Msg("service failed")
        wg.Done()
      } else {
        wg.Add(-1)
      }
    }()

    go func() {
      if err := runApi(cCtx, nc, js); err != nil {
        log.Panic().Err(err).Msg("api failed")
        wg.Done()
      } else {
        wg.Add(-1)
      }
    }()

    wg.Wait()
    return nil
  },
}
