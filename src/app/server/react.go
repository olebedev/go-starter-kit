package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"runtime"
	"time"

	"github.com/labstack/echo"
	"github.com/nu7hatch/gouuid"
	"gopkg.in/olebedev/go-duktape-fetch.v1"
	"gopkg.in/olebedev/go-duktape.v1"
)

// NewReact initialized React struct
func NewReact(filePath string, debug bool, server http.Handler) *React {
	r := &React{
		debug: debug,
		path:  filePath,
	}
	if !debug {
		r.pool = newDuktapePool(filePath, runtime.NumCPU(), server)
	} else {
		// Use onDemandPool to load full react
		// app each time for any http requests.
		// Useful to debug the app.
		r.pool = &onDemandPool{
			path:   filePath,
			engine: server,
		}
	}
	return r
}

// React struct is contains duktape
// pool to serve HTTP requests and
// separates some domain specific
// resources.
type React struct {
	pool
	debug bool
	path  string
}

// Handle handles all HTTP requests which
// have no been caught via static file
// handler or other middlewares.
func (r *React) Handle(c *echo.Context) error {
	UUID := c.Get("uuid").(*uuid.UUID)
	defer func() {
		if r := recover(); r != nil {
			c.Render(http.StatusInternalServerError, "react.html", resp{
				UUID:  UUID.String(),
				Error: r.(string),
			})
		}
	}()

	vm := r.get()
	vm.Lock()

	vm.PushGlobalObject()
	vm.GetPropString(-1, "__router__")
	vm.PushString("renderToString")

	req := func() string {
		b, _ := json.Marshal(map[string]interface{}{
			"url":     c.Request().URL.String(),
			"headers": c.Request().Header,
			"uuid":    UUID.String(),
		})
		return string(b)
	}()
	vm.PushString(req)
	vm.JsonDecode(-1)
	ch := make(chan *resp, 1)
	vm.PushGoFunction(func(ctx *duktape.Context) int {

		// Getting response object via json
		r := func() *resp {
			var re resp
			json.Unmarshal([]byte(vm.JsonEncode(-1)), &re)
			return &re
		}()

		// Unlock handler
		ch <- r
		// Return nothing into duktape context
		return 0
	})
	vm.PcallProp(1, 2)
	vm.Unlock()

	// Lock handler and wait for js app response
	select {
	case re := <-ch:
		// Hold the context. This call is really important
		// because async calls is possible. And we cannot
		// allow to break the context stack.
		vm.Lock()
		// Clean duktape vm stack
		vm.PopN(vm.GetTop())

		// Drop any futured async calls
		vm.ResetTimers()
		// Release the context
		vm.Unlock()
		// Return vm back to the pool
		r.put(vm)
		// Handle the response
		if len(re.Redirect) == 0 && len(re.Error) == 0 {
			// If no redirection and no errors
			return c.Render(http.StatusOK, "react.html", re)
			// If redirect
		} else if len(re.Redirect) != 0 {
			return c.Redirect(http.StatusMovedPermanently, re.Redirect)
			// If internal error
		} else if len(re.Error) != 0 {
			return c.Render(http.StatusInternalServerError, "react.html", re)
		}
	case <-time.After(2 * time.Second):
		// release duktape context
		r.drop(vm)
		return c.Render(http.StatusInternalServerError, "react.html", resp{
			UUID:  UUID.String(),
			Error: "time is out",
		})
	}
	return nil
}

// Resp is a struct for convinient
// react app response parsing.
// Feel free to add any other keys to this struct
// and return value for this key at ecmascript side.
// Keep it sync with: src/app/client/router/toString.js:23
type resp struct {
	UUID     string `json:"uuid"`
	Error    string `json:"error"`
	Redirect string `json:"redirect"`
	App      string `json:"app"`
	Title    string `json:"title"`
	Meta     string `json:"meta"`
	Initial  string `json:"initial"`
}

func (r resp) HTMLApp() template.HTML {
	return template.HTML(r.App)
}

func (r resp) HTMLMeta() template.HTML {
	return template.HTML(r.Meta)
}

// Interface to serve React app on demand or from prepared pool.
type pool interface {
	get() *duktape.Context
	put(*duktape.Context)
	drop(*duktape.Context)
}

// NewDuktapePool return new duktape contexts pool.
func newDuktapePool(filePath string, size int, engine http.Handler) *duktapePool {
	pool := &duktapePool{
		path:   filePath,
		ch:     make(chan *duktape.Context, size),
		engine: engine,
	}

	go func() {
		for i := 0; i < size; i++ {
			pool.ch <- newDuktapeContext(filePath, engine)
		}
	}()

	return pool
}

// NewDuktapeContext loads bundle.js to context.
func newDuktapeContext(filePath string, engine http.Handler) *duktape.Context {
	vm := duktape.New()
	vm.PevalString(`var console = {log:print,warn:print,error:print,info:print}`)
	fetch.Define(vm, engine)
	app, err := Asset(filePath)
	Must(err)
	fmt.Printf("%s loaded\n", filePath)
	if err := vm.PevalString(string(app)); err != nil {
		derr := err.(*duktape.Error)
		fmt.Printf("\n\n\n%v\n%v\n\n\n", derr.FileName, derr.LineNumber)
		panic(derr.Message)
	}
	vm.PopN(vm.GetTop())
	return vm
}

// Pool's implementations

type onDemandPool struct {
	path   string
	engine http.Handler
}

func (f *onDemandPool) get() *duktape.Context {
	return newDuktapeContext(f.path, f.engine)
}

func (f onDemandPool) put(c *duktape.Context) {
	c.Lock()
	c.ResetTimers()
	c.Gc(0)
	c.DestroyHeap()
}

func (f *onDemandPool) drop(c *duktape.Context) {
	f.put(c)
}

type duktapePool struct {
	ch     chan *duktape.Context
	path   string
	engine http.Handler
}

func (o *duktapePool) get() *duktape.Context {
	return <-o.ch
}

func (o *duktapePool) put(ot *duktape.Context) {
	ot.Gc(0)
	o.ch <- ot
}

func (o *duktapePool) drop(ot *duktape.Context) {
	ot.Lock()
	ot.ResetTimers()
	ot.DestroyHeap()
	ot = nil
	o.ch <- newDuktapeContext(o.path, o.engine)
}
