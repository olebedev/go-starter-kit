package main

import (
	"os"

	"./server"

	"github.com/codegangsta/cli"
)

// declare vars to change it with ldflags
var (
	debug      = "false"
	commitHash = "0"
)

func main() {
	server.Debug = debug
	server.CommitHash = commitHash
	Run(os.Args)
}

// Run creates, configures and runs
// main cli.App
func Run(args []string) {

	app := cli.NewApp()
	app.Name = "app"
	app.Usage = "React server application"

	app.Commands = []cli.Command{
		{
			Name:   "run",
			Usage:  "Runs server",
			Action: RunServer,
		},
	}
	app.Run(args)
}

// RunServer creates, configures and runs
// main server.App
func RunServer(c *cli.Context) {
	app := server.NewApp(server.AppOptions{
	// see server/app.go:150
	})
	app.Run()
}
