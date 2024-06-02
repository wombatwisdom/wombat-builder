package main

import (
  "fmt"
  "github.com/rs/zerolog/log"
  "github.com/urfave/cli/v2"
  "os"
)

func main() {
  app := &cli.App{
    Commands: []*cli.Command{
      {
        Name:  "builder",
        Usage: "run the builder",
        Action: func(cCtx *cli.Context) error {
          fmt.Println("added task: ", cCtx.Args().First())
          return nil
        },
      },
      {
        Name:  "service",
        Usage: "run the builder service",
        Flags: []cli.Flag{
          &cli.StringFlag{
            Name:    "nats-url",
            Usage:   "the nats url to connect to",
            Value:   "tls://connect.ngs.global",
            EnvVars: []string{"NATS_URL"},
          },
        },
        Action: func(cCtx *cli.Context) error {
          fmt.Println("completed task: ", cCtx.Args().First())
          return nil
        },
      },
    },
  }

  if err := app.Run(os.Args); err != nil {
    log.Fatal().Err(err)
  }
}
