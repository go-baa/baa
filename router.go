package baa

import (
	"fmt"
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
func (r *Router) add(method string, pattern string, handlers []HandlerFunc) *Route {
	if pattern == "" {
		panic("route pattern can not be emtpy!")
	}
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
	ru.handlers = append(ru.handlers, handlers...)

	return ru
}

// NotFound set the route not match result.
// Configurable http.HandlerFunc which is called when no matching route is
// found. If it is not set, http.NotFound is used.
// Be sure to set 404 response code in your handler.
func (r *Router) NotFound(h HandlerFunc) {
	r.notFoundHandler = h
}

// GetNotFoundHandler ...
func (r *Router) GetNotFoundHandler() HandlerFunc {
	return r.notFoundHandler
}

// Match match the uri for handler
func (r *Router) Match(method, uri string, c *Context) *Route {
	for p := range r.routeMap[method] {
		if p == uri {
			return r.routeMap[method][p]
		}
	}
	return nil
}

// URLFor use named route return format url
func (r *Router) URLFor(name string, args ...interface{}) string {
	if name == "" {
		return ""
	}
	url := r.routeNamedMap[name]
	if url == "" {
		return ""
	}
	return fmt.Sprintf(url, args...)
}

// Name set name of route
func (r *Route) Name(name string) {
	if name == "" {
		return
	}
	pRune := []rune(r.pattern)
	p := make([]rune, 0, len(pRune))
	var j int
	for i, c := range pRune {
		if i < j {
			continue
		}
		j++
		if c == ':' {
			p = append(p, '%')
			p = append(p, 'v')
			for ; j < len(pRune); j++ {
				if !isParamChar(pRune[j]) {
					break
				}
			}
			continue
		}
		p = append(p, c)
	}
	r.router.routeNamedMap[name] = string(p)
	r.router.baa.Logger().Printf("debug route.name, %s \t %s", name, r.router.routeNamedMap[name])
}

// handle if ether handle return not nil then break aother handle
func (r *Route) handle(c *Context) error {
	for _, h := range r.handlers {
		err := h(c)
		if err != nil {
			return err
		}
	}
	return nil
}

// isParamChar check the char can used for route params
func isParamChar(c rune) bool {
	if (c >= 65 && c <= 90) || (c >= 97 && c <= 122) || (c >= 48 && c <= 57) {
		return true
	}
	return false
}
