package react

import (
	"app/server/data"
	. "app/server/utils"
	"fmt"

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
	var v string
	vm := r.pool.get()
	vm.PevalString(`React.renderToString(React.createElement(App, {}));`)
	v = vm.SafeToString(-1)
	vm.PopN(vm.GetTop())
	r.pool.put(vm)

	c.Writer.WriteHeader(200)
	c.Writer.Header().Add("Content-Type", "text/html")
	c.Writer.Write([]byte("<!doctype html>\n" + v))
	c.Abort()
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
		panic(err.(*duktape.Error).Message)
	}
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
