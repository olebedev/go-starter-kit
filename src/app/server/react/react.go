package react

import (
	"app/server/data"
	"app/server/utils"
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

func Bind(kit *utils.Kit) {
	r := react{kit: kit}
	r.init()
	kit.Engine.NoRoute(r.handle)
}

type react struct {
	pool pool
	kit  *utils.Kit
}

func (r *react) init() {
	if !r.kit.Conf.UBool("debug") {
		r.pool = newDuktapePool(runtime.NumCPU(), r.kit.Engine)
	} else {
		// Use onDemandPool to load full react
		// app each time for any http requests.
		// Useful to debug the app.
		r.pool = &onDemandPool{r.kit.Engine}
	}
}

// Handle handles all HTTP requests which
// have no been caught via static file
// handler or other middlewares
func (r *react) handle(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			UUID := c.MustGet("uuid").(*uuid.UUID)
			c.Writer.WriteHeader(http.StatusInternalServerError)
			c.Writer.Header().Add("Content-Type", "text/plain")
			c.Writer.Write([]byte(fmt.Sprintf("req uuid: %s\n%#v", UUID, r)))
			c.Abort()
		}
	}()

	vm := r.pool.get()
	vm.Lock()

	vm.PushGlobalObject()
	vm.GetPropString(-1, "__router__")
	// vm.Replace(-2)
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
		// because async calls is possible. So, we cannot
		// allow to break the context stack.
		vm.Lock()
		// Clean duktape vm stack
		vm.PopN(vm.GetTop())

		// Drop any futured async calls
		vm.ResetTimers()
		// Release the context
		vm.Unlock()
		// Return vm back to the pool
		r.pool.put(vm)

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
	case <-time.After(5 * time.Second):
		r.pool.drop(vm)
		UUID := c.MustGet("uuid").(*uuid.UUID)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.Header().Add("Content-Type", "text/plain")
		c.Writer.Write([]byte(fmt.Sprintf("req uuid: %s\ntime is out", UUID)))
		c.Abort()
	}
}

type resp struct {
	Error    string `json:"error"`
	Redirect string `json:"redirect"`
	Body     string `json:"body"`
}

// Interface to serve React app on demand or from prepared pool
type pool interface {
	get() *duktape.Context
	put(*duktape.Context)
	drop(*duktape.Context)
}

func newDuktapePool(size int, engine *gin.Engine) *duktapePool {
	pool := &duktapePool{
		ch:     make(chan *duktape.Context, size),
		engine: engine,
	}

	go func() {
		for i := 0; i < size; i++ {
			pool.ch <- newDuktapeContext(engine)
		}
	}()

	return pool
}

// Loads bundle.js to context
func newDuktapeContext(engine *gin.Engine) *duktape.Context {
	vm := duktape.New()
	vm.PevalString(`var console = {log:print,warn:print,error:print,info:print}`)
	fetch.Define(vm, engine)
	app, err := data.Asset("static/build/bundle.js")
	utils.Must(err)
	fmt.Println("static/build/bundle.js loaded")
	if err := vm.PevalString(string(app)); err != nil {
		derr := err.(*duktape.Error)
		fmt.Printf("\n\n\n%v\n%v\n\n\n", derr.FileName, derr.LineNumber)
		panic(derr.Message)
	}
	vm.PopN(vm.GetTop())
	return vm
}

type onDemandPool struct {
	engine *gin.Engine
}

func (f *onDemandPool) get() *duktape.Context {
	return newDuktapeContext(f.engine)
}

func (_ onDemandPool) put(c *duktape.Context) {
	c.Gc(0)
	c.DestroyHeap()
}

func (on *onDemandPool) drop(c *duktape.Context) {
	on.put(c)
}

type duktapePool struct {
	ch     chan *duktape.Context
	engine *gin.Engine
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
	o.ch <- newDuktapeContext(o.engine)
}
