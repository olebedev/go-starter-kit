# go-duktape-fetch [![wercker status](https://app.wercker.com/status/fb4b5e19e7981f6aa9b0426deeaa1406/s "wercker status")](https://app.wercker.com/project/bykey/fb4b5e19e7981f6aa9b0426deeaa1406)

Fetch polyfill for [go-duktape](https://github.com/olebedev/go-duktape).

### Usage

First of all install the package `go get gopkg.in/olebedev/go-duktape-fetch.v2`.

```go
package main

import (
  "fmt"

  "gopkg.in/olebedev/go-duktape.v2"
  "gopkg.in/olebedev/go-duktape-fetch.v2"
)

func main() {
  // create an ecmascript context
  ctx := duktape.New()
  // push fetch into the global scope
  fetch.PushGlobal(ctx, nil)
  ch := make(chan string)
  ctx.PushGlobalGoFunction("cbk", func(c *duktape.Context) int {
    ch <- c.SafeToString(-1)
    return 0
  })
  // make a request
  ctx.PevalString(`
    fetch('http://ya.ru').then(function(resp) {
      return resp.text();
    }).then(function(body) {
      // release channel
      cbk(body.slice(0, 15));
    });
  `)
  // print head line of the response body
  fmt.Println(<-ch)
}
```
This program will prodice `<!DOCTYPE html>` into stdout.

Also you can define fetch with some other good possibilities. In 
particular you can specify local http instance with `http.Handler` 
interface as a the second parameter. It will be used for all local
requests which url starts with `/`(single slash). See [tests](https://github.com/olebedev/go-duktape-fetch/blob/master/fetch_test.go) 
for more detail.
