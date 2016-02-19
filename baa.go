// Package baa provider a fast & simple Go web framework, routing, middleware, dependency injection, http context.
/*
app := baa.Classic()
app.Get("/", function(c *baa.Context) {
    c.String(200, "Hello World!")
})
app.Run(":8001")
*/
package baa

import (
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

const (
	// MiddlewareMaxSize set size limit of middleware
	MiddlewareMaxSize = 20
	// RouteMaxLength set length limit of route pattern
	RouteMaxLength = 256
	// RouterParamMaxLength set length limit of route pattern param
	RouterParamMaxLength = 64

	// DEV mode
	DEV = "development"
	// PROD mode
	PROD = "production"
	// TEST mode
	TEST = "test"
)

// Baa provlider an application
type Baa struct {
	debug        bool
	name         string
	di           *DI
	router       *Router
	logger       Logger
	render       Renderer
	pool         sync.Pool
	errorHandler ErrorHandleFunc
	middleware   []HandlerFunc
}

// HandlerFunc context handler
type HandlerFunc func(*Context)

// ErrorHandleFunc HTTP error handleFunc
type ErrorHandleFunc func(error, *Context)

// HTTPError represents an error that occured while handling a request.
type HTTPError struct {
	code    int
	message string
}

// default application for baa
var _defaultAPP *Baa

func init() {
	_defaultAPP = Classic()
}

// App return the default baa instance
func App() *Baa {
	return _defaultAPP
}

// Classic create a baa application with default config.
func Classic() *Baa {
	b := New()
	b.SetRender(NewRender())
	b.SetErrorHandler(b.DefaultErrorHandler)
	return b
}

// New create a baa application without any config.
func New() *Baa {
	b := new(Baa)
	b.middleware = make([]HandlerFunc, 0)
	b.pool.New = func() interface{} {
		return NewContext(nil, nil, b)
	}
	b.SetLogger(log.New(os.Stderr, "[Baa] ", log.LstdFlags))
	b.SetDIer(NewDI())
	b.SetRouter(NewRouter())
	return b
}

// Server returns the internal *http.Server.
func (b *Baa) Server(addr string) *http.Server {
	s := &http.Server{Addr: addr}
	return s
}

// Run runs a server.
func (b *Baa) Run(addr string) {
	b.run(b.Server(addr))
}

// RunTLS runs a server with TLS configuration.
func (b *Baa) RunTLS(addr, certfile, keyfile string) {
	b.run(b.Server(addr), certfile, keyfile)
}

// RunServer runs a custom server.
func (b *Baa) RunServer(s *http.Server) {
	b.run(s)
}

// RunTLSServer runs a custom server with TLS configuration.
func (b *Baa) RunTLSServer(s *http.Server, crtFile, keyFile string) {
	b.run(s, crtFile, keyFile)
}

func (b *Baa) run(s *http.Server, files ...string) {
	s.Handler = b
	if len(files) == 0 {
		b.logger.Printf("Listen %s", s.Addr)
		b.logger.Fatal(s.ListenAndServe())
	} else if len(files) == 2 {
		b.logger.Printf("Listen %s with TLS", s.Addr)
		b.logger.Fatal(s.ListenAndServeTLS(files[0], files[1]))
	} else {
		b.logger.Fatal("invalid TLS configuration")
	}
}

func (b *Baa) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := b.pool.Get().(*Context)
	defer b.pool.Put(c)
	c.reset(w, r)

	var h HandlerFunc
	route := b.router.Match(r.Method, r.URL.Path, c)
	if route == nil {
		// notFound
		h = b.router.GetNotFoundHandler()
		if h == nil {
			h = func(c *Context) {
				http.NotFound(c.Resp, c.Req)
			}
		}
	} else {
		h = route.handle
	}

	// concat middleware handle and route handle to a list
	c.mw = append(c.mw, b.middleware...)
	c.mw = append(c.mw, h)
	c.Next()
}

// SetDebug set baa debug
func (b *Baa) SetDebug(v bool) {
	b.debug = v
}

// SetDIer registers a Baa.DI
func (b *Baa) SetDIer(di *DI) {
	b.di = di
}

// SetLogger registers a Baa.Logger
func (b *Baa) SetLogger(logger Logger) {
	b.logger = logger
}

// Logger return baa logger
func (b *Baa) Logger() Logger {
	return b.logger
}

// SetRender registers a Baa.Renderer
func (b *Baa) SetRender(r Renderer) {
	b.render = r
}

// SetRouter registers a Baa.Router
func (b *Baa) SetRouter(r *Router) {
	b.router = r
}

// SetErrorHandler registers a custom Baa.ErrorHandleFunc.
func (b *Baa) SetErrorHandler(h ErrorHandleFunc) {
	b.errorHandler = h
}

// DefaultErrorHandler invokes the default HTTP error handler.
func (b *Baa) DefaultErrorHandler(err error, c *Context) {
	code := http.StatusInternalServerError
	msg := http.StatusText(code)
	if he, ok := err.(*HTTPError); ok {
		code = he.code
		msg = he.message
	}
	if b.debug {
		msg = err.Error()
	}
	http.Error(c.Resp, msg, code)
}

// Use registers a middleware
func (b *Baa) Use(m HandlerFunc) {
	if len(b.middleware) > MiddlewareMaxSize {
		b.logger.Printf("middleware num overhead, limit %d\n", MiddlewareMaxSize)
		return
	}
	b.middleware = append(b.middleware, m)
}

// SetDI registers a dependency injection
func (b *Baa) SetDI(name string, h interface{}) {
	b.di.Set(name, h)
}

// GetDI fetch a registered dependency injection
func (b *Baa) GetDI(name string) interface{} {
	return b.di.Get(name)
}

// SetAutoHead sets the value who determines whether add HEAD method automatically
// when GET method is added. Combo router will not be affected by this value.
func (b *Baa) SetAutoHead(v bool) {
	b.router.autoHead = v
}

// Route is a shortcut for same handlers but different HTTP methods.
//
// Example:
// 		baa.route("/", "GET,POST", h)
func (b *Baa) Route(pattern, methods string, h ...HandlerFunc) *Route {
	var rs *Route
	for _, m := range strings.Split(methods, ",") {
		rs = b.router.add(strings.TrimSpace(m), pattern, h)
	}
	return rs
}

// Group registers a list of same prefix route
func (b *Baa) Group(pattern string, fn func(), h ...HandlerFunc) {

}

// Get is a shortcut for b.router.handle("GET", pattern, handlers)
func (b *Baa) Get(pattern string, h ...HandlerFunc) *Route {
	rs := b.router.add("GET", pattern, h)
	if b.router.autoHead {
		b.Head(pattern, h...)
	}
	return rs
}

// Patch is a shortcut for b.router.handle("PATCH", pattern, handlers)
func (b *Baa) Patch(pattern string, h ...HandlerFunc) *Route {
	return b.router.add("PATCH", pattern, h)
}

// Post is a shortcut for b.router.handle("POST", pattern, handlers)
func (b *Baa) Post(pattern string, h ...HandlerFunc) *Route {
	return b.router.add("POST", pattern, h)
}

// Put is a shortcut for b.router.handle("PUT", pattern, handlers)
func (b *Baa) Put(pattern string, h ...HandlerFunc) *Route {
	return b.router.add("PUT", pattern, h)
}

// Delete is a shortcut for b.router.handle("DELETE", pattern, handlers)
func (b *Baa) Delete(pattern string, h ...HandlerFunc) *Route {
	return b.router.add("DELETE", pattern, h)
}

// Options is a shortcut for b.router.handle("OPTIONS", pattern, handlers)
func (b *Baa) Options(pattern string, h ...HandlerFunc) *Route {
	return b.router.add("OPTIONS", pattern, h)
}

// Head is a shortcut for b.router.handle("HEAD", pattern, handlers)
func (b *Baa) Head(pattern string, h ...HandlerFunc) *Route {
	return b.router.add("HEAD", pattern, h)
}

// Any is a shortcut for b.router.handle("*", pattern, handlers)
func (b *Baa) Any(pattern string, h ...HandlerFunc) *Route {
	return b.router.add("*", pattern, h)
}

// NotFound set 404 router
func (b *Baa) NotFound(h HandlerFunc) {
	b.router.NotFound(h)
}

// URLFor use named route return format url
func (b *Baa) URLFor(name string, args ...interface{}) string {
	return b.router.URLFor(name, args...)
}

// NewHTTPError creates a new HTTPError instance.
func NewHTTPError(code int, msg ...string) *HTTPError {
	e := &HTTPError{code: code, message: http.StatusText(code)}
	if len(msg) > 0 {
		m := msg[0]
		e.message = m
	}
	return e
}

// SetCode sets code.
func (e *HTTPError) SetCode(code int) {
	e.code = code
}

// Code returns code.
func (e *HTTPError) Code() int {
	return e.code
}

// Error returns message.
func (e *HTTPError) Error() string {
	return e.message
}
