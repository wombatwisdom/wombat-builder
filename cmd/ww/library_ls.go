package main

import (
	"github.com/urfave/cli/v2"
)

func ListLibraryCommand() *cli.Command {
	return &cli.Command{
		Name:    "ls",
		Usage:   "list libraries",
		Aliases: []string{"list"},
		Flags: []cli.Flag{
			LogFlag,
		},
		Action: func(c *cli.Context) error {
			GlobalLogLevelFromFlag(c)

			lib := LibFromEnv()
			libs, err := lib.Libraries()
			if err != nil {
				return cli.Exit(err, 1)
			}

			for _, l := range libs {
				_, _ = c.App.Writer.Write([]byte(l + "\n"))
			}

			return nil
		},
	}
}
