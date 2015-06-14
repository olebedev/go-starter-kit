package server

import (
	"app/server/data"
	"app/server/react"
	. "app/server/utils"

	"github.com/codegangsta/cli"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
)

func Run(args []string) {

	app := cli.NewApp()
	app.Name = "app"
	app.Usage = "lmbd landing server application"

	configFlag := cli.StringFlag{
		Name:   "config, c",
		Value:  "local",
		Usage:  "configuration section name",
		EnvVar: "CONFIG",
	}

	app.Commands = []cli.Command{
		{
			Name:   "run",
			Usage:  "Runs server",
			Action: runServer,
			Flags:  []cli.Flag{configFlag},
		},
	}
	app.Run(args)
}

func runServer(c *cli.Context) {
	kit := NewKit(c, conf)
	kit.Engine = gin.Default()

	kit.Engine.StaticFS("/static", &assetfs.AssetFS{
		Asset:    data.Asset,
		AssetDir: data.AssetDir,
		Prefix:   "static",
	})

	// routes.bind(r)
	react.Bind(kit)
	Must(kit.Engine.Run(":" + kit.Conf.UString("port")))
}
