package server

import (
	"app/server/api"
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

	kit.Engine.Use(func(c *gin.Context) {
		c.Set("kit", kit)
	})

	// Avoid favicon react handling
	kit.Engine.GET("/favicon.ico", func(c *gin.Context) {
		c.Redirect(301, "/static/images/favicon.ico")
	})

	kit.Engine.GET("/api/v1/conf", api.ConfHandler)
	react.Bind(kit)
	Must(kit.Engine.Run(":" + kit.Conf.UString("port")))
}
