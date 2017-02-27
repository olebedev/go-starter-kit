package fetch

import (
	"net/http"
	"sync"
	"testing"

	goproxy "gopkg.in/elazarl/goproxy.v1"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/eventloop"
	"github.com/stretchr/testify/require"
)

func TestEnable(t *testing.T) {
	loop := eventloop.NewEventLoop()
	loop.Start()
	defer loop.Stop()

	Enable(loop, goproxy.NewProxyHttpServer())

	var v goja.Value
	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	loop.RunOnLoop(func(vm *goja.Runtime) {
		v, err = vm.RunString(`typeof fetch`)
		wg.Done()
	})

	wg.Wait()
	require.Nil(t, err)
	require.NotNil(t, v)
	require.Equal(t, "function", v.ToString().String())
}

func TestRequest(t *testing.T) {
	loop := eventloop.NewEventLoop()
	loop.Start()
	defer loop.Stop()

	Enable(loop, goproxy.NewProxyHttpServer())

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

	require.Equal(t, "<!DOCTYPE html>", <-wait)
}

func TestCustom(t *testing.T) {
	loop := eventloop.NewEventLoop()
	loop.Start()
	defer loop.Stop()

	Enable(loop, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(`{"error": "method not allowed"}`))
	}))

	wait := make(chan string, 3)
	loop.RunOnLoop(func(vm *goja.Runtime) {
		vm.Set("callback", func(call goja.FunctionCall) goja.Value {
			wait <- call.Argument(0).ToString().String()
			return nil
		})

		vm.RunString(`
			fetch('https://ya.ru').then(function(resp){
				callback(resp.ok);
				callback(resp.status);
				callback(resp.url);
				callback(resp.method);
				callback(resp.headers.get('content-type'));
				return resp.json();
			}).then(function(resp){
				callback(resp.error);
			});
		`)
	})

	require.Equal(t, "false", <-wait)
	require.Equal(t, "405", <-wait)
	require.Equal(t, "https://ya.ru", <-wait)
	require.Equal(t, "GET", <-wait)
	require.Equal(t, "application/json", <-wait)
	require.Equal(t, "method not allowed", <-wait)
}
