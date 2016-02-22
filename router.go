package baa

import (
	"fmt"
	"strings"
	"sync"
)

const (
	// RouteMaxLength set length limit of route pattern
	RouteMaxLength = 256
	// RouterParamMaxLength set length limit of route pattern param
	RouterParamMaxLength = 32
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
	radix := _radix[:0]
	var i, k int
	var tru *Route
	for i = 0; i < len(pattern); i++ {
		//param route
		if pattern[i] == ':' {
			// clear static route
			if len(radix) > 0 {
				root = r.insert(root, newRoute(string(radix), nil, nil))
				radix = _radix[:0]
			}
			// set param route
			param := _param[:0]
			k = 0
			for i = i + 1; i < len(pattern); i++ {
				if !isParamChar(pattern[i]) {
					i--
					break
				}
				param = append(param, pattern[i])
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
		radix = append(radix, pattern[i])
	}

	// static route
	if len(radix) > 0 {
		tru = newRoute(string(radix), handlers, r)
		r.insert(root, tru)
	}

	return newRoute(pattern, handlers, r)
}

// insert build the route tree
func (r *Router) insert(root *Route, ru *Route) *Route {
	// same route
	if root.pattern == ru.pattern {
		if ru.handlers != nil {
			root.handlers = ru.handlers
		}
		return root
	}

	// param route
	if root.hasParam && ru.hasParam {
		if root.parent == nil {
			panic("Router.insert error route has no parent")
		}
		ru.parent = root.parent
		root.parent.children[ru.pattern] = ru
		return ru
	}

	var k string
	if root.hasParam && !ru.hasParam {
		for k = range root.children {
			if root.children[k].pattern == ru.pattern {
				if ru.handlers != nil {
					root.children[k].handlers = ru.handlers
				}
				return root.children[k]
			}
			if hasPrefix(root.children[k].pattern, ru.pattern) > 0 {
				return r.insert(root.children[k], ru)
			}
		}

		ru.parent = root
		root.children[ru.pattern] = ru
		return ru
	}

	var ok bool
	if !root.hasParam && ru.hasParam {
		if _, ok = root.children[ru.pattern]; ok {
			if ru.handlers != nil {
				root.children[ru.pattern].handlers = ru.handlers
			}
		} else {
			ru.parent = root
			root.children[ru.pattern] = ru
		}
		return root.children[ru.pattern]
	}

	// find radix
	pos := hasPrefix(root.pattern, ru.pattern)
	if pos == 0 {
		panic("Router.insert error root[" + root.pattern + "] and node[" + ru.pattern + "] not have both prefix")
	}
	if pos == len(ru.pattern) {
		ru.parent = root.parent
		ru.children[root.pattern] = root
		delete(root.parent.children, root.pattern)
		root.parent.children[ru.pattern] = ru
		root.pattern = root.pattern[pos:]
		root.parent = ru
		return ru
	}

	ru.pattern = ru.pattern[pos:]
	if pos == len(root.pattern) {
		for k = range root.children {
			if root.children[k].pattern == ru.pattern {
				if ru.handlers != nil {
					root.children[k].handlers = ru.handlers
				}
				return root.children[k]
			}
			if hasPrefix(root.children[k].pattern, ru.pattern) > 0 {
				return r.insert(root.children[k], ru)
			}
		}

		ru.parent = root
		root.children[ru.pattern] = ru
		return ru
	}

	delete(root.parent.children, root.pattern)
	_root := newRoute(root.pattern[:pos], nil, r)
	_root.parent = root.parent
	ru.parent = _root
	newRoot := newRoute(root.pattern[pos:], root.handlers, r)
	newRoot.children = root.children
	newRoot.parent = _root
	_root.children[newRoot.pattern] = newRoot
	_root.children[ru.pattern] = ru
	root.parent.children[_root.pattern] = _root
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
				// find allow pattern contains :
				if !isParamChar(pattern[i]) && pattern[i] != ':' {
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
			if ru.handlers != nil {
				return ru
			}
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
		return
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

// isParamChar check the char can used for route params
func isParamChar(c byte) bool {
	if (c >= 65 && c <= 90) || (c >= 97 && c <= 122) || (c >= 48 && c <= 57) || c == 95 {
		return true
	}
	return false
}

// hasPrefix returns the same prefix position
func hasPrefix(s1, s2 string) int {
	l := len(s1)
	var i int
	for i = 0; i < len(s2); i++ {
		if i >= l || s2[i] != s1[i] {
			break
		}
	}
	return i
}
