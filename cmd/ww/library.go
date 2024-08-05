package main

import (
	"github.com/urfave/cli/v2"
)

func LibraryCommand() *cli.Command {
	return &cli.Command{
		Name:  "library",
		Usage: "Manage libraries",
		Subcommands: []*cli.Command{
			AddLibraryCommand(),
			ListLibraryCommand(),
		},
	}
}
