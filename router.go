package baa

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// METHODS declare support HTTP method
var METHODS = map[string]bool{
	"GET":     true,
	"POST":    true,
	"PUT":     true,
	"DELETE":  true,
	"PATCH":   true,
	"OPTIONS": true,
	"HEAD":    true,
}

// optimize ...
var _radix [RouteMaxLength]byte
var _param [RouterParamMaxLength]byte

// Router provlider router for baa
type Router struct {
	autoHead        bool
	mu              sync.RWMutex
	notFoundHandler HandlerFunc
	routeMap        map[string]*Route
	routeNamedMap   map[string]string
	group           *Group
}

// Route is a tree node
// route use radix tree
type Route struct {
	pattern  string
	hasParam bool
	parent   *Route
	router   *Router
	children map[string]*Route
	handlers []HandlerFunc
}

// Group route
type Group struct {
	pattern  string
	handlers []HandlerFunc
	mu       sync.RWMutex
}

// newRouter create a router instance
func newRouter() *Router {
	r := new(Router)
	r.routeMap = make(map[string]*Route)
	for m := range METHODS {
		r.routeMap[m] = newRoute("/", nil, nil)
	}
	r.routeNamedMap = make(map[string]string)
	r.group = newGroup()
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

// newGroup create a group router
func newGroup() *Group {
	g := new(Group)
	g.handlers = make([]HandlerFunc, 0)
	return g
}

// setNotFoundHandler set the route not match result.
// Configurable http.HandlerFunc which is called when no matching route is
// found. If it is not set, http.NotFound is used.
// Be sure to set 404 response code in your handler.
func (r *Router) setNotFoundHandler(h HandlerFunc) {
	r.notFoundHandler = h
}

// NotFound execute 404 handler
func (r *Router) NotFound(c *Context) {
	if r.notFoundHandler != nil {
		r.notFoundHandler(c)
		return
	}
	http.NotFound(c.Resp, c.Req)
}

// Match match the uri for handler
func (r *Router) match(method, uri string, c *Context) *Route {
	ru := r.lookup(uri, r.routeMap[method], c)
	if ru != nil && ru.handlers != nil {
		return ru
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

// groupAdd add a group route has same prefix and handle chain
func (r *Router) groupAdd(pattern string, f func(), handlers []HandlerFunc) {
	r.group.mu.Lock()
	defer r.group.mu.Unlock()

	r.group.pattern = pattern
	r.group.handlers = handlers

	f()

	r.group.reset()
}

// Handle registers a new request handle with the given pattern, method and handlers.
func (r *Router) add(method string, pattern string, handlers []HandlerFunc) *Route {
	if _, ok := METHODS[method]; !ok {
		panic("unsupport http method [" + method + "]")
	}
	if pattern == "" {
		panic("route pattern can not be emtpy!")
	}
	if pattern[0] != '/' {
		panic("route pattern must begin /")
	}
	if len(pattern) > RouteMaxLength {
		panic(fmt.Sprintf("route pattern max length limit %d", RouteMaxLength))
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	// check group set, not concurrent safe
	if r.group.pattern != "" {
		pattern = r.group.pattern + pattern
		if len(r.group.handlers) > 0 {
			h := make([]HandlerFunc, 0, len(r.group.handlers)+len(handlers))
			h = append(h, r.group.handlers...)
			h = append(h, handlers...)
			handlers = h
		}
	}

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
			p := ":" + string(param[:k])
			if i == len(pattern) {
				tru = newRoute(p, handlers, r)
			} else {
				tru = newRoute(p, nil, nil)
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

	// param route
	if root.hasParam {
		if ru.hasParam {
			if parent := root.parent; parent != nil {
				ru.parent = parent
				parent.children[ru.pattern] = ru
			}
		} else {
			ru.parent = root
			root.children[ru.pattern] = ru
		}
		return ru
	}

	// find radix
	var i, l int
	l = len(root.pattern)
	if root.pattern[0] == ':' && ru.pattern[0] == ':' {

	} else {
		for i = 0; i < len(ru.pattern); i++ {
			// param route can not split
			if i >= l || ru.pattern[i] != root.pattern[i] {
				break
			}
		}
	}
	if i > 0 && i < l {
		// has radix, and not child, reset root
		var newRu *Route
		if i <= l {
			newRu = newRoute(ru.pattern[:i], ru.handlers, r)
			ru = newRu
		} else {
			newRu = newRoute(ru.pattern[:i], nil, nil)
			ru.pattern = ru.pattern[i:]
			ru.parent = newRu
			newRu.children[ru.pattern] = ru
		}
		newRu.children[root.pattern[i:]] = &Route{
			pattern:  root.pattern[i:],
			handlers: root.handlers,
			children: root.children,
			router:   r,
			parent:   newRu,
		}
		root.reset(newRu)
		return ru
	}

	// reset ru pattern wipe out radix
	ru.pattern = ru.pattern[i:]

	// has radix and ru is child, children is note empty, continue check children radix
	i = 0
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
		ru.parent = root
		root.children[ru.pattern] = ru
	}
	return ru
}

func (r *Router) lookup(pattern string, root *Route, c *Context) *Route {
	var ru *Route
	// static route
	if !root.hasParam {
		if pattern == root.pattern {
			return root
		}
		if strings.HasPrefix(pattern, root.pattern) {
			pattern = pattern[len(root.pattern):]
		} else {
			return nil
		}
	} else {
		var i int
		if len(root.children) == 0 {
			i = len(pattern)
		} else {
			for i = 0; i < len(pattern); i++ {
				if !isParamChar(pattern[i]) {
					break
				}
			}
		}
		c.SetParam(root.pattern[1:], pattern[:i])
		if i == len(pattern) {
			if root.handlers != nil {
				return root
			}
			return nil
		}
		pattern = pattern[i:]
	}
	if len(root.children) == 0 {
		return nil
	}

	// first, static route
	for _, v := range root.children {
		if ru = r.lookup(pattern, v, c); ru != nil {
			//if ru.handlers != nil {
			return ru
			//}
		}
	}

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

// reset group data for next group set
func (g *Group) reset() {
	g.pattern = ""
	g.handlers = g.handlers[:0]
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

// handle route hadnle chain
// if something wrote to http, break chain and return
func (r *Route) handle(c *Context) {
	for _, h := range r.handlers {
		h(c)
		if c.Resp.Wrote() {
			return
		}
	}
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
	if (c >= 65 && c <= 90) || (c >= 97 && c <= 122) || (c >= 48 && c <= 57) || c == 95 {
		return true
	}
	return false
}
