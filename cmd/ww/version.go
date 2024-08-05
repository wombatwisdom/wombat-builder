package main

import (
	"github.com/urfave/cli/v2"
)

func VersionCommand() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Manage Versions",
		Subcommands: []*cli.Command{
			AddVersionCommand(),
			ListVersionCommand(),
		},
	}
}
