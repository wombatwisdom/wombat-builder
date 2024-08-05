package main

import (
	"github.com/urfave/cli/v2"
)

func PackageCommand() *cli.Command {
	return &cli.Command{
		Name:  "package",
		Usage: "Manage Packages",
		Subcommands: []*cli.Command{
			AddPackageCommand(),
			ListPackageCommand(),
			BuildPackageCommand(),
		},
	}
}
