package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"runtime"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/nu7hatch/gouuid"
	"gopkg.in/olebedev/go-duktape-fetch.v2"
	"gopkg.in/olebedev/go-duktape.v2"
)

// React struct is contains duktape
// pool to serve HTTP requests and
// separates some domain specific
// resources.
type React struct {
	pool
	debug bool
	path  string
}

// NewReact initialized React struct
func NewReact(filePath string, debug bool, server http.Handler) *React {
	r := &React{
		debug: debug,
		path:  filePath,
	}
	if !debug {
		r.pool = newDuktapePool(filePath, runtime.NumCPU()+1, server)
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

// Handle handles all HTTP requests which
// have no been caught via static file
// handler or other middlewares.
func (r *React) Handle(c echo.Context) error {

	UUID := c.Get("uuid").(*uuid.UUID)
	defer func() {
		if r := recover(); r != nil {
			c.Render(http.StatusInternalServerError, "react.html", Resp{
				UUID:  UUID.String(),
				Error: r.(string),
			})
		}
	}()

	vm := r.get()

	start := time.Now()
	select {
	case re := <-vm.Handle(map[string]interface{}{
		"url":     c.Request().(*standard.Request).URI(),
		"headers": c.Request().(*standard.Request).Request.Header,
		"uuid":    UUID.String(),
	}):
		re.RenderTime = time.Since(start)
		// Return vm back to the pool
		r.put(vm)
		// Handle the Response
		if len(re.Redirect) == 0 && len(re.Error) == 0 {
			// If no redirection and no errors
			c.Response().Header().Set("X-React-Render-Time", fmt.Sprintf("%s", re.RenderTime))
			return c.Render(http.StatusOK, "react.html", re)
			// If redirect
		} else if len(re.Redirect) != 0 {
			return c.Redirect(http.StatusMovedPermanently, re.Redirect)
			// If internal error
		} else if len(re.Error) != 0 {
			c.Response().Header().Set("X-React-Render-Time", fmt.Sprintf("%s", re.RenderTime))
			return c.Render(http.StatusInternalServerError, "react.html", re)
		}
	case <-time.After(2 * time.Second):
		// release duktape context
		r.drop(vm)
		return c.Render(http.StatusInternalServerError, "react.html", Resp{
			UUID:  UUID.String(),
			Error: "time is out",
		})
	}
	return nil
}

// Resp is a struct for convinient
// react app Response parsing.
// Feel free to add any other keys to this struct
// and return value for this key at ecmascript side.
// Keep it sync with: src/app/client/router/toString.js:23
type Resp struct {
	UUID       string        `json:"uuid"`
	Error      string        `json:"error"`
	Redirect   string        `json:"redirect"`
	App        string        `json:"app"`
	Title      string        `json:"title"`
	Meta       string        `json:"meta"`
	Initial    string        `json:"initial"`
	RenderTime time.Duration `json:"-"`
}

// HTMLApp returns a application template
func (r Resp) HTMLApp() template.HTML {
	return template.HTML(r.App)
}

// HTMLTitle returns a title data
func (r Resp) HTMLTitle() template.HTML {
	return template.HTML(r.Title)
}

// HTMLMeta returns a meta data
func (r Resp) HTMLMeta() template.HTML {
	return template.HTML(r.Meta)
}

// Interface to serve React app on demand or from prepared pool.
type pool interface {
	get() *ReactVM
	put(*ReactVM)
	drop(*ReactVM)
}

// NewDuktapePool return new duktape contexts pool.
func newDuktapePool(filePath string, size int, engine http.Handler) *duktapePool {
	pool := &duktapePool{
		path:   filePath,
		ch:     make(chan *ReactVM, size),
		engine: engine,
	}

	go func() {
		for i := 0; i < size; i++ {
			pool.ch <- newReactVM(filePath, engine)
		}
	}()

	return pool
}

// newReactVM loads bundle.js to context.
func newReactVM(filePath string, engine http.Handler) *ReactVM {

	vm := &ReactVM{
		Context: duktape.New(),
		ch:      make(chan Resp, 1),
	}

	vm.PevalString(`var console = {log:print,warn:print,error:print,info:print}`)
	fetch.PushGlobal(vm.Context, engine)
	app, err := Asset(filePath)
	Must(err)

	// Reduce CGO calls
	vm.PushGlobalGoFunction("__goServerCallback__", func(ctx *duktape.Context) int {
		result := ctx.SafeToString(-1)
		vm.ch <- func() Resp {
			var re Resp
			json.Unmarshal([]byte(result), &re)
			return re
		}()
		return 0
	})

	fmt.Printf("%s loaded\n", filePath)
	if err := vm.PevalString(string(app)); err != nil {
		derr := err.(*duktape.Error)
		panic(derr.Message)
	}
	vm.PopN(vm.GetTop())
	return vm
}

// ReactVM wraps duktape.Context
type ReactVM struct {
	*duktape.Context
	ch chan Resp
}

// Handle handles http requests
func (r *ReactVM) Handle(req map[string]interface{}) <-chan Resp {
	b, err := json.Marshal(req)
	Must(err)
	// Keep it sync with `src/app/client/index.js:4`
	r.PevalString(`main(` + string(b) + `, __goServerCallback__)`)
	return r.ch
}

// DestroyHeap destroys the context's heap
func (r *ReactVM) DestroyHeap() {
	close(r.ch)
	r.Context.DestroyHeap()
}

// Pool's implementations

type onDemandPool struct {
	path   string
	engine http.Handler
}

func (f *onDemandPool) get() *ReactVM {
	return newReactVM(f.path, f.engine)
}

func (f onDemandPool) put(c *ReactVM) {
	c.Lock()
	c.FlushTimers()
	c.Gc(0)
	c.DestroyHeap()
}

func (f *onDemandPool) drop(c *ReactVM) {
	f.put(c)
}

type duktapePool struct {
	ch     chan *ReactVM
	path   string
	engine http.Handler
}

func (o *duktapePool) get() *ReactVM {
	return <-o.ch
}

func (o *duktapePool) put(ot *ReactVM) {
	// Drop any futured async calls
	ot.Lock()
	ot.FlushTimers()
	ot.Unlock()
	o.ch <- ot
}

func (o *duktapePool) drop(ot *ReactVM) {
	ot.Lock()
	ot.FlushTimers()
	ot.Gc(0)
	ot.DestroyHeap()
	ot = nil
	o.ch <- newReactVM(o.path, o.engine)
}
