// Package http contains the basic framework on top of which the Misirlou API
// is built.
package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/erikdubbelboer/fasthttp"
	"github.com/thehowl/fasthttprouter"
	"zxq.co/ripple/misirlou-api/models"
)

type reqHandler struct {
	method  string
	path    string
	handler func(c *Context)
}

var handlers = make([]reqHandler, 0, 16)

// GET registers a handler for a GET request.
func GET(path string, handler func(c *Context)) {
	handlers = append(handlers, reqHandler{"GET", path, handler})
}

// POST registers a handler for a POST request.
func POST(path string, handler func(c *Context)) {
	handlers = append(handlers, reqHandler{"POST", path, handler})
}

// Options is a struct which is embedded in every context and contains
// information passed to the handlers directly by the main package.
type Options struct {
	DB *models.DB
}

// Handler creates an HTTP request handler using httprouter.
func Handler(o Options) fasthttp.RequestHandler {
	r := fasthttprouter.New()

	for _, h := range handlers {
		r.Handle(h.method, h.path, wrapper(o, h.handler))
	}

	return r.Handler
}

func wrapper(o Options, f func(*Context)) fasthttp.RequestHandler {
	return func(r *fasthttp.RequestCtx) {
		start := time.Now()
		ctx := &Context{
			Options: o,
			ctx:     r,
		}
		ctx.ctx.SetContentType("text/plain; charset=utf-8")
		defer func() {
			err := recover()
			if err != nil {
				fmt.Println("RECOVERED FROM PANIC:", err)
				debug.PrintStack()
				ctx.SetCode(500)
			}
			code := ctx.ctx.Response.StatusCode()

			elapsed := time.Since(start)
			// successful request, or not our problem anyways
			color := "42"
			if code >= 500 && code < 600 {
				// fatal error
				color = "43"
			}

			fmt.Printf(
				"%s | %15s |\x1b[%sm %3d \x1b[0m %-7s %s\n",
				time.Now().Format("2006-01-02 15:04:05"),
				elapsed, color, code, r.Method(), r.URI().Path(),
			)
		}()

		f(ctx)
	}
}

// Context is the information passed to each request about the requested method.
type Context struct {
	Options

	ctx *fasthttp.RequestCtx
	ip  net.IP
}

// Header retrieves a header from the request.
func (c *Context) Header(s string) []byte {
	return c.ctx.Request.Header.PeekBytes(s2b(s))
}

// SetHeader sets an header in the response.
func (c *Context) SetHeader(k, v string) {
	c.ctx.Response.Header.Set(k, v)
}

// SetCode sets the status code for the response, and sends all the header to
// the response.
func (c *Context) SetCode(i int) {
	c.ctx.SetStatusCode(i)
}

// Error closes the request with a 500 code and prints the error to stderr.
func (c *Context) Error(err error) {
	c.SetCode(500)
	c.WriteString("Internal Server Error")
	fmt.Fprintln(os.Stderr, err)
}

// WriteString writes s to the response. We provide WriteString and not Write
// because c.ctx.Write really makes a b2s conversion and then calls
// AppendBodyString. To avoid this indirection, we do not provide Write, and
// instead provide a WriteString with no return arguments - that is because even
// if there was (int, error), they would always be len(s) and nil.
func (c *Context) WriteString(s string) {
	c.ctx.Response.AppendBodyString(s)
}

// SetBody changes the existing response body with the passed body.
func (c *Context) SetBody(b []byte) {
	c.ctx.Response.SetBody(b)
}

// SetJSON sets the response body to the given JSON value.
func (c *Context) SetJSON(code int, v interface{}) {
	c.ctx.Response.ResetBody()
	c.ctx.SetContentType("application/json; charset=utf-8")
	c.SetCode(code)
	err := json.NewEncoder(c.ctx.Response.BodyWriter()).Encode(v)
	if err != nil {
		c.Error(err)
	}
}

// Query retrieves a value from the query string.
func (c *Context) Query(s string) []byte {
	return c.ctx.QueryArgs().PeekBytes(s2b(s))
}

// QueryInt retrieves a value from the query int, and parses it as an int.
func (c *Context) QueryInt(s string) int {
	i, _ := strconv.Atoi(b2s(c.Query(s)))
	return i
}

var ipHeaders = [...][]byte{
	[]byte("X-Forwarded-For"),
	[]byte("X-Real-IP"),
}

// IP retrieves the IP address of the request. If the remote addr is loopback
// or invalid (e.g. an unix socket), then it is allowed to change the IP address
// by using the request header X-Forwarded-For or X-Real-IP.
func (c *Context) IP() net.IP {
	if len(c.ip) != 0 {
		return c.ip
	}

	ip := c.ctx.RemoteIP()
	// if it is zero or the loopback, it means that it probably is the same
	// computer that called this in the first place, so we can allow them
	// to set the IP address using HTTP headers (X-Forwarded-For).
	// TODO(howl): we should probably also handle 192.168.*.* for the local
	// network, although I'm not too sure on how we should get around that,
	// so I'll only allow unix sockets and loopback for the moment.
	if !(ip.Equal(net.IPv4zero) || ip.IsLoopback()) {
		c.ip = ip
		return ip
	}

	// check if there is an IP in any of the headers we know could contain them.
	for _, hKey := range ipHeaders {
		h := c.ctx.Request.Header.PeekBytes(hKey)
		if len(h) > 0 {
			pos := bytes.LastIndexByte(h, ',')
			if pos >= 0 {
				c.ip = net.ParseIP(b2s(h[pos+1:]))
			} else {
				c.ip = net.ParseIP(b2s(h))
			}
			return c.ip
		}
	}

	c.ip = ip
	return ip
}
