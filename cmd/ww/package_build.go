package main

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/wombatwisdom/wombat-builder/builder"
	"github.com/wombatwisdom/wombat-builder/docgen"
	"github.com/wombatwisdom/wombat-builder/library"
)

func BuildPackageCommand() *cli.Command {
	return &cli.Command{
		Name:  "build",
		Usage: "build a package",
		Description: `
build a package and store it.

The location of the library is determined by the LIBRARY_DIR environment variable. If not set, the default
location is the library folder in the current directory.
`,
		Args:      true,
		ArgsUsage: " library version pkg",
		Flags: []cli.Flag{
			LogFlag,
			&cli.StringFlag{
				Name:  "goos",
				Usage: "go OS to build for",
			},
			&cli.StringFlag{
				Name:  "goarch",
				Usage: "go architecture to build for",
			},
			&cli.StringFlag{
				Name:  "go-exec",
				Usage: "the go executable to use",
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

			if !c.IsSet("goos") {
				return cli.Exit("goos must be provided", 1)
			}

			if !c.IsSet("goarch") {
				return cli.Exit("goarch must be provided", 1)
			}

			goExec := "go"
			if c.IsSet("go-exec") {
				goExec = c.String("go-exec")
			}

			lc := LibFromEnv()
			if err := add(context.Background(), lc, ref, goExec, c.String("goos"), c.String("goarch")); err != nil {
				return cli.Exit(err, 1)
			}

			return nil
		},
	}
}

func add(ctx context.Context, lc library.Client, ref builder.PackageRef, goExec string, goos string, goarch string) error {
	logger := log.With().Str("library", ref.Library).Str("version", ref.Version).Str("package", ref.Package).Logger()

	lib, err := lc.Library(ref.Library)
	if err != nil {
		return err
	}

	pkg, err := lc.Package(ref.Library, ref.Version, ref.Package)
	if err != nil {
		return err
	}

	spec := builder.PackageBuildSpec{
		PackageRef: ref,
		Module:     lib.Module,
		Path:       pkg.Fqn,
		Os:         goos,
		Arch:       goarch,
	}

	bc := builder.NewBuilder(lc)
	if err := bc.BuildPackage(ctx, spec, goExec); err != nil {
		return err
	}
	log.Debug().Str("os", goos).Str("arch", goarch).Msg("package built")

	gen := docgen.NewDocsGenerator(lc)
	if err := gen.GenerateForPackage(ctx, ref.Library, ref.Version, ref.Package); err != nil {
		return err
	}
	logger.Info().Msg("docs generated")

	return nil
}
