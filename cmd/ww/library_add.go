package main

import (
	"github.com/urfave/cli/v2"
	"github.com/wombatwisdom/wombat-builder/library"
)

func AddLibraryCommand() *cli.Command {
	return &cli.Command{
		Name:      "add",
		Usage:     "add a library",
		Args:      true,
		ArgsUsage: " <name>",
		Flags: []cli.Flag{
			LogFlag,
			&cli.StringFlag{
				Name:    "module",
				Aliases: []string{"m"},
				Usage:   "the module name, without the version",
			},
		},
		Action: func(c *cli.Context) error {
			GlobalLogLevelFromFlag(c)

			// -- parse the arguments
			if c.NArg() != 1 {
				return cli.Exit("the name of the library needs to be provided ", 1)
			}

			if !c.IsSet("module") {
				return cli.Exit("the module name must be provided", 1)
			}

			libSpec := library.Spec{
				Name:   c.Args().First(),
				Module: c.String("module"),
			}

			lib := LibFromEnv()
			if err := lib.AddLibrary(libSpec); err != nil {
				return cli.Exit(err, 1)
			}

			return nil
		},
	}
}
