package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"text/template"

	"github.com/codegangsta/cli"
	"gopkg.in/fsnotify.v1"
)

func main() {
	app := cli.NewApp()
	app.Usage = "cli for fsnotify"
	app.UsageText = `on [options] [path]

   Path could be mix of folders and filepaths, default is '.'.
   Regular usecase is watch the file's changes and execute some
   bash script or command line tool. It could be done in this way:

   on | xargs -n1 -I{} <program>`

	app.Author = "olebedev <ole6edev@gmail.com>"
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "template,t",
			Usage: "output template to render received event, see: https://godoc.org/gopkg.in/fsnotify.v1#Event",
			Value: "{{.Name}}",
		},
		cli.IntFlag{
			Name:  "mask,m",
			Usage: "event's bitwise mask, see: https://godoc.org/gopkg.in/fsnotify.v1#Op",
			Value: 15,
		},
		cli.BoolFlag{
			Name:  "r",
			Usage: "watch given paths recursively",
		},
	}
	app.Action = func(c *cli.Context) {
		t, err := template.New("output").Parse(c.String("template"))
		if err != nil {
			log.Fatal(err)
		}
		id := fsnotify.Op(c.Int("mask"))

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		done := make(chan os.Signal)
		signal.Notify(done, os.Interrupt)

		go func() {
			for {
				select {
				case event := <-watcher.Events:
					// if bitwise mask covers the event
					if event.Op&id == event.Op && event.Name != "" {
						fmt.Println(render(t, event))
					}
				case err := <-watcher.Errors:
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
					}
				}
			}
		}()

		args := c.Args()
		if len(args) == 0 {
			args = []string{"."}
		}

		for _, arg := range args {
			if err := addPath(watcher, arg, c.Bool("r")); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		<-done
		watcher.Close()
	}

	app.Run(os.Args)
}

func render(t *template.Template, e fsnotify.Event) string {
	var doc bytes.Buffer
	t.Execute(&doc, e)

	return doc.String()
}

func addPath(w *fsnotify.Watcher, name string, recursively bool) error {
	f, err := os.Stat(name)
	if err != nil {
		return err
	}

	if !f.IsDir() || !recursively {
		return w.Add(name)
	}

	return filepath.Walk(name, func(p string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if f.IsDir() {
			return w.Add(p)
		}
		return nil
	})
}
