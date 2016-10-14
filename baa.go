package baa

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

const (
	// DEV mode
	DEV = "development"
	// PROD mode
	PROD = "production"
	// TEST mode
	TEST = "test"
)

// Env default application runtime environment
var Env string

// Baa provlider an application
type Baa struct {
	debug           bool
	name            string
	di              DIer
	router          Router
	pool            sync.Pool
	errorHandler    ErrorHandleFunc
	notFoundHandler HandlerFunc
	middleware      []HandlerFunc
}

// Middleware middleware handler
type Middleware interface{}

// HandlerFunc context handler func
type HandlerFunc func(*Context)

// ErrorHandleFunc HTTP error handleFunc
type ErrorHandleFunc func(error, *Context)

// appInstances storage application instances
var appInstances map[string]*Baa

// defaultAppName default application name
const defaultAppName = "_default_"

// New create a baa application without any config.
func New() *Baa {
	b := new(Baa)
	b.middleware = make([]HandlerFunc, 0)
	b.pool = sync.Pool{
		New: func() interface{} {
			return NewContext(nil, nil, b)
		},
	}
	if Env != PROD {
		b.debug = true
	}
	b.SetDIer(NewDI())
	b.SetDI("router", NewTree(b))
	b.SetDI("logger", log.New(os.Stderr, "[Baa] ", log.LstdFlags))
	b.SetDI("render", newRender())
	b.SetNotFound(b.DefaultNotFoundHandler)
	return b
}

// Instance register or returns named application
func Instance(name string) *Baa {
	if name == "" {
		name = defaultAppName
	}
	if appInstances[name] == nil {
		appInstances[name] = New()
		appInstances[name].name = defaultAppName
	}
	return appInstances[name]
}

// Default initial a default app then returns
func Default() *Baa {
	return Instance(defaultAppName)
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
	b.Logger().Printf("Run mode: %s", Env)
	if len(files) == 0 {
		b.Logger().Printf("Listen %s", s.Addr)
		b.Logger().Fatal(s.ListenAndServe())
	} else if len(files) == 2 {
		b.Logger().Printf("Listen %s with TLS", s.Addr)
		b.Logger().Fatal(s.ListenAndServeTLS(files[0], files[1]))
	} else {
		panic("invalid TLS configuration")
	}
}

func (b *Baa) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := b.pool.Get().(*Context)
	c.Reset(w, r)

	// build handler chain
	h := b.Router().Match(r.Method, r.URL.Path, c)
	// notFound
	if h == nil {
		c.handlers = append(c.handlers, b.notFoundHandler)
	} else {
		c.handlers = append(c.handlers, h...)
	}

	c.Next()

	b.pool.Put(c)
}

// SetDIer set baa di
func (b *Baa) SetDIer(v DIer) {
	b.di = v
}

// SetDebug set baa debug
func (b *Baa) SetDebug(v bool) {
	b.debug = v
}

// Debug returns baa debug state
func (b *Baa) Debug() bool {
	return b.debug
}

// Logger return baa logger
func (b *Baa) Logger() Logger {
	return b.GetDI("logger").(Logger)
}

// Render return baa render
func (b *Baa) Render() Renderer {
	return b.GetDI("render").(Renderer)
}

// Router return baa router
func (b *Baa) Router() Router {
	if b.router == nil {
		b.router = b.GetDI("router").(Router)
	}
	return b.router
}

// Use registers a middleware
func (b *Baa) Use(m ...Middleware) {
	for i := range m {
		if m[i] != nil {
			b.middleware = append(b.middleware, wrapMiddleware(m[i]))
		}
	}
}

// SetDI registers a dependency injection
func (b *Baa) SetDI(name string, h interface{}) {
	switch name {
	case "logger":
		if _, ok := h.(Logger); !ok {
			panic("DI logger must be implement interface baa.Logger")
		}
	case "render":
		if _, ok := h.(Renderer); !ok {
			panic("DI render must be implement interface baa.Renderer")
		}
	case "router":
		if _, ok := h.(Router); !ok {
			panic("DI router must be implement interface baa.Router")
		}
	}
	b.di.Set(name, h)
}

// GetDI fetch a registered dependency injection
func (b *Baa) GetDI(name string) interface{} {
	return b.di.Get(name)
}

// Static set static file route
// h used for set Expries ...
func (b *Baa) Static(prefix string, dir string, index bool, h HandlerFunc) {
	if prefix == "" {
		panic("baa.Static prefix can not be empty")
	}
	if dir == "" {
		panic("baa.Static dir can not be empty")
	}
	staticHandler := newStatic(prefix, dir, index, h)
	b.Get(prefix, staticHandler)
	b.Get(prefix+"*", staticHandler)
}

// SetAutoHead sets the value who determines whether add HEAD method automatically
// when GET method is added. Combo router will not be affected by this value.
func (b *Baa) SetAutoHead(v bool) {
	b.Router().SetAutoHead(v)
}

// SetAutoTrailingSlash optional trailing slash.
func (b *Baa) SetAutoTrailingSlash(v bool) {
	b.Router().SetAutoTrailingSlash(v)
}

// Route is a shortcut for same handlers but different HTTP methods.
//
// Example:
// 		baa.Route("/", "GET,POST", h)
func (b *Baa) Route(pattern, methods string, h ...HandlerFunc) RouteNode {
	var ru RouteNode
	var ms []string
	if methods == "*" {
		for m := range RouterMethods {
			ms = append(ms, m)
		}
	} else {
		ms = strings.Split(methods, ",")
	}
	for _, m := range ms {
		ru = b.Router().Add(strings.TrimSpace(m), pattern, h)
	}
	return ru
}

// Group registers a list of same prefix route
func (b *Baa) Group(pattern string, f func(), h ...HandlerFunc) {
	b.Router().GroupAdd(pattern, f, h)
}

// Any is a shortcut for b.Router().handle("*", pattern, handlers)
func (b *Baa) Any(pattern string, h ...HandlerFunc) RouteNode {
	var ru RouteNode
	for m := range RouterMethods {
		ru = b.Router().Add(m, pattern, h)
	}
	return ru
}

// Delete is a shortcut for b.Route(pattern, "DELETE", handlers)
func (b *Baa) Delete(pattern string, h ...HandlerFunc) RouteNode {
	return b.Router().Add("DELETE", pattern, h)
}

// Get is a shortcut for b.Route(pattern, "GET", handlers)
func (b *Baa) Get(pattern string, h ...HandlerFunc) RouteNode {
	return b.Router().Add("GET", pattern, h)
}

// Head is a shortcut forb.Route(pattern, "Head", handlers)
func (b *Baa) Head(pattern string, h ...HandlerFunc) RouteNode {
	return b.Router().Add("HEAD", pattern, h)
}

// Options is a shortcut for b.Route(pattern, "Options", handlers)
func (b *Baa) Options(pattern string, h ...HandlerFunc) RouteNode {
	return b.Router().Add("OPTIONS", pattern, h)
}

// Patch is a shortcut for b.Route(pattern, "PATCH", handlers)
func (b *Baa) Patch(pattern string, h ...HandlerFunc) RouteNode {
	return b.Router().Add("PATCH", pattern, h)
}

// Post is a shortcut for b.Route(pattern, "POST", handlers)
func (b *Baa) Post(pattern string, h ...HandlerFunc) RouteNode {
	return b.Router().Add("POST", pattern, h)
}

// Put is a shortcut for b.Route(pattern, "Put", handlers)
func (b *Baa) Put(pattern string, h ...HandlerFunc) RouteNode {
	return b.Router().Add("PUT", pattern, h)
}

// SetNotFound set not found route handler
func (b *Baa) SetNotFound(h HandlerFunc) {
	b.notFoundHandler = h
}

// NotFound execute not found handler
func (b *Baa) NotFound(c *Context) {
	if b.notFoundHandler != nil {
		b.notFoundHandler(c)
		return
	}
	http.NotFound(c.Resp, c.Req)
}

// SetError set error handler
func (b *Baa) SetError(h ErrorHandleFunc) {
	b.errorHandler = h
}

// Error execute internal error handler
func (b *Baa) Error(err error, c *Context) {
	if err == nil {
		err = errors.New("Internal Server Error")
	}
	if b.errorHandler != nil {
		b.errorHandler(err, c)
		return
	}
	code := http.StatusInternalServerError
	msg := http.StatusText(code)
	if b.debug {
		msg = err.Error()
	}
	b.Logger().Println(err)
	http.Error(c.Resp, msg, code)
}

// DefaultNotFoundHandler invokes the default HTTP error handler.
func (b *Baa) DefaultNotFoundHandler(c *Context) {
	code := http.StatusNotFound
	msg := http.StatusText(code)
	http.Error(c.Resp, msg, code)
}

// URLFor use named route return format url
func (b *Baa) URLFor(name string, args ...interface{}) string {
	return b.Router().URLFor(name, args...)
}

// wrapMiddleware wraps middleware.
func wrapMiddleware(m Middleware) HandlerFunc {
	switch m := m.(type) {
	case HandlerFunc:
		return m
	case func(*Context):
		return m
	case http.Handler, http.HandlerFunc:
		return WrapHandlerFunc(func(c *Context) {
			m.(http.Handler).ServeHTTP(c.Resp, c.Req)
		})
	case func(http.ResponseWriter, *http.Request):
		return WrapHandlerFunc(func(c *Context) {
			m(c.Resp, c.Req)
		})
	default:
		panic("unknown middleware")
	}
}

// WrapHandlerFunc wrap for context handler chain
func WrapHandlerFunc(h HandlerFunc) HandlerFunc {
	return func(c *Context) {
		h(c)
		c.Next()
	}
}

func init() {
	appInstances = make(map[string]*Baa)
	Env = os.Getenv("BAA_ENV")
	if Env == "" {
		Env = DEV
	}
}
