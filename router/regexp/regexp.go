package regexp

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/go-baa/baa"
)

// Router provlider router for baa with regexp
type Router struct {
	autoHead          bool
	autoTrailingSlash bool
	mu                sync.RWMutex
	groups            []*group
	nodes             [baa.RouteLength][]*Node
	baa               *baa.Baa
	namedNodes        map[string]*Node
}

// Node is a router node
type Node struct {
	hasParam bool
	pattern  string
	format   []byte
	params   []string
	re       *regexp.Regexp
	root     *Router
	handlers []baa.HandlerFunc
}

// group route
type group struct {
	pattern  string
	handlers []baa.HandlerFunc
}

// New create a router instance
func New(b *baa.Baa) baa.Router {
	r := new(Router)
	for i := 0; i < len(r.nodes); i++ {
		r.nodes[i] = make([]*Node, 0)
	}
	r.namedNodes = make(map[string]*Node)
	r.groups = make([]*group, 0)
	r.baa = b
	return r
}

// newNode create a route node
func newNode(pattern string, handles []baa.HandlerFunc, root *Router) *Node {
	n := new(Node)
	n.pattern = pattern
	n.handlers = handles
	n.root = root
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
func (r *Router) SetAutoHead(v bool) {
	r.autoHead = v
}

// SetAutoTrailingSlash optional trailing slash.
func (r *Router) SetAutoTrailingSlash(v bool) {
	r.autoTrailingSlash = v
}

// Match match the route
func (r *Router) Match(method, uri string, c *baa.Context) baa.RouteNode {
	for _, n := range r.nodes[baa.RouterMethods[method]] {
		if !n.hasParam {
			if n.pattern == uri {
				return n
			}
			continue
		}
		data := n.re.FindStringSubmatch(uri)
		if len(data) != len(n.params)+1 {
			continue
		}
		for i := range n.params {
			c.SetParam(n.params[i], data[i+1])
		}
		return n
	}
	return nil
}

// URLFor use named route return format url
func (r *Router) URLFor(name string, args ...interface{}) string {
	if name == "" {
		return ""
	}
	node := r.namedNodes[name]
	if node == nil || len(node.format) == 0 {
		return ""
	}
	format := make([]byte, len(node.format))
	copy(format, node.format)
	for i := len(node.params) + 1; i <= len(args); i++ {
		format = append(format, "%v"...)
	}
	return fmt.Sprintf(string(format), args...)
}

// Add registers a new handle with the given method, pattern and handlers.
// add check training slash option.
func (r *Router) Add(method, pattern string, handlers []baa.HandlerFunc) baa.RouteNode {
	if method == "GET" && r.autoHead {
		r.add("HEAD", pattern, handlers)
	}
	if r.autoTrailingSlash && (len(pattern) > 1 || len(r.groups) > 0) {
		if pattern[len(pattern)-1] == '/' {
			pattern = pattern[:len(pattern)-1]
		}
		r.add(method, pattern+"/", handlers)
	}
	return r.add(method, pattern, handlers)
}

// GroupAdd add a group route has same prefix and handle chain
func (r *Router) GroupAdd(pattern string, f func(), handlers []baa.HandlerFunc) {
	g := newGroup()
	g.pattern = pattern
	g.handlers = handlers
	r.groups = append(r.groups, g)

	f()

	r.groups = r.groups[:len(r.groups)-1]
}

// add registers a new request handle with the given method, pattern and handlers.
func (r *Router) add(method, pattern string, handlers []baa.HandlerFunc) *Node {
	if _, ok := baa.RouterMethods[method]; !ok {
		panic("Router.add: unsupport http method [" + method + "]")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// check group set
	if len(r.groups) > 0 {
		var gpattern string
		var ghandlers []baa.HandlerFunc
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
		panic("Router.add: pattern can not be emtpy!")
	}
	if pattern[0] != '/' {
		panic("Router.add: pattern must begin /")
	}

	for i := 0; i < len(handlers); i++ {
		handlers[i] = baa.WrapHandlerFunc(handlers[i])
	}

	tn := newNode("", handlers, r)
	var i, j, k int
	var param, restr string
	var tpattern []byte
	for i = 0; i < len(pattern); i++ {
		if pattern[i] == ':' {
			// in param route
			for j = i + 1; j < len(pattern) && baa.IsParamChar(pattern[j]); j++ {
			}
			param = pattern[i+1 : j]
			if param == "" {
				panic("Router.add: pattern param is empty")
			}
			i = j - 1
			// check regexp
			restr = ""
			if j < len(pattern) && pattern[j] == '(' {
				for k = j + 1; k < len(pattern) && pattern[k] != ')'; k++ {
				}
				restr = pattern[j+1 : k]
				i = k
			}
			if restr == "" {
				restr = "([^\\/]+)"
			} else if restr == "int" {
				restr = "([\\d]+)"
			} else if restr == "string" {
				restr = "([\\w]+)"
			} else {
				restr = "(" + restr + ")"
			}
			tpattern = append(tpattern, restr...)
			tn.format = append(tn.format, "%v"...)
			tn.params = append(tn.params, param)
			tn.hasParam = true
			continue
		}
		tpattern = append(tpattern, pattern[i])
		tn.format = append(tn.format, pattern[i])
	}
	tn.pattern = string(tpattern)
	if tn.hasParam {
		var err error
		tn.re, err = regexp.Compile(tn.pattern + "$")
		if err != nil {
			panic("Router.add: " + err.Error())
		}
	}
	// check repeat route
	if r.baa.Debug() && tn.pattern != "/" {
		for _, n := range r.nodes[baa.RouterMethods[method]] {
			if n.pattern == tn.pattern {
				panic("Router.add: route already exist -> " + tn.pattern)
			}
		}
	}
	r.nodes[baa.RouterMethods[method]] = append(r.nodes[baa.RouterMethods[method]], tn)
	return tn
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
	n.root.namedNodes[name] = n
}
