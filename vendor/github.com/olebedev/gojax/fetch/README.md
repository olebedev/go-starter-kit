# fetch [![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/olebedev/gojax/fetch)

> a window.fetch JavaScript polyfill

### Usage

Install via `go get https://github.com/olebedev/gojax/fetch`.

```go
package main

import (
	"fmt"

	goproxy "gopkg.in/elazarl/goproxy.v1"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/eventloop"
	"github.com/olebedev/gojax/fetch"
)

func main() {
	loop := eventloop.NewEventLoop()
	loop.Start()
	defer loop.Stop()

	fetch.Enable(loop, goproxy.NewProxyHttpServer())

	wait := make(chan string, 1)
	loop.RunOnLoop(func(vm *goja.Runtime) {
		vm.Set("callback", func(call goja.FunctionCall) goja.Value {
			wait <- call.Argument(0).ToString().String()
			return nil
		})

		vm.RunString(`
			fetch('https://ya.ru').then(function(resp){
				return resp.text();
			}).then(function(resp){
				callback(resp.slice(0, 15));
			});
		`)
	})
	fmt.Println(<-wait)
}
```

This program will prints `<!DOCTYPE html>` into stdout. See `fetch_test.go` for more examples.

