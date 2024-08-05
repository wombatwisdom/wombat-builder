package main

import (
	"github.com/urfave/cli/v2"
)

func ListPackageCommand() *cli.Command {
	return &cli.Command{
		Name:      "ls",
		Usage:     "list packages for a library versions",
		Args:      true,
		ArgsUsage: " <library> <version>",
		Flags: []cli.Flag{
			LogFlag,
		},
		Action: func(c *cli.Context) error {
			GlobalLogLevelFromFlag(c)

			// -- parse the arguments
			if c.NArg() != 2 {
				return cli.Exit("no library and version provided ", 1)
			}
			libId := c.Args().Get(0)
			verId := c.Args().Get(1)

			lib := LibFromEnv()
			pkgs, err := lib.Packages(libId, verId)
			if err != nil {
				return cli.Exit(err, 1)
			}

			for _, v := range pkgs {
				_, _ = c.App.Writer.Write([]byte(v + "\n"))
			}

			return nil
		},
	}
}
