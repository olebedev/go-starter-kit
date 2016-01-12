package main

import (
	"fmt"
	"os"

	"./server"

	"github.com/codegangsta/cli"
)

func main() {
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
		{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "Prints app's version",
			Action:  Version,
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

var githash, gittag, buildstamp string

// Version prints git commit hash,
// date time and vestion from tag
func Version(c *cli.Context) {
	fmt.Printf(`Git tag: %s
Git Commit Hash: %s
UTC Build Time: %s
`, gittag, githash, buildstamp)
}
