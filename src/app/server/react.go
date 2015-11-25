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

	start := time.Now()
	select {
	case re := <-vm.Handle(map[string]interface{}{
		"url":     c.Request().URL.String(),
		"headers": c.Request().Header,
		"uuid":    UUID.String(),
	}):
		re.RenderTime = time.Since(start)
		// Return vm back to the pool
		r.put(vm)
		// Handle the response
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
	UUID       string        `json:"uuid"`
	Error      string        `json:"error"`
	Redirect   string        `json:"redirect"`
	App        string        `json:"app"`
	Title      string        `json:"title"`
	Meta       string        `json:"meta"`
	Initial    string        `json:"initial"`
	RenderTime time.Duration `json:"-"`
}

func (r resp) HTMLApp() template.HTML {
	return template.HTML(r.App)
}

func (r resp) HTMLMeta() template.HTML {
	return template.HTML(r.Meta)
}

// Interface to serve React app on demand or from prepared pool.
type pool interface {
	get() *ReactVm
	put(*ReactVm)
	drop(*ReactVm)
}

// NewDuktapePool return new duktape contexts pool.
func newDuktapePool(filePath string, size int, engine http.Handler) *duktapePool {
	pool := &duktapePool{
		path:   filePath,
		ch:     make(chan *ReactVm, size),
		engine: engine,
	}

	go func() {
		for i := 0; i < size; i++ {
			pool.ch <- newReactVm(filePath, engine)
		}
	}()

	return pool
}

// newReactVm loads bundle.js to context.
func newReactVm(filePath string, engine http.Handler) *ReactVm {

	vm := &ReactVm{
		Context: duktape.New(),
		ch:      make(chan resp, 1),
	}

	vm.PevalString(`var console = {log:print,warn:print,error:print,info:print}`)
	fetch.PushGlobal(vm.Context, engine)
	app, err := Asset(filePath)
	Must(err)

	// Reduce CGO calls
	vm.PushGlobalGoFunction("__goServerCallback__", func(ctx *duktape.Context) int {
		result := ctx.SafeToString(-1)
		vm.ch <- func() resp {
			var re resp
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

// ReactVm wraps duktape.Context
type ReactVm struct {
	*duktape.Context
	ch chan resp
}

func (r *ReactVm) Handle(req map[string]interface{}) <-chan resp {
	b, err := json.Marshal(req)
	Must(err)
	// Keep it sync with `src/app/client/index.js:1`
	r.PevalString(`__router__.renderToString(` + string(b) + `, __goServerCallback__)`)
	return r.ch
}

func (r *ReactVm) DestroyHeap() {
	close(r.ch)
	r.Context.DestroyHeap()
}

// Pool's implementations

type onDemandPool struct {
	path   string
	engine http.Handler
}

func (f *onDemandPool) get() *ReactVm {
	return newReactVm(f.path, f.engine)
}

func (f onDemandPool) put(c *ReactVm) {
	c.Lock()
	c.FlushTimers()
	c.Gc(0)
	c.DestroyHeap()
}

func (f *onDemandPool) drop(c *ReactVm) {
	f.put(c)
}

type duktapePool struct {
	ch     chan *ReactVm
	path   string
	engine http.Handler
}

func (o *duktapePool) get() *ReactVm {
	return <-o.ch
}

func (o *duktapePool) put(ot *ReactVm) {
	// Drop any futured async calls
	ot.Lock()
	ot.FlushTimers()
	ot.Unlock()
	o.ch <- ot
}

func (o *duktapePool) drop(ot *ReactVm) {
	ot.Lock()
	ot.FlushTimers()
	ot.Gc(0)
	ot.DestroyHeap()
	ot = nil
	o.ch <- newReactVm(o.path, o.engine)
}
