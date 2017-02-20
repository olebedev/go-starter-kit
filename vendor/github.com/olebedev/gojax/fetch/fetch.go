package fetch

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/textproto"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/eventloop"
	"github.com/olebedev/gojax/fetch/internal/data"
	"github.com/pkg/errors"
)

// Enable enables fetch for the instance. Loop instance is required instead of
// flat goja's. B/c fetch polyfill uses timeouts for promises.
//
// The second parameter could be any http handler. Even you local instance,
// to handle http requests locally programmatically.
func Enable(loop *eventloop.EventLoop, proxy http.Handler) error {
	if proxy == nil {
		return errors.New("proxy handler cannot be nil")
	}

	script := string(data.MustAsset("dist/bundle.js"))
	prg, err := goja.Compile("fetch.js", script, false)
	if err != nil {
		return errors.Wrap(err, "compile script")
	}
	loop.RunOnLoop(func(vm *goja.Runtime) {
		vm.Set("__fetch__", request(loop, proxy))
		_, err := vm.RunProgram(prg)
		if err != nil {
			panic(err)
		}
	})

	return nil
}

func request(loop *eventloop.EventLoop, proxy http.Handler) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if fn, ok := goja.AssertFunction(call.Argument(2)); ok {
			u := call.Argument(0).String()
			o := call.Argument(1).Export().(map[string]interface{})

			go func() {
				var body io.Reader
				method := http.MethodGet
				header := make(http.Header)

				if h, ex := o["headers"]; ex {
					if he, ok := h.(http.Header); ok {
						for key, value := range he {
							header[textproto.CanonicalMIMEHeaderKey(key)] = value
						}
					}
				}

				if b, ex := o["body"]; ex {
					if bo, ok := b.(string); ok {
						body = bytes.NewBufferString(bo)
					}
				}

				if m, ex := o["method"]; ex {
					if me, ok := m.(string); ok {
						method = me
					}
				}

				var toRet map[string]interface{}

				res := httptest.NewRecorder()
				req, err := http.NewRequest(method, u, body)
				if err != nil {
					toRet = map[string]interface{}{
						"body":    fmt.Sprintf("Internal Server Error: %s", err.Error()),
						"headers": make(map[string][]string),
						"status":  http.StatusInternalServerError,
						"method":  method,
						"url":     u,
					}
				} else {
					req.Header = header
					proxy.ServeHTTP(res, req)
					toRet = map[string]interface{}{
						"body":    res.Body.String(),
						"headers": map[string][]string(res.Header()),
						"status":  res.Code,
						"method":  method,
						"url":     u,
					}
				}
				loop.RunOnLoop(func(vm *goja.Runtime) { fn(nil, vm.ToValue(toRet)) })
			}()
		}
		return nil
	}
}
