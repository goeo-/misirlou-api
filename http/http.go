// Package http contains the basic framework on top of which the Misirlou API
// is built.
package http

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/erikdubbelboer/fasthttp"
	"github.com/go-redis/redis"
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

// PUT registers a handler for a PUT request.
func PUT(path string, handler func(c *Context)) {
	handlers = append(handlers, reqHandler{"PUT", path, handler})
}

// Options is a struct which is embedded in every context and contains
// information passed to the handlers directly by the main package.
type Options struct {
	DB    *models.DB
	Redis *redis.Client

	OAuth2ClientID     string
	OAuth2ClientSecret string

	BaseURL        string
	StoreTokensURL string
	HTTPS          bool
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

	sessCached bool
	sess       *models.Session
	ctx        *fasthttp.RequestCtx
	ip         net.IP
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

func (c *Context) reportError(err error) {
	fmt.Fprintln(os.Stderr, err)
}

// ResponseError can be passed to Error, and instead of returning a 500, it will
// return a response with the code and the message specified.
type ResponseError struct {
	Code    int
	Message string
}

func (re *ResponseError) Error() string {
	return re.Message
}

// Error closes the request with a 500 code and prints the error to stderr.
func (c *Context) Error(err error) {
	if re, ok := err.(*ResponseError); ok {
		c.SetCode(re.Code)
		c.WriteString(re.Message + "\n")
		return
	}
	c.SetCode(500)
	c.WriteString("Internal Server Error\n")
	c.reportError(err)
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

var newline = []byte{'\n'}

// SetJSONWithCode sets the response body to the given JSON value, as well
// as the HTTP response code.
func (c *Context) SetJSONWithCode(v interface{}, code int) {
	c.ctx.Response.ResetBody()
	c.ctx.SetContentType("application/json; charset=utf-8")
	c.SetCode(code)
	w := c.ctx.Response.BodyWriter()
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		c.Error(err)
	}
	c.ctx.Response.AppendBody(newline)
}

// SetJSON sets the response body to the given JSON value. The second value
// defines whether the response code should be 404 of 200. It's useful to
// quickly set the code to 404 when the value is null: c.SetJSON(t, t == nil)
func (c *Context) SetJSON(v interface{}, is404 bool) {
	code := 200
	if is404 {
		code = 404
	}
	c.SetJSONWithCode(v, code)
}

// SetCookie sets a cookie in the user's browser.
func (c *Context) SetCookie(k, v string, expire time.Duration) {
	cookie := fasthttp.AcquireCookie()
	cookie.SetKey(k)
	cookie.SetValue(v)
	cookie.SetExpire(time.Now().Add(expire))
	cookie.SetSecure(c.HTTPS)
	c.ctx.Response.Header.SetCookie(cookie)
	fasthttp.ReleaseCookie(cookie)
}

// DeleteCookie removes a cookie from the client.
func (c *Context) DeleteCookie(k string) {
	c.ctx.Response.Header.DelClientCookie(k)
}

// Redirect redirects the client to the given page, using the given response code.
func (c *Context) Redirect(code int, location string) {
	c.ctx.Redirect(location, code)
}

// Query retrieves a value from the query string.
func (c *Context) Query(s string) string {
	return string(c.ctx.QueryArgs().PeekBytes(s2b(s)))
}

// QueryInt retrieves a value from the querystring, and parses it as an int.
func (c *Context) QueryInt(s string) int {
	i, _ := strconv.Atoi(c.Query(s))
	return i
}

// QueryID retrieves a value from the querystring, and parses it as an ID.
func (c *Context) QueryID(s string) (i models.ID) {
	i.UnmarshalText(c.ctx.QueryArgs().PeekBytes(s2b(s)))
	return i
}

// ParamInt takes a named parameter set in the route of the request, and parses
// it as an int.
func (c *Context) ParamInt(s string) int {
	i, _ := strconv.Atoi(c.ctx.UserValue(s).(string))
	return i
}

// ParamID takes a named parameter set in the route of the request, and parses
// it as an ID.
func (c *Context) ParamID(s string) (i models.ID) {
	i.UnmarshalText([]byte(c.ctx.UserValue(s).(string)))
	return i
}

// JSON unmarshals the request's body into v, and returns any error.
func (c *Context) JSON(v interface{}) error {
	err := json.Unmarshal(c.ctx.Request.Body(), v)
	if err != nil {
		c.SetCode(400)
		c.WriteString("Bad JSON: " + err.Error())
	}
	return err
}

// Cookie returns the requested cookie.
func (c *Context) Cookie(s string) string {
	return string(c.ctx.Request.Header.Cookie(s))
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

// Session retrieves the Session related to this context.
func (c *Context) Session() *models.Session {
	if c.sessCached {
		return c.sess
	}
	c.sessCached = true
	// hash the authorization header, which contains the key
	hash := sha256.Sum256(c.ctx.Request.Header.Peek("Authorization"))
	encoded := hex.EncodeToString(hash[:])
	// find the session in the db
	sess, err := c.DB.Session(encoded)
	if err != nil {
		c.reportError(err)
		return nil
	}
	// cache it
	c.sess = sess
	return sess
}
