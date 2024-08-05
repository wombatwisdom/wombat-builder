package main

import (
	"context"
	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"github.com/wombatwisdom/wombat-builder/builder"
	"github.com/wombatwisdom/wombat-builder/library"
	"os"
	"path"
)

func BuildCommand() *cli.Command {
	return &cli.Command{
		Name:  "build",
		Usage: "build a package",
		Description: `
Build a package from a specific library version.

The library and version must be specified with the --library and --version flags as does the package with 
the --package flag.

The location of the library is determined by the LIBRARY_DIR environment variable. If not set, the default
location is the library folder in the current directory.
`,
		Args:      true,
		ArgsUsage: "<library> <version> <package>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output_dir",
				Aliases: []string{"o"},
				Usage:   "the output directory",
			},
			&cli.StringFlag{
				Name:  "os",
				Usage: "the target operating system",
			},
			&cli.StringFlag{
				Name:  "arch",
				Usage: "the target architecture",
			},
			&cli.StringFlag{
				Name:  "go-exec",
				Usage: "the path to the go binary to use",
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

			if !c.IsSet("os") {
				return cli.Exit("os must be set", 1)
			}

			if !c.IsSet("arch") {
				return cli.Exit("arch must be set", 1)
			}

			goExec := "go"
			if c.IsSet("go-exec") {
				goExec = c.String("go-exec")
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

			spec := builder.PackageBuildSpec{
				PackageRef: builder.PackageRef{
					Library: libId,
					Version: verId,
					Package: pkg,
				},
				Os:   c.String("os"),
				Arch: c.String("arch"),
			}

			lib := library.NewFsClient(libDir)
			bc := builder.NewBuilder(lib)

			ctx := context.Background()
			if err := bc.BuildPackage(ctx, spec, goExec); err != nil {
				return cli.Exit(err, 1)
			}

			color.Green("%s %s %s built for %s/%s", libId, verId, pkg, c.String("os"), c.String("arch"))

			return nil
		},
	}
}
