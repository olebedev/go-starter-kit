package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nu7hatch/gouuid"
	"github.com/olebedev/go-duktape"
	"github.com/olebedev/go-duktape-fetch"
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
func (r *React) Handle(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			UUID := c.MustGet("uuid").(*uuid.UUID)
			c.Writer.WriteHeader(http.StatusInternalServerError)
			c.Writer.Header().Add("Content-Type", "text/plain")
			c.Writer.Write([]byte(fmt.Sprintf("req uuid: %s\n%#v", UUID, r)))
			c.Abort()
		}
	}()

	vm := r.get()
	vm.Lock()

	vm.PushGlobalObject()
	vm.GetPropString(-1, "__router__")
	vm.PushString("renderToString")

	req := func() string {
		b, _ := json.Marshal(map[string]interface{}{
			"url":     c.Request.URL.String(),
			"headers": c.Request.Header,
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
			// If no redirection and no error
			c.Writer.WriteHeader(http.StatusOK)
			c.Writer.Header().Add("Content-Type", "text/html")
			c.Writer.Write([]byte("<!doctype html>\n" + re.Body))
			c.Abort()
			// If redirect
		} else if len(re.Redirect) != 0 {
			c.Redirect(http.StatusMovedPermanently, re.Redirect)
			// If internal error
		} else if len(re.Error) != 0 {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			c.Writer.Header().Add("Content-Type", "text/plain")
			c.Writer.Write([]byte(re.Error))
			c.Abort()
		}
	case <-time.After(2 * time.Second):
		// release duktape context
		r.drop(vm)
		UUID := c.MustGet("uuid").(*uuid.UUID)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.Header().Add("Content-Type", "text/html")
		c.Writer.Write([]byte(fmt.Sprintf(`
			<!DOCTYPE html>
			<html>
				<head>
					<meta charset=charset"UTF-8">
					<link rel="stylesheet" href="/static/build/bundle.css">
					<title>Internal Server Error</title>
				</head>
				<body>
					<h1>Internal Server Error</h1>
					<p>uuid: %s</p>
					<div id="app"></div>
					<script async src="/static/build/bundle.js"></script>
				</body>
			</html>
		`, UUID)))
		c.Abort()
	}
}

// Resp is a struct for convinient
// react app response parsing.
type resp struct {
	Error    string `json:"error"`
	Redirect string `json:"redirect"`
	Body     string `json:"body"`
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
	ot.DestroyHeap()
	ot = nil
	o.ch <- newDuktapeContext(o.path, o.engine)
}
