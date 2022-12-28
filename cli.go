package main

import (
	"github.com/urfave/cli/v2"
)

const appName = "certinfo"

var appVersion = "undefined"

func newApp() *cli.App {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Version = appVersion
	app.Name = appName
	app.Usage = "TLS certificate info tool"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "host, h",
			Usage: "FQDN of server to get certificate information",
		},
		&cli.IntFlag{
			Name:  "port, p",
			Usage: "Port number",
			Value: 443,
		},
		&cli.StringFlag{
			Name:  "tls",
			Usage: "Force client TLS version. Valid values are 1.0 to 1.3",
		},
		&cli.BoolFlag{
			Name:  "insecure",
			Usage: "Ignore certificate errors",
		},
	}
	app.Action = getCertInfo
	return app
}
