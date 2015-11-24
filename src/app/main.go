package main

import (
	"app/server"
	"fmt"
	"os"

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
		Config: c.String("config"),
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
