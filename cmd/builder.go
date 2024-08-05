package cmd

import (
  "fmt"
  "github.com/rs/zerolog/log"
  "github.com/urfave/cli/v2"
  "github.com/wombatwisdom/wombat-builder/internal/builder"
  "github.com/wombatwisdom/wombat-builder/internal/cmd"
  "github.com/wombatwisdom/wombat-builder/internal/store"
  "runtime"
)

var BuilderCommand = &cli.Command{
	Name:  "builder",
	Usage: "run the builder",
	Description: `
The builder contains serveral workers which take up the task of building artifacts.
The number of workers can be configured using the --workers flag and is set to the number of cpu's by default'.
    `,
	Flags: append(cmd.NatsFlags, []cli.Flag{
		&cli.IntFlag{
			Name:    "workers",
			Usage:   "the number of workers to run",
			Value:   runtime.NumCPU(),
			EnvVars: []string{"WORKERS"},
		},
	}...),
	Action: func(cCtx *cli.Context) error {
		nc, js, err := cmd.ConnectNats(cCtx)
		if err != nil {
			return fmt.Errorf("failed to connect to nats: %w", err)
		}
		defer nc.Close()

		s, err := store.NewStore(js, false)
		if err != nil {
			return fmt.Errorf("failed to create store: %w", err)
		}

		return runBuilder(cCtx, s)
	},
}

func runBuilder(cCtx *cli.Context, s *store.Store) error {
	bldr, err := builder.NewBuilder(s, cCtx.Int("workers"))
	if err != nil {
		return err
	}

	if err := bldr.Run(cCtx.Context); err != nil {
		return err
	}

	log.Info().Msg("builder finished")
	return nil
}
