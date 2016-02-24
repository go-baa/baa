package baa

import (
	"fmt"
	"sync"
)

const (
	// RouteMaxLength set length limit of route pattern
	routeMaxLength = 256
	// RouterParamMaxLength set length limit of route pattern param
	routerParamMaxLength = 32
)

const (
	// method key in routeMap
	GET int = iota
	POST
	PUT
	DELETE
	PATCH
	OPTIONS
	HEAD
	// RouteLength route table length
	RouteLength
)

// methods declare support HTTP method
var methods = map[string]bool{
	"GET":     true,
	"POST":    true,
	"PUT":     true,
	"DELETE":  true,
	"PATCH":   true,
	"OPTIONS": true,
	"HEAD":    true,
}

// methodKeys declare method key in routeMap
var methodKeys = map[string]int{
	"GET":     GET,
	"POST":    POST,
	"PUT":     PUT,
	"DELETE":  DELETE,
	"PATCH":   PATCH,
	"OPTIONS": OPTIONS,
	"HEAD":    HEAD,
}

// optimize ...
var _radix [routeMaxLength]byte
var _param [routerParamMaxLength]byte

// Router provlider router for baa
type Router struct {
	autoHead        bool
	mu              sync.RWMutex
	notFoundHandler HandlerFunc
	routeMap        [RouteLength]*Route
	routeNamedMap   map[string]string
	groups          []*group
}

// Route is a tree node
// route use radix tree
type Route struct {
	pattern  string
	hasParam bool
	parent   *Route
	router   *Router
	children []*Route
	handlers []HandlerFunc
}

// group route
type group struct {
	pattern  string
	handlers []HandlerFunc
}

// newRouter create a router instance
func newRouter() *Router {
	r := new(Router)
	for i := 0; i < len(r.routeMap); i++ {
		r.routeMap[i] = newRoute("/", nil, nil)
	}
	r.routeNamedMap = make(map[string]string)
	r.groups = make([]*group, 0)
	return r
}

// newRoute create a route item
func newRoute(pattern string, handles []HandlerFunc, router *Router) *Route {
	r := new(Route)
	r.pattern = pattern
	r.handlers = handles
	r.router = router
	r.children = make([]*Route, 0)
	return r
}

// newGroup create a group router
func newGroup() *group {
	g := new(group)
	g.handlers = make([]HandlerFunc, 0)
	return g
}

// Match match the uri for handler
func (r *Router) match(method, uri string, c *Context) *Route {
	return r.lookup(uri, r.routeMap[methodKeys[method]], c)
}

// urlFor use named route return format url
func (r *Router) urlFor(name string, args ...interface{}) string {
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
	g := newGroup()
	g.pattern = pattern
	g.handlers = handlers
	r.groups = append(r.groups, g)

	f()

	r.groups = r.groups[:len(r.groups)-1]
}

// Handle registers a new request handle with the given pattern, method and handlers.
func (r *Router) add(method string, pattern string, handlers []HandlerFunc) *Route {
	if _, ok := methods[method]; !ok {
		panic("unsupport http method [" + method + "]")
	}
	if pattern == "" {
		panic("route pattern can not be emtpy!")
	}
	if pattern[0] != '/' {
		panic("route pattern must begin /")
	}
	if len(pattern) > routeMaxLength {
		panic(fmt.Sprintf("route pattern max length limit %d", routeMaxLength))
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	// check group set
	if len(r.groups) > 0 {
		var gpattern string
		var ghandlers []HandlerFunc
		for i := range r.groups {
			gpattern += r.groups[i].pattern
			if len(r.groups[i].handlers) > 0 {
				ghandlers = append(ghandlers, r.groups[i].handlers...)
			}
		}
		pattern = gpattern + pattern
		ghandlers = append(ghandlers, handlers...)
		handlers = ghandlers
	}

	for i := 0; i < len(handlers); i++ {
		handlers[i] = wrapHandlerFunc(handlers[i])
	}

	root := r.routeMap[methodKeys[method]]
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
			if k > routerParamMaxLength {
				panic(fmt.Sprintf("route pattern param max length limit %d", routerParamMaxLength))
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
func (r *Router) insert(root *Route, node *Route) *Route {
	fmt.Printf("->insert: root: %s, node: %s\n", root.pattern, node.pattern)
	// same route
	if root.pattern == node.pattern {
		if node.handlers != nil {
			root.handlers = node.handlers
		}
		return root
	}

	// param route
	if root.hasParam && node.hasParam {
		if root.parent == nil {
			panic("Router.insert error route has no parent")
		}
		return root.parent.insertChild(node)
	}

	if !root.hasParam && node.hasParam {
		return root.insertChild(node)
	}

	var i int
	if root.hasParam && !node.hasParam {
		for i = range root.children {
			if root.children[i].pattern == node.pattern {
				if node.handlers != nil {
					root.children[i].handlers = node.handlers
				}
				return root.children[i]
			}
			if root.children[i].hasPrefix(node) > 0 {
				return r.insert(root.children[i], node)
			}
		}

		root.insertChild(node)
		return node
	}

	// find radix
	pos := root.hasPrefix(node)
	if pos == 0 {
		panic("Router.insert error root[" + root.pattern + "] and node[" + node.pattern + "] not have both prefix")
	}
	if pos == len(node.pattern) {
		root.parent.deleteChild(root)
		root.parent.insertChild(node)
		root.pattern = root.pattern[pos:]
		node.insertChild(root)
		return node
	}

	node.pattern = node.pattern[pos:]
	if pos == len(root.pattern) {
		for i = range root.children {
			if root.children[i].pattern == node.pattern {
				if node.handlers != nil {
					root.children[i].handlers = node.handlers
				}
				return root.children[i]
			}
			if root.children[i].hasPrefix(node) > 0 {
				return r.insert(root.children[i], node)
			}
		}

		root.insertChild(node)
		return node
	}

	// _parent root and ru has new parent
	_parent := newRoute(root.pattern[:pos], nil, r)
	root.parent.deleteChild(root)
	root.parent.insertChild(_parent)
	root.pattern = root.pattern[pos:]
	_parent.insertChild(root)
	_parent.insertChild(node)

	return node
}

func (r *Router) lookup(pattern string, root *Route, c *Context) *Route {
	var ru *Route
	var i int
	// static route
	if !root.hasParam {
		if pattern == root.pattern {
			return root
		}
		if len(pattern) >= len(root.pattern) && pattern[0:len(root.pattern)] == root.pattern {
			pattern = pattern[len(root.pattern):]
		} else {
			return nil
		}
	} else {
		if len(root.children) == 0 {
			i = len(pattern)
		} else {
			for i = 0; i < len(pattern) && isParamChar(pattern[i]); i++ {
			}
		}
		c.SetParam(root.pattern[1:], pattern[:i])
		if i == len(pattern) {
			return root
		}
		pattern = pattern[i:]
	}
	if len(root.children) == 0 {
		return nil
	}

	// first, static route
	for i = range root.children {
		if ru = r.lookup(pattern, root.children[i], c); ru != nil {
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

// deleteChild find child and delete from root route
func (r *Route) deleteChild(child *Route) {
	for i := 0; i < len(r.children); i++ {
		if r.children[i].pattern != child.pattern {
			continue
		}
		if len(r.children) == 1 {
			r.children = r.children[:0]
			return
		}
		if i == 0 {
			r.children = r.children[1:]
			return
		}
		if i+1 == len(r.children) {
			r.children = r.children[:i]
			return
		}
		r.children = append(r.children[:i], r.children[i+1:]...)
		return
	}
}

// insertChild insert child into root route, and returns the child route
func (r *Route) insertChild(child *Route) *Route {
	for i := 0; i < len(r.children); i++ {
		if r.children[i].pattern == child.pattern {
			if child.handlers != nil {
				r.children[i].handlers = child.handlers
			}
			return r.children[i]
		}
	}
	child.parent = r
	r.children = append(r.children, child)
	return child
}

// hasChild check root has child, if yes return child route, or reutrn nil
func (r *Route) hasChild(child *Route) *Route {
	for i := 0; i < len(r.children); i++ {
		if r.children[i].pattern == child.pattern {
			return r.children[i]
		}
	}
	return nil
}

// hasPrefix returns the same prefix position, if none return 0
func (r *Route) hasPrefix(ru *Route) int {
	l := len(r.pattern)
	var i int
	for i = 0; i < len(ru.pattern) && i < l && ru.pattern[i] == r.pattern[i]; i++ {
	}
	return i
}

// wrapHandlerFunc wrap for context handler chain
func wrapHandlerFunc(h HandlerFunc) HandlerFunc {
	return func(c *Context) {
		h(c)
		c.Next()
	}
}

// isParamChar check the char can used for route params
// a-z->65:90, A-Z->97:122, 0-9->48->57, _->95, :->58
func isParamChar(c byte) bool {
	if (c >= 65 && c <= 90) || (c >= 97 && c <= 122) || (c >= 48 && c <= 58) || c == 95 {
		return true
	}
	return false
}
