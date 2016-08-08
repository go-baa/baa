package baa

import (
	"fmt"
	"sync"
)

// Tree provlider router for baa with radix tree
type Tree struct {
	autoHead          bool
	autoTrailingSlash bool
	mu                sync.RWMutex
	groups            []*group
	nodes             [RouteLength]*Leaf
	baa               *Baa
	namedNodes        map[string]string
}

// Leaf is a tree node
type Leaf struct {
	hasParam bool
	alpha    byte
	pattern  string
	param    string
	root     *Tree
	parent   *Leaf
	children []*Leaf
	handlers []HandlerFunc
}

// group route
type group struct {
	pattern  string
	handlers []HandlerFunc
}

// NewTree create a router instance
func NewTree(b *Baa) Router {
	t := new(Tree)
	for i := 0; i < len(t.nodes); i++ {
		t.nodes[i] = newLeaf("/", nil, nil)
	}
	t.namedNodes = make(map[string]string)
	t.groups = make([]*group, 0)
	t.baa = b
	return t
}

// newLeaf create a route node
func newLeaf(pattern string, handlers []HandlerFunc, root *Tree) *Leaf {
	l := new(Leaf)
	l.pattern = pattern
	l.alpha = pattern[0]
	l.handlers = handlers
	l.root = root
	l.children = make([]*Leaf, 0)
	return l
}

// newGroup create a group router
func newGroup() *group {
	g := new(group)
	g.handlers = make([]HandlerFunc, 0)
	return g
}

// SetAutoHead sets the value who determines whether add HEAD method automatically
// when GET method is added. Combo router will not be affected by this value.
func (t *Tree) SetAutoHead(v bool) {
	t.autoHead = v
}

// SetAutoTrailingSlash optional trailing slash.
func (t *Tree) SetAutoTrailingSlash(v bool) {
	t.autoTrailingSlash = v
}

// Match find matched route and returns handlerss
func (t *Tree) Match(method, pattern string, c *Context) []HandlerFunc {
	var i, l int
	var root, nn *Leaf
	root = t.nodes[RouterMethods[method]]

	for {
		// static route
		if !root.hasParam {
			l = len(root.pattern)
			if l <= len(pattern) && root.pattern == pattern[:l] {
				if l == len(pattern) {
					if root.handlers == nil {
						return nil
					}
					return root.handlers
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
				return root.handlers
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

// URLFor use named route return format url
func (t *Tree) URLFor(name string, args ...interface{}) string {
	if name == "" {
		return ""
	}
	url := t.namedNodes[name]
	if url == "" {
		return ""
	}
	return fmt.Sprintf(url, args...)
}

// Add registers a new handle with the given method, pattern and handlers.
// add check training slash option.
func (t *Tree) Add(method, pattern string, handlers []HandlerFunc) RouteNode {
	if method == "GET" && t.autoHead {
		t.add("HEAD", pattern, handlers)
	}
	if t.autoTrailingSlash && (len(pattern) > 1 || len(t.groups) > 0) {
		if pattern[len(pattern)-1] == '/' {
			pattern = pattern[:len(pattern)-1]
		}
		t.add(method, pattern+"/", handlers)
	}
	return t.add(method, pattern, handlers)
}

// GroupAdd add a group route has same prefix and handle chain
func (t *Tree) GroupAdd(pattern string, f func(), handlers []HandlerFunc) {
	g := newGroup()
	g.pattern = pattern
	g.handlers = handlers
	t.groups = append(t.groups, g)

	f()

	t.groups = t.groups[:len(t.groups)-1]
}

// add registers a new request handle with the given method, pattern and handlers.
func (t *Tree) add(method, pattern string, handlers []HandlerFunc) *Leaf {
	if _, ok := RouterMethods[method]; !ok {
		panic("unsupport http method [" + method + "]")
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	// check group set
	if len(t.groups) > 0 {
		var gpattern string
		var ghandlers []HandlerFunc
		for i := range t.groups {
			gpattern += t.groups[i].pattern
			if len(t.groups[i].handlers) > 0 {
				ghandlers = append(ghandlers, t.groups[i].handlers...)
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
		handlers[i] = WrapHandlerFunc(handlers[i])
	}

	root := t.nodes[RouterMethods[method]]
	var radix []byte
	var param []byte
	var i, k int
	var tl *Leaf
	for i = 0; i < len(pattern); i++ {
		//param route
		if pattern[i] == ':' {
			// clear static route
			if len(radix) > 0 {
				root = t.insert(root, newLeaf(string(radix), nil, nil))
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
				tl = newLeaf(":", handlers, t)
			} else {
				tl = newLeaf(":", nil, nil)
			}
			tl.param = string(param[:k])
			tl.hasParam = true
			root = t.insert(root, tl)
			continue
		}
		radix = append(radix, pattern[i])
	}

	// static route
	if len(radix) > 0 {
		tl = newLeaf(string(radix), handlers, t)
		t.insert(root, tl)
	}

	return newLeaf(pattern, handlers, t)
}

// insert build the route tree
func (t *Tree) insert(root *Leaf, node *Leaf) *Leaf {
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
				return t.insert(root.children[i], node)
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
				return t.insert(root.children[i], node)
			}
		}

		root.insertChild(node)
		return node
	}

	// _parent root and ru has new parent
	_parent := newLeaf(root.pattern[:pos], nil, t)
	root.parent.deleteChild(root)
	root.parent.insertChild(_parent)
	root.resetPattern(root.pattern[pos:])
	_parent.insertChild(root)
	_parent.insertChild(node)

	return node
}

// Name set name of route
func (l *Leaf) Name(name string) {
	if name == "" {
		return
	}
	p := make([]byte, 0, len(l.pattern))
	for i := 0; i < len(l.pattern); i++ {
		if l.pattern[i] != ':' {
			p = append(p, l.pattern[i])
			continue
		}
		p = append(p, '%')
		p = append(p, 'v')
		for i = i + 1; i < len(l.pattern); i++ {
			if l.pattern[i] == '/' {
				i--
				break
			}
		}
	}
	l.root.namedNodes[name] = string(p)
}

// findChild find child static route
func (l *Leaf) findChild(b byte) *Leaf {
	var i int
	var j = len(l.children)
	for ; i < j; i++ {
		if l.children[i].alpha == b && !l.children[i].hasParam {
			return l.children[i]
		}
	}
	return nil
}

// deleteChild find child and delete from root route
func (l *Leaf) deleteChild(node *Leaf) {
	for i := 0; i < len(l.children); i++ {
		if l.children[i].pattern != node.pattern {
			continue
		}
		if len(l.children) == 1 {
			l.children = l.children[:0]
			return
		}
		if i == 0 {
			l.children = l.children[1:]
			return
		}
		if i+1 == len(l.children) {
			l.children = l.children[:i]
			return
		}
		l.children = append(l.children[:i], l.children[i+1:]...)
		return
	}
}

// insertChild insert child into root route, and returns the child route
func (l *Leaf) insertChild(node *Leaf) *Leaf {
	var i int
	for ; i < len(l.children); i++ {
		if l.children[i].pattern == node.pattern {
			if l.children[i].hasParam && node.hasParam && l.children[i].param != node.param {
				panic("Router Tree.insert error cannot use two param [:" + l.children[i].param + ", :" + node.param + "]with same prefix!")
			}
			if node.handlers != nil {
				l.children[i].handlers = node.handlers
			}
			return l.children[i]
		}
	}
	node.parent = l
	l.children = append(l.children, node)

	i = len(l.children) - 1
	if i > 0 && l.children[i].hasParam {
		l.children[0], l.children[i] = l.children[i], l.children[0]
		return l.children[0]
	}
	return node
}

// resetPattern reset route pattern and alpha
func (l *Leaf) resetPattern(pattern string) {
	l.pattern = pattern
	l.alpha = pattern[0]
}

// hasPrefixString returns the same prefix position, if none return 0
func (l *Leaf) hasPrefixString(s string) int {
	var i, j int
	j = len(l.pattern)
	if len(s) < j {
		j = len(s)
	}
	for i = 0; i < j && s[i] == l.pattern[i]; i++ {
	}
	return i
}
