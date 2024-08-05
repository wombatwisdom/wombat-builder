package main

import (
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := initApp()

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("Failed to run app")
	}
}

func initApp() *cli.App {
	app := &cli.App{
		Name:  "ww",
		Usage: "Wombat Wisdom Tooling",
		Commands: []*cli.Command{
			LibraryCommand(),
			VersionCommand(),
			PackageCommand(),
		},
	}

	return app
}
