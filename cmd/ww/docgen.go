package main

import (
	"context"
	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"github.com/wombatwisdom/wombat-builder/docgen"
	"github.com/wombatwisdom/wombat-builder/library"
	"os"
	"path"
)

func DocGenCommand() *cli.Command {
	return &cli.Command{
		Name:      "docgen",
		Usage:     "generate docs for a package",
		Args:      true,
		ArgsUsage: "<library> <version> <package>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output_dir",
				Aliases: []string{"o"},
				Usage:   "the output directory",
			},
			&cli.StringFlag{
				Name:  "loglevel",
				Usage: "set the log level",
				Value: "warn",
			},
		},
		Action: func(c *cli.Context) error {
			libDir := os.Getenv("LIBRARY_DIR")
			if libDir == "" {
				libDir = "library/libraries"
			}

			// -- parse the arguments
			if c.NArg() != 3 {
				return cli.Exit("library, version and package must be provided as arguments", 1)
			}

			libId := c.Args().Get(0)
			verId := c.Args().Get(1)
			pkg := c.Args().Get(2)

			outputDir, err := os.Getwd()
			if err != nil {
				return cli.Exit(err, 1)
			}
			if c.IsSet("output_dir") {
				outputDir = path.Join(outputDir, c.String("output_dir"))
			}

			ll := "warn"
			if c.IsSet("loglevel") {
				ll = c.String("loglevel")
			}
			logLevel, err := zerolog.ParseLevel(ll)
			if err != nil {
				return cli.Exit("invalid log level", 1)
			}
			zerolog.SetGlobalLevel(logLevel)

			lib := library.NewFsClient(libDir)
			dg := docgen.NewDocsGenerator(lib)

			if err := dg.GenerateForPackage(context.Background(), libId, verId, pkg); err != nil {
				return cli.Exit(err, 1)
			}

			color.Green("documentation generated for %s %s %s", libId, verId, pkg)
			return nil
		},
	}
}
