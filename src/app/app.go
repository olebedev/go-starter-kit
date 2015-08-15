package main

import (
	app "app/server"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	app.Run(os.Args)
}
