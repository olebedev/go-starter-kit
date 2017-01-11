package fetch

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/parnurzeal/gorequest"
	"gopkg.in/olebedev/go-duktape.v2"
)

var bundle string

func init() {
	b, err := Asset("dist/bundle.js")
	must(err)
	bundle = string(b)
}

func PushGlobal(c *duktape.Context, server http.Handler) {
	c.PushTimers()
	must(c.PevalString(bundle))
	c.Pop()

	c.PushGlobalObject()
	c.GetPropString(-1, "fetch")
	c.PushGoFunction(goFetchSync(server))
	c.PutPropString(-2, "goFetchSync")
	c.Pop2()

}

type options struct {
	Url     string      `json:"url"`
	Method  string      `json:"method"`
	Headers http.Header `json:"headers"`
	Body    string      `json:"body"`
}

func goFetchSync(server http.Handler) func(*duktape.Context) int {
	return func(c *duktape.Context) int {
		url := c.SafeToString(0)
		opts := options{
			Method:  gorequest.GET,
			Headers: http.Header{},
		}
		err := json.Unmarshal([]byte(c.JsonEncode(1)), &opts)
		must(err)

		var resp response
		if strings.HasPrefix(url, "http") || strings.HasPrefix(url, "//") {
			resp = fetchHttp(url, opts)
		} else if strings.HasPrefix(url, "/") {
			resp = fetchHandlerFunc(server, url, opts)
		} else {
			return duktape.ErrRetURI
		}

		j, err := json.MarshalIndent(resp, "", "  ")
		must(err)
		c.Pop3()
		c.PushString(string(j))
		c.JsonDecode(-1)
		return 1
	}
}

type response struct {
	options
	Status     int     `json:"status"`
	StatusText string  `json:"statusText,omitempty"`
	Errors     []error `json:"errors"`
}

func fetchHttp(url string, opts options) response {
	var resp gorequest.Response
	var body string
	var errs []error
	client := gorequest.New()
	switch opts.Method {
	case gorequest.HEAD:
		resp, body, errs = client.Head(url).End()
	case gorequest.GET:
		resp, body, errs = client.Get(url).End()
	case gorequest.POST:
		resp, body, errs = client.Post(url).Query(opts.Body).End()
	case gorequest.PUT:
		resp, body, errs = client.Put(url).Query(opts.Body).End()
	case gorequest.PATCH:
		resp, body, errs = client.Patch(url).Query(opts.Body).End()
	case gorequest.DELETE:
		resp, body, errs = client.Delete(url).End()
	}

	result := response{
		options:    opts,
		Status:     resp.StatusCode,
		StatusText: resp.Status,
		Errors:     errs,
	}
	result.Body = body
	result.Headers = resp.Header
	return result
}

func fetchHandlerFunc(server http.Handler, url string, opts options) response {
	result := response{
		options: opts,
		Errors:  []error{},
	}

	if server == nil {
		result.Errors = append(result.Errors, errors.New("`http.Handler` isn't set yet"))
		result.Status = http.StatusInternalServerError
	}

	b := bytes.NewBufferString(opts.Body)
	res := httptest.NewRecorder()
	req, err := http.NewRequest(opts.Method, url, b)

	if err != nil {
		result.Errors = []error{err}
		return result
	}

	req.Header = opts.Headers
	server.ServeHTTP(res, req)
	result.Status = res.Code
	result.Headers = res.Header()
	result.Body = res.Body.String()
	return result
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
