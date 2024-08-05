package cmd

import (
  "fmt"
  "github.com/nats-io/nats.go"
  "github.com/nats-io/nats.go/jetstream"
  "github.com/rs/zerolog/log"
  "github.com/urfave/cli/v2"
  "github.com/wombatwisdom/wombat-builder/internal/api"
  "github.com/wombatwisdom/wombat-builder/internal/cmd"
)

var ApiCommand = &cli.Command{
	Name:  "api",
	Usage: "run the builder api",
	Description: `
The builder api exposes a rest api that can be used to manage the process of building artifacts. 
`,
	Flags: append(cmd.NatsFlags, []cli.Flag{
		&cli.IntFlag{
			Name:  "port",
			Usage: "the port to run the api on",
			Value: 4430,
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

		return runApi(cCtx, nc, js)
	},
}

func runApi(cCtx *cli.Context, nc *nats.Conn, js jetstream.JetStream) error {
	a, err := api.NewApi(nc, js, cCtx.Int("port"), cCtx.Bool("ui"))
	if err != nil {
		return err
	}
	if err := a.Run(cCtx.Context); err != nil {
		return err
	}

	log.Info().Msg("api finished")
	return nil
}
