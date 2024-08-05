package main

import (
	"github.com/urfave/cli/v2"
	"github.com/wombatwisdom/wombat-builder/library"
)

func AddVersionCommand() *cli.Command {
	return &cli.Command{
		Name:      "add",
		Usage:     "add a version to a library",
		Args:      true,
		ArgsUsage: " <library> <name>",
		Flags: []cli.Flag{
			LogFlag,
		},
		Action: func(c *cli.Context) error {
			GlobalLogLevelFromFlag(c)

			// -- parse the arguments
			if c.NArg() != 2 {
				return cli.Exit("the library and the version need to be provided ", 1)
			}
			libId := c.Args().Get(0)

			verSpec := library.VersionSpec{
				Name: c.Args().Get(1),
			}

			lib := LibFromEnv()
			if err := lib.AddVersion(libId, verSpec); err != nil {
				return cli.Exit(err, 1)
			}

			return nil
		},
	}
}
