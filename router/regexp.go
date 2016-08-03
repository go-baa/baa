package router

import (
	"sync"

	"github.com/go-baa/baa"
)

// Reuter provlider router for baa with regexp
type Reuter struct {
	autoHead          bool
	autoTrailingSlash bool
	mu                sync.RWMutex
	notFoundHandler   baa.HandlerFunc
	groups            []*group
	nodes             [baa.RouteLength]*Node
	namedNodes        map[string]string
}

// Node is a router node
type Node struct {
	hasParam bool
	alpha    byte
	pattern  string
	param    string
	root     *Reuter
	parent   *Node
	children []*Node
	handlers []baa.HandlerFunc
}

// group route
type group struct {
	pattern  string
	handlers []baa.HandlerFunc
}

// newTree create a router instance
func NewRegexp() baa.Router {
	r := new(Reuter)
	for i := 0; i < len(r.nodes); i++ {
		r.nodes[i] = newNode("/", nil, nil)
	}
	r.namedNodes = make(map[string]string)
	r.groups = make([]*group, 0)
	return r
}

// newNode create a route node
func newNode(pattern string, handles []baa.HandlerFunc, root *Reuter) *Node {
	n := new(Node)
	n.pattern = pattern
	n.alpha = pattern[0]
	n.handlers = handles
	n.root = root
	n.children = make([]*Node, 0)
	return n
}

// newGroup create a group router
func newGroup() *group {
	g := new(group)
	g.handlers = make([]baa.HandlerFunc, 0)
	return g
}

// SetAutoHead sets the value who determines whether add HEAD method automatically
// when GET method is added. Combo router will not be affected by this value.
func (r *Reuter) SetAutoHead(v bool) {
	r.autoHead = v
}

// SetAutoTrailingSlash optional trailing slash.
func (r *Reuter) SetAutoTrailingSlash(v bool) {
	r.autoTrailingSlash = v
}

// Match match the route
func (r *Reuter) Match(method, uri string, c *baa.Context) baa.RouteNode {
	return nil
}

// URLFor use named route return format url
func (r *Reuter) URLFor(name string, args ...interface{}) string {
	return ""
}

// Add registers a new handle with the given method, pattern and handlers.
func (r *Reuter) Add(method, pattern string, handlers []baa.HandlerFunc) baa.RouteNode {
	return nil
}

// GroupAdd registers a list of same prefix route
func (r *Reuter) GroupAdd(pattern string, f func(), handlers []baa.HandlerFunc) {

}

// Handlers returns handlers bond with leaf
func (n *Node) Handlers() []baa.HandlerFunc {
	return n.handlers
}

// Name set name of route
func (n *Node) Name(name string) {
	if name == "" {
		return
	}
	p := make([]byte, 0, len(n.pattern))
	for i := 0; i < len(n.pattern); i++ {
		if n.pattern[i] != ':' {
			p = append(p, n.pattern[i])
			continue
		}
		p = append(p, '%')
		p = append(p, 'v')
		for i = i + 1; i < len(n.pattern); i++ {
			if n.pattern[i] == '/' {
				i--
				break
			}
		}
	}
	n.root.namedNodes[name] = string(p)
}
