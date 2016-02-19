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

const (
	// RouteMaxLength set length limit of route pattern
	RouteMaxLength = 256
	// RouterParamMaxLength set length limit of route pattern param
	RouterParamMaxLength = 64
)

var _radix [RouteMaxLength]byte
var _param [RouterParamMaxLength]byte

// Router provlider router for baa
type Router struct {
	autoHead        bool
	mu              sync.RWMutex
	notFoundHandler HandlerFunc
	routeMap        map[string]*Route
	routeNamedMap   map[string]string
}

// Route is a tree node
// route use radix tree
type Route struct {
	pattern  string
	hasParam bool
	router   *Router
	children map[string]*Route
	handlers []HandlerFunc
}

// NewRouter create a router instance
func NewRouter() *Router {
	r := new(Router)
	r.routeMap = make(map[string]*Route)
	for m := range METHODS {
		r.routeMap[m] = newRoute("/", nil, nil)
	}
	r.routeNamedMap = make(map[string]string)
	return r
}

// newRoute create a route item
func newRoute(pattern string, handles []HandlerFunc, router *Router) *Route {
	r := new(Route)
	r.pattern = pattern
	r.handlers = handles
	r.router = router
	r.children = make(map[string]*Route)
	return r
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
	ru := r.lookup(uri, r.routeMap[method])
	return ru
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

// Handle registers a new request handle with the given pattern, method and handlers.
func (r *Router) add(method string, pattern string, handlers []HandlerFunc) *Route {
	if pattern == "" {
		panic("route pattern can not be emtpy!")
	}
	if pattern[0] != '/' {
		panic("route pattern must begin /")
	}
	if len(pattern) > RouteMaxLength {
		panic(fmt.Sprintf("route pattern max length limit %d", RouteMaxLength))
	}
	if _, ok := METHODS[method]; !ok {
		panic("unsupport http method [" + method + "]")
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	root := r.routeMap[method]
	radix := _radix[0:]
	var j int
	var k int
	var tru *Route
	for i := 0; i < len(pattern); i++ {
		//param route
		if pattern[i] == ':' {
			// clear static route
			if j > 0 {
				root = r.insert(root, newRoute(string(radix[:j]), nil, nil))
				j = 0
			}
			// set param route
			root = r.insert(root, newRoute(":", nil, nil))
			param := _param[0:]
			k = 0
			for i = i + 1; i < len(pattern); i++ {
				if !isParamChar(pattern[i]) {
					i--
					break
				}
				param[k] = pattern[i]
				k++
			}
			if k == 0 {
				panic("route pattern param is empty")
			}
			if k > RouterParamMaxLength {
				panic(fmt.Sprintf("route pattern param max length limit %d", RouterParamMaxLength))
			}
			// check last character
			if i+1 == len(pattern) {
				tru = newRoute(string(param[:k]), handlers, r)
			} else {
				tru = newRoute(string(param[:k]), nil, nil)
			}
			tru.hasParam = true
			root = r.insert(root, tru)
			continue
		}
		radix[j] = pattern[i]
		j++
	}

	// static route
	if j > 0 {
		tru = newRoute(string(radix[:j]), handlers, r)
		r.insert(root, tru)
	}

	return newRoute(pattern, handlers, r)
}

func (r *Router) insert(root *Route, ru *Route) *Route {
	// same route reset root
	if root.pattern == ru.pattern {
		if ru.handlers != nil {
			ru.children = root.children
			root.reset(ru)
		}
		return root
	}

	// find radix
	var i, l int
	l = len(root.pattern)
	for i = 0; i < len(ru.pattern); i++ {
		if i >= l || ru.pattern[i] != root.pattern[i] {
			break
		}
	}
	if i > 0 && i < l {
		// has radix, and not child, reset root
		var newRu *Route
		if i == l {
			newRu = newRoute(string(ru.pattern[:i]), ru.handlers, r)
			ru = newRu
		} else {
			newRu = newRoute(string(ru.pattern[:i]), nil, nil)
			ru.pattern = ru.pattern[i:]
			newRu.children[ru.pattern] = ru
		}
		newRu.children[string(root.pattern[i:])] = &Route{
			pattern:  string(root.pattern[i:]),
			handlers: root.handlers,
			children: root.children,
			router:   r,
		}
		root.reset(newRu)
		return ru
	}

	// reset ru pattern wipe out radix
	ru.pattern = ru.pattern[i:]

	// has radix and ru is child, children is empty , driectly set
	if len(root.children) == 0 {
		root.children[ru.pattern] = ru
		return ru
	}
	// has radix and ru is child, children is note empty, continue check children radix
	for j := range root.children {
		l = len(root.children[j].pattern) - 1
		for i = 0; i < len(ru.pattern); i++ {
			if i > l || ru.pattern[i] != root.children[j].pattern[i] {
				break
			}
		}
		if i > 0 {
			ru = r.insert(root.children[j], ru)
			break
		}
	}
	// has radix and ru is child, children is note empty, but none children has radix with ru, let ru be a child
	if i == 0 {
		root.children[ru.pattern] = ru
	}
	return ru
}

func (r *Router) lookup(pattern string, root *Route) *Route {
	return nil
}

// print the route map
func (r *Router) print(prefix string, root *Route) {
	if root == nil {
		for m := range r.routeMap {
			fmt.Println(m)
			r.print("", r.routeMap[m])
		}
	}
	fmt.Println(prefix + " -> " + root.pattern)
	for i := range root.children {
		r.print(prefix+" -> "+root.pattern, root.children[i])
	}
}

// Name set name of route
func (r *Route) Name(name string) {
	if name == "" {
		return
	}
	p := make([]byte, 0, len(r.pattern))
	for i := 0; i < len(r.pattern); i++ {
		if r.pattern[i] != ':' {
			p = append(p, r.pattern[i])
			continue
		}
		p = append(p, '%')
		p = append(p, 'v')
		for i = i + 1; i < len(r.pattern); i++ {
			if !isParamChar(r.pattern[i]) {
				i--
				break
			}
		}
	}
	r.router.routeNamedMap[name] = string(p)
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

// reset route handle
func (r *Route) reset(ru *Route) {
	r.pattern = ru.pattern
	r.children = ru.children
	r.handlers = ru.handlers
	r.router = ru.router
}

// isParamChar check the char can used for route params
func isParamChar(c byte) bool {
	if (c >= 65 && c <= 90) || (c >= 97 && c <= 122) || (c >= 48 && c <= 57) {
		return true
	}
	return false
}
