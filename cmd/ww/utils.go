package main

import (
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"github.com/wombatwisdom/wombat-builder/library"
	"os"
)

var (
	LogFlag = &cli.StringFlag{
		Name:  "loglevel",
		Usage: "set the log level",
		Value: "warn",
	}
)

func GlobalLogLevelFromFlag(c *cli.Context) {
	ll := "warn"
	if c.IsSet("loglevel") {
		ll = c.String("loglevel")
	}
	logLevel, err := zerolog.ParseLevel(ll)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLevel)
}

func LibFromEnv() library.Client {
	libDir := os.Getenv("LIBRARY_DIR")
	if libDir == "" {
		libDir = "library/libraries"
	}
	return library.NewFsClient(libDir)
}
