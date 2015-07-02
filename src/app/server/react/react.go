package react

import (
	"app/server/data"
	. "app/server/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olebedev/go-duktape"
	"github.com/olebedev/go-duktape-fetch"
)

func Bind(kit *Kit) {
	r := react{kit: kit}
	r.init()
	kit.Engine.NoRoute(r.handle)
}

type react struct {
	pool   pool
	engine *gin.Engine
	kit    *Kit
}

func (r *react) init() {
	if r.kit.Conf.UBool("duktape.pool.use") {
		r.pool = newDuktapePool(r.kit.Conf.UInt("duktape.pool.size", 1), r.engine)
	} else {
		r.pool = &onDemandPool{r.engine}
	}
}

func (r *react) handle(c *gin.Context) {
	vm := r.pool.get()

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
	ch := make(chan struct{}, 1)
	vm.PushGoFunction(func(ctx *duktape.Context) int {

		// Getting response object via json
		r := func() *resp {
			var re resp
			json.Unmarshal([]byte(vm.JsonEncode(-1)), &re)
			return &re
		}()

		// Handle the response
		if len(r.Redirect) == 0 && len(r.Error) == 0 {
			// If no redirection and no error
			c.Writer.WriteHeader(http.StatusOK)
			c.Writer.Header().Add("Content-Type", "text/html")
			c.Writer.Write([]byte("<!doctype html>\n" + r.Body))
			c.Abort()
			// If redirect
		} else if len(r.Redirect) != 0 {
			c.Redirect(http.StatusMovedPermanently, r.Redirect)
			// If internal error
		} else if len(r.Error) != 0 {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			c.Writer.Header().Add("Content-Type", "text/plain")
			c.Writer.Write([]byte(r.Error))
			c.Abort()
		}

		// Unlock handler
		ch <- struct{}{}
		// Return nothing into duktape context
		return 0
	})

	// Duktape stack -> [ {global}, __router__, "renderToString", {\"url\":\"...\"}, {url:...}, {func: true} ]
	vm.PcallProp(1, 2)
	// Lock handler and wait for app response
	<-ch
	// Clean stack
	if i := vm.GetTop(); i > 0 {
		vm.PopN(i)
	}
	// Return vm back to the pool
	r.pool.put(vm)
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
}

func newDuktapePool(size int, engine *gin.Engine) *duktapePool {
	pool := &duktapePool{
		ch: make(chan *duktape.Context, size),
	}
loop:
	for {
		select {
		case pool.ch <- newDuktapeContext(engine):
		default:
			break loop
		}
	}
	return pool
}

// Loads bundle.js to context
func newDuktapeContext(engine *gin.Engine) *duktape.Context {
	vm := duktape.New()
	vm.PevalString(`var console = {log:print,warn:print,error:print,info:print}`)
	fetch.Define(vm, engine)
	app, err := data.Asset("static/build/bundle.js")
	Must(err)
	fmt.Println("static/build/bundle.js loaded")
	if err := vm.PevalString(string(app)); err != nil {
		derr := err.(*duktape.Error)
		fmt.Printf("\n\n\n%v\n%v\n\n\n", derr.FileName, derr.LineNumber)
		panic(derr.Message)
	}
	vm.PopN(vm.GetTop())
	return vm
}

// Loads file pre request
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

// Prefetched pool
type duktapePool struct {
	ch chan *duktape.Context
}

func (o *duktapePool) get() *duktape.Context {
	return <-o.ch
}

func (o *duktapePool) put(ot *duktape.Context) {
	o.ch <- ot
}
