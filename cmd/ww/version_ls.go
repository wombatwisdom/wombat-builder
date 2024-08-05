package main

import (
	"github.com/urfave/cli/v2"
)

func ListVersionCommand() *cli.Command {
	return &cli.Command{
		Name:      "ls",
		Usage:     "list library versions",
		Args:      true,
		ArgsUsage: " <library>",
		Flags: []cli.Flag{
			LogFlag,
		},
		Action: func(c *cli.Context) error {
			GlobalLogLevelFromFlag(c)

			// -- parse the arguments
			if c.NArg() != 1 {
				return cli.Exit("no library provided ", 1)
			}
			libId := c.Args().First()

			lib := LibFromEnv()
			vers, err := lib.Versions(libId)
			if err != nil {
				return cli.Exit(err, 1)
			}

			for _, v := range vers {
				_, _ = c.App.Writer.Write([]byte(v + "\n"))
			}

			return nil
		},
	}
}
