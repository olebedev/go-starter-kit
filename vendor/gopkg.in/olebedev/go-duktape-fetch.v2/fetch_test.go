package fetch

import (
	"encoding/json"
	"testing"

	. "gopkg.in/check.v1"

	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
	"gopkg.in/olebedev/go-duktape.v2"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type FetchSuite struct {
	ctx     *duktape.Context
	goFetch func(*duktape.Context) int
}

func (s *FetchSuite) SetUpSuite(c *C) {
	gin.SetMode(gin.ReleaseMode)
}

func (s *FetchSuite) SetUpTest(c *C) {
	s.goFetch = goFetchSync(nil)
	s.ctx = duktape.New()
	PushGlobal(s.ctx, nil)
}

var _ = Suite(&FetchSuite{})

func (s *FetchSuite) TestStackAroundGoFetchSync(c *C) {
	s.ctx.PushString("http://ya.ru")
	s.ctx.PushObject() // options
	s.ctx.PushObject() // headers => [ url options headers ]
	s.ctx.PutPropString(-2, "headers")

	s.ctx.PushContextDump()
	c.Assert(s.ctx.SafeToString(-1), Equals, "ctx: top=2, stack=[\"http://ya.ru\",{headers:{}}]")
	c.Assert(s.goFetch(s.ctx), Equals, 1)

	c.Assert(s.ctx.GetTop(), Equals, 1)

	resp := response{}
	c.Assert(json.Unmarshal([]byte(s.ctx.JsonEncode(-1)), &resp), IsNil)
}

func (s *FetchSuite) TestGoFetchSyncExternal(c *C) {

	s.ctx.PushString("http://ya.ru")
	s.ctx.PushObject()                 // options => [ url {} ]
	s.ctx.PushObject()                 // headers => [ url {} {} ]
	s.ctx.PutPropString(-2, "headers") // => [ url {headers: {}} ]

	c.Assert(s.goFetch(s.ctx), Equals, 1)

	resp := response{}
	c.Assert(json.Unmarshal([]byte(s.ctx.JsonEncode(-1)), &resp), IsNil)
	c.Assert(resp.Method, Equals, gorequest.GET)
	c.Assert(resp.Status, Equals, 200)
	c.Assert(resp.StatusText, Equals, "200 Ok")
	c.Assert(resp.Errors, HasLen, 0)
	c.Assert(resp.Body[:15], Equals, "<!DOCTYPE html>")
	c.Assert(resp.Headers.Get("Content-Type"), Equals, "text/html; charset=UTF-8")
}

func (s *FetchSuite) TestGoFetchInternal404(c *C) {
	PushGlobal(s.ctx, gin.Default())
	c.Assert(s.ctx.PevalString(`
		fetch.goFetchSync('/404', {});
	`), IsNil)
	s.ctx.JsonEncode(-1)
	respString := s.ctx.SafeToString(-1)
	resp := response{}
	json.Unmarshal([]byte(respString), &resp)
	c.Assert(resp.Status, Equals, 404)
	c.Assert(resp.Method, Equals, gorequest.GET)
	c.Assert(resp.Body, Equals, "404 page not found")
}

func (s *FetchSuite) TestGoFetchPromise(c *C) {
	PushGlobal(s.ctx, gin.Default())
	ch := make(chan string)
	s.ctx.PushGlobalGoFunction("cbk", func(co *duktape.Context) int {
		ch <- co.SafeToString(-1)
		return 0
	})

	c.Assert(s.ctx.PevalString(`
		fetch('/404')
			.then(function(resp){
				return resp.text();
			}).then(cbk);
	`), IsNil)
	c.Assert(<-ch, Equals, "404 page not found")
}

func (s *FetchSuite) TestGoFetchJson(c *C) {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(201, map[string]interface{}{
			"hello": "world",
		})
	})

	PushGlobal(s.ctx, r)

	ch := make(chan bool)
	s.ctx.PushGlobalGoFunction("cbk", func(co *duktape.Context) int {
		ch <- co.GetType(-1).IsObject()
		return 0
	})

	c.Assert(s.ctx.PevalString(`
		fetch('/')
			.then(function(resp){
				return resp.json();
			}).then(cbk);
	`), IsNil)
	c.Assert(<-ch, Equals, true)
}

func (s *FetchSuite) TestGlobals(c *C) {

	// fetch
	c.Assert(s.ctx.PevalString(`typeof fetch;`), IsNil)
	c.Assert(s.ctx.SafeToString(-1), Equals, "function")
	s.ctx.Pop()

	// fetch.goFetchSync
	c.Assert(s.ctx.PevalString(`typeof fetch.goFetchSync;`), IsNil)
	c.Assert(s.ctx.SafeToString(-1), Equals, "function")
	s.ctx.Pop()

	// fetch.Promise
	c.Assert(s.ctx.PevalString(`typeof fetch.Promise;`), IsNil)
	c.Assert(s.ctx.SafeToString(-1), Equals, "function")
	s.ctx.Pop()

}

func (s *FetchSuite) TearDownTest(c *C) {
	s.ctx.DestroyHeap()
}
