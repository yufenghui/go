package main

import (
	"github.com/urfave/cli"
	"github.com/yufenghui/go/gravitee/cmd"
	"log"
	"os"
)

var (
	configFile string
	app        *cli.App
)

func init() {

	app = cli.NewApp()

	app.Name = "gravitee-oauth2-server"
	app.Usage = "Gravitee OAuth 2.0 Server"
	app.Version = "1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config",
			Value:       "config.yml",
			Destination: &configFile,
		},
	}

}

func main() {

	app.Commands = []cli.Command{
		{
			Name:    "runserver",
			Aliases: []string{"run"},
			Usage:   "run web server",
			Action: func(c *cli.Context) error {
				return cmd.RunServer(configFile)
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
