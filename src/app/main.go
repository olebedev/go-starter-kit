package main

import (
	"app/server"
	"os"
	"runtime"

	"github.com/codegangsta/cli"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	Run(os.Args)
}

// Run creates, configures and runs
// main cli.App
func Run(args []string) {

	app := cli.NewApp()
	app.Name = "app"
	app.Usage = "React server application"

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
			Action: RunServer,
			Flags:  []cli.Flag{configFlag},
		},
	}
	app.Run(args)
}

// RunServer creates, configures and runs
// main server.App
func RunServer(c *cli.Context) {
	app := server.NewApp(server.AppOptions{
		Config: c.String("config"),
	})
	app.Run()
}
