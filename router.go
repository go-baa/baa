package baa

import (
	"fmt"
	"sync"
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

// Router provlider router for baa
type Router struct {
	autoHead          bool
	autoTrailingSlash bool
	mu                sync.RWMutex
	notFoundHandler   HandlerFunc
	groups            []*group
	routeMap          [RouteLength]*Route
	routeNamedMap     map[string]string
}

// Route is a tree node
// route use radix tree
type Route struct {
	hasParam bool
	alpha    byte
	pattern  string
	param    string
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
	r.alpha = pattern[0]
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

// methods returns all support methods
func (r *Router) methods() []string {
	var ms []string
	for m := range methods {
		ms = append(ms, m)
	}
	return ms
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

	// check pattern (for training slash move behind group check)
	if pattern == "" {
		panic("route pattern can not be emtpy!")
	}
	if pattern[0] != '/' {
		panic("route pattern must begin /")
	}

	for i := 0; i < len(handlers); i++ {
		handlers[i] = wrapHandlerFunc(handlers[i])
	}

	root := r.routeMap[methodKeys[method]]
	var radix []byte
	var param []byte
	var i, k int
	var tru *Route
	for i = 0; i < len(pattern); i++ {
		//param route
		if pattern[i] == ':' {
			// clear static route
			if len(radix) > 0 {
				root = r.insert(root, newRoute(string(radix), nil, nil))
				radix = radix[:0]
			}
			// set param route
			param = param[:0]
			k = 0
			for i = i + 1; i < len(pattern); i++ {
				if pattern[i] == '/' {
					i--
					break
				}
				param = append(param, pattern[i])
				k++
			}
			if k == 0 {
				panic("route pattern param is empty")
			}
			// check last character
			if i == len(pattern) {
				tru = newRoute(":", handlers, r)
			} else {
				tru = newRoute(":", nil, nil)
			}
			tru.param = string(param[:k])
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
	// same route
	if root.pattern == node.pattern {
		if node.handlers != nil {
			root.handlers = node.handlers
		}
		return root
	}

	// param route
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
			if root.children[i].hasPrefixString(node.pattern) > 0 {
				return r.insert(root.children[i], node)
			}
		}

		root.insertChild(node)
		return node
	}

	// find radix
	pos := root.hasPrefixString(node.pattern)
	if pos == len(node.pattern) {
		root.parent.deleteChild(root)
		root.parent.insertChild(node)
		root.resetPattern(root.pattern[pos:])
		node.insertChild(root)
		return node
	}

	node.resetPattern(node.pattern[pos:])
	if pos == len(root.pattern) {
		for i = range root.children {
			if root.children[i].pattern == node.pattern {
				if node.handlers != nil {
					root.children[i].handlers = node.handlers
				}
				return root.children[i]
			}
			if root.children[i].hasPrefixString(node.pattern) > 0 {
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
	root.resetPattern(root.pattern[pos:])
	_parent.insertChild(root)
	_parent.insertChild(node)

	return node
}

// match match the route
func (r *Router) match(method, pattern string, c *Context) *Route {
	var i, l int
	var root, nn *Route
	root = r.routeMap[methodKeys[method]]

	for {
		// static route
		if !root.hasParam {
			l = len(root.pattern)
			if l <= len(pattern) && root.pattern == pattern[:l] {
				if l == len(pattern) {
					return root
				}
				if len(root.children) == 0 {
					return nil
				}
				pattern = pattern[l:]
			} else {
				return nil
			}
		} else {
			// params route
			l = len(pattern)
			if len(root.children) == 0 {
				i = l
			} else {
				for i = 0; i < l && pattern[i] != '/'; i++ {
				}
			}
			c.SetParam(root.param, pattern[:i])
			if i == l {
				return root
			}
			pattern = pattern[i:]
		}

		// only one child
		if len(root.children) == 1 {
			if root.children[0].hasParam || root.children[0].alpha == pattern[0] {
				root = root.children[0]
				continue
			}
			break
		}

		// children static route
		if nn = root.findChild(pattern[0]); nn != nil {
			root = nn
			continue
		}

		// children param route
		if root.children[0].hasParam {
			root = root.children[0]
			continue
		}

		break
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
			if r.pattern[i] == '/' {
				i--
				break
			}
		}
	}
	r.router.routeNamedMap[name] = string(p)
}

// findChild find child static route
func (r *Route) findChild(b byte) *Route {
	var i int
	var l = len(r.children)
	for ; i < l; i++ {
		if r.children[i].alpha == b && !r.children[i].hasParam {
			return r.children[i]
		}
	}
	return nil
}

// deleteChild find child and delete from root route
func (r *Route) deleteChild(node *Route) {
	for i := 0; i < len(r.children); i++ {
		if r.children[i].pattern != node.pattern {
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
func (r *Route) insertChild(node *Route) *Route {
	var i int
	for ; i < len(r.children); i++ {
		if r.children[i].pattern == node.pattern {
			if r.children[i].hasParam && node.hasParam && r.children[i].param != node.param {
				panic("Router.insert error cannot use two param [:" + r.children[i].param + ", :" + node.param + "]with same prefix!")
			}
			if node.handlers != nil {
				r.children[i].handlers = node.handlers
			}
			return r.children[i]
		}
	}
	node.parent = r
	r.children = append(r.children, node)

	i = len(r.children) - 1
	if i > 0 && r.children[i].hasParam {
		r.children[0], r.children[i] = r.children[i], r.children[0]
		return r.children[0]
	}
	return node
}

// resetPattern reset route pattern and alpha
func (r *Route) resetPattern(pattern string) {
	r.pattern = pattern
	r.alpha = pattern[0]
}

// hasPrefixString returns the same prefix position, if none return 0
func (r *Route) hasPrefixString(s string) int {
	var i, l int
	l = len(r.pattern)
	if len(s) < l {
		l = len(s)
	}
	for i = 0; i < l && s[i] == r.pattern[i]; i++ {
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
