package main

import (
	"github.com/urfave/cli/v2"
	"github.com/wombatwisdom/wombat-builder/builder"
	"github.com/wombatwisdom/wombat-builder/library"
)

func AddPackageCommand() *cli.Command {
	return &cli.Command{
		Name:  "add",
		Usage: "add a package",
		Description: `
add a package to the store.

The full package name must be provided. This is the full package name as it would be in a go import statement.

The library and version must be specified with the --library and --version flags as does the package with 
the --package flag.

The location of the library is determined by the LIBRARY_DIR environment variable. If not set, the default
location is the library folder in the current directory.
`,
		Args:      true,
		ArgsUsage: " library version pkg",
		Flags: []cli.Flag{
			LogFlag,
			&cli.StringFlag{
				Name:  "path",
				Usage: "the path added to the module to target the package",
			},
		},
		Action: func(c *cli.Context) error {
			GlobalLogLevelFromFlag(c)

			// -- parse the arguments
			if c.NArg() != 3 {
				return cli.Exit("library, version and package ids must be provided", 1)
			}
			ref := builder.PackageRef{
				Library: c.Args().Get(0),
				Version: c.Args().Get(1),
				Package: c.Args().Get(2),
			}

			if !c.IsSet("path") {
				return cli.Exit("path must be provided", 1)
			}

			lc := LibFromEnv()
			pkgSpec := library.PackageSpec{
				Fqn:  c.String("path"),
				Name: ref.Package,
			}

			if err := lc.AddPackage(ref.Library, ref.Version, pkgSpec); err != nil {
				return err
			}

			return nil
		},
	}
}
