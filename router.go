package baa

import (
	//"net/http"
	"sync"
)

// METHODS 定义支持的HTTP method
var METHODS = map[string]bool{
	"GET":     true,
	"POST":    true,
	"PUT":     true,
	"DELETE":  true,
	"PATCH":   true,
	"OPTIONS": true,
	"HEAD":    true,
}

// Router provlider router for baa
type Router struct {
	mu              sync.RWMutex
	autoHead        bool
	routeMap        map[string]map[string]*Route
	routeNamedMap   map[string]string
	notFoundHandler HandlerFunc
	baa             *Baa
}

// Route is a single route
type Route struct {
	pattern  string
	handlers []HandlerFunc
	router   *Router
}

// NewRouter create a router instance
func NewRouter(b *Baa) *Router {
	r := new(Router)
	r.baa = b
	r.routeMap = make(map[string]map[string]*Route)
	for m := range METHODS {
		r.routeMap[m] = make(map[string]*Route)
	}
	r.routeNamedMap = make(map[string]string)
	return r
}

// Handle registers a new request handle with the given pattern, method and handlers.
func (r *Router) add(method string, pattern string, handlers []Handler) *Route {
	var ru *Route
	var ok bool
	if _, ok = METHODS[method]; !ok {
		panic("unsupport http method [" + method + "]")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if ru, ok = r.routeMap[method][pattern]; !ok {
		ru = &Route{
			pattern: pattern,
			router:  r,
		}
		ru.handlers = make([]HandlerFunc, 0, 1)
		r.routeMap[method][pattern] = ru
	}
	for _, h := range handlers {
		ru.handlers = append(ru.handlers, wrapHandler(h))
	}
	return ru
}

// NotFound set the route not match result.
// Configurable http.HandlerFunc which is called when no matching route is
// found. If it is not set, http.NotFound is used.
// Be sure to set 404 response code in your handler.
func (r *Router) NotFound(h Handler) {
	r.notFoundHandler = wrapHandler(h)
}

// GetNotFoundHandler ...
func (r *Router) GetNotFoundHandler() HandlerFunc {
	return r.notFoundHandler
}

// Match match the uri for handler
func (r *Router) Match(method, uri string) *Route {
	for p := range r.routeMap[method] {
		if p == uri {
			return r.routeMap[method][p]
		}
	}
	return nil
}

// ServeHTTP implements the Handler interface and can be registered to a HTTP server
// func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
// 	c := NewContext(w, req, r.baa)
// 	route := r.Match(req.Method, req.URL.Path)
// 	if route != nil {
// 		route.handle(c)
// 		return
// 	}

// 	// 404
// 	if r.notFoundHandler == nil {
// 		http.NotFound(w, req)
// 	} else {
// 		r.notFoundHandler(c)
// 	}
// }

// Name set name of route
func (r *Route) Name(name string) {
	if name == "" {
		return
	}
	r.router.routeNamedMap[name] = r.pattern
}

// handle ...
func (r *Route) handle(c *Context) error {
	for _, h := range r.handlers {
		err := h(c)
		if err != nil {
			return err
		}
	}
	return nil
}
