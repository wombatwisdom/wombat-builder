package main

import (
  "github.com/rs/zerolog/log"
  "github.com/urfave/cli/v2"
  "github.com/wombatwisdom/wombat-builder/cmd"
  "os"
)

func main() {
	app := &cli.App{
		Name:  "wombat-builder",
		Usage: "a collection of tools for generating wombat distributions",
		Commands: []*cli.Command{
			cmd.BuilderCommand,
			cmd.ServiceCommand,
			cmd.ApiCommand,
			cmd.AllCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("failed to run")
	}
}
