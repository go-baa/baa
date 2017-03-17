package baa

import (
	"fmt"
	"sync"
)

const (
	leafKindStatic uint = iota
	leafKindParam
	leafKindWide
)

// Tree provlider router for baa with radix tree
type Tree struct {
	autoHead          bool
	autoTrailingSlash bool
	mu                sync.RWMutex
	groups            []*group
	nodes             [RouteLength]*leaf
	baa               *Baa
	nameNodes         map[string]*Node
}

// Node is struct for named route
type Node struct {
	paramNum int
	pattern  string
	format   string
	name     string
	root     *Tree
}

// Leaf is a tree node
type leaf struct {
	kind        uint
	pattern     string
	param       string
	handlers    []HandlerFunc
	children    []*leaf
	childrenNum uint
	paramChild  *leaf
	wideChild   *leaf
	parent      *leaf
	root        *Tree
	nameNode    *Node
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
		t.nodes[i] = newLeaf("/", nil, t)
	}
	t.nameNodes = make(map[string]*Node)
	t.groups = make([]*group, 0)
	t.baa = b
	return t
}

// NewNode create a route node
func NewNode(pattern string, root *Tree) *Node {
	return &Node{
		pattern: pattern,
		root:    root,
	}
}

// newLeaf create a tree leaf
func newLeaf(pattern string, handlers []HandlerFunc, root *Tree) *leaf {
	l := new(leaf)
	l.pattern = pattern
	l.handlers = handlers
	l.root = root
	l.kind = leafKindStatic
	l.children = make([]*leaf, 128)
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

// Match find matched route then returns handlers and name
func (t *Tree) Match(method, pattern string, c *Context) ([]HandlerFunc, string) {
	var i, l int
	var root, nl *leaf
	root = t.nodes[RouterMethods[method]]
	current := root

	for {
		switch current.kind {
		case leafKindStatic:
			// static route
			l = len(current.pattern)
			if l > len(pattern) {
				break
			}
			i := l - 1
			for ; i >= 0; i-- {
				if current.pattern[i] != pattern[i] {
					break
				}
			}
			if i >= 0 {
				break
			}
			if len(pattern) == l || current.children[pattern[l]] != nil ||
				current.paramChild != nil ||
				current.wideChild != nil {
				pattern = pattern[l:]
				root = current
			}
		case leafKindParam:
			// params route
			l = len(pattern)
			if current.childrenNum == 0 {
				i = l
			} else {
				for i = 0; i < l && pattern[i] != '/'; i++ {
				}
			}
			c.SetParam(current.param, pattern[:i])
			pattern = pattern[i:]
			root = current
		case leafKindWide:
			// wide route
			c.SetParam(current.param, pattern)
			pattern = pattern[:0]
		default:
		}

		if len(pattern) == 0 {
			if current.handlers != nil {
				if current.nameNode != nil {
					return current.handlers, current.nameNode.name
				}
				return current.handlers, ""
			}
			if root.paramChild == nil && root.wideChild == nil {
				return nil, ""
			}
		} else {
			// children static route
			if current == root {
				if nl = root.children[pattern[0]]; nl != nil {
					current = nl
					continue
				}
			}
		}

		// param route
		if root.paramChild != nil {
			current = root.paramChild
			continue
		}

		// wide route
		if root.wideChild != nil {
			current = root.wideChild
			continue
		}

		break
	}

	return nil, ""
}

// URLFor use named route return format url
func (t *Tree) URLFor(name string, args ...interface{}) string {
	if name == "" {
		return ""
	}
	node := t.nameNodes[name]
	if node == nil || len(node.format) == 0 {
		return ""
	}
	format := make([]byte, len(node.format))
	copy(format, node.format)
	for i := node.paramNum + 1; i <= len(args); i++ {
		format = append(format, "%v"...)
	}
	return fmt.Sprintf(string(format), args...)
}

// Add registers a new handle with the given method, pattern and handlers.
// add check training slash option.
func (t *Tree) Add(method, pattern string, handlers []HandlerFunc) RouteNode {
	if method == "GET" && t.autoHead {
		t.add("HEAD", pattern, handlers)
	}
	if t.autoTrailingSlash && (len(pattern) > 1 || len(t.groups) > 0) {
		if pattern[len(pattern)-1] == '/' {
			t.add(method, pattern[:len(pattern)-1], handlers)
		} else if pattern[len(pattern)-1] == '*' {
			// wideChild not need trail slash
		} else {
			t.add(method, pattern+"/", handlers)
		}
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
func (t *Tree) add(method, pattern string, handlers []HandlerFunc) RouteNode {
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
	origPattern := pattern
	nameNode := NewNode(origPattern, t)

	// specialy route = /
	if len(pattern) == 1 {
		root.handlers = handlers
		root.nameNode = nameNode
		return nameNode
	}

	// left trim slash, because root is slash /
	pattern = pattern[1:]

	var radix []byte
	var param []byte
	var i, k int
	var tl *leaf
	for i = 0; i < len(pattern); i++ {
		// wide route
		if pattern[i] == '*' {
			// clear static route
			if len(radix) > 0 {
				root = root.insertChild(newLeaf(string(radix), nil, t))
				radix = radix[:0]
			}
			tl = newLeaf("*", handlers, t)
			tl.kind = leafKindWide
			tl.nameNode = nameNode
			root.insertChild(tl)
			break
		}

		// param route
		if pattern[i] == ':' {
			// clear static route
			if len(radix) > 0 {
				root = root.insertChild(newLeaf(string(radix), nil, t))
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
				tl.nameNode = nameNode
			} else {
				tl = newLeaf(":", nil, t)
			}
			tl.param = string(param[:k])
			tl.kind = leafKindParam
			root = root.insertChild(tl)
			continue
		}
		radix = append(radix, pattern[i])
	}

	// static route
	if len(radix) > 0 {
		tl = newLeaf(string(radix), handlers, t)
		tl.nameNode = nameNode
		root.insertChild(tl)
	}

	return nameNode
}

// insertChild insert child into root route, and returns the child route
func (l *leaf) insertChild(node *leaf) *leaf {
	// wide route
	if node.kind == leafKindWide {
		if l.wideChild != nil {
			panic("Router Tree.insert error: cannot set two wide route with same prefix!")
		}
		l.wideChild = node
		return node
	}

	// param route
	if node.kind == leafKindParam {
		if l.paramChild == nil {
			l.paramChild = node
			return l.paramChild
		}
		if l.paramChild.param != node.param {
			panic("Router Tree.insert error cannot use two param [:" + l.paramChild.param + ", :" + node.param + "] with same prefix!")
		}
		if node.handlers != nil {
			if l.paramChild.handlers != nil {
				panic("Router Tree.insert error: cannot twice set handler for same route")
			}
			l.paramChild.handlers = node.handlers
			l.paramChild.nameNode = node.nameNode
		}
		return l.paramChild
	}

	// static route
	child := l.children[node.pattern[0]]
	if child == nil {
		// new child
		l.children[node.pattern[0]] = node
		l.childrenNum++
		return node
	}

	pos := child.hasPrefixString(node.pattern)
	pre := node.pattern[:pos]
	if pos == len(child.pattern) {
		// same route
		if pos == len(node.pattern) {
			if node.handlers != nil {
				if child.handlers != nil {
					panic("Router Tree.insert error: cannot twice set handler for same route")
				}
				child.handlers = node.handlers
				child.nameNode = node.nameNode
			}
			return child
		}

		// child is prefix or node
		node.pattern = node.pattern[pos:]
		return child.insertChild(node)
	}

	newChild := newLeaf(child.pattern[pos:], child.handlers, child.root)
	newChild.nameNode = child.nameNode
	newChild.children = child.children
	newChild.childrenNum = child.childrenNum
	newChild.paramChild = child.paramChild
	newChild.wideChild = child.wideChild

	// node is prefix of child
	if pos == len(node.pattern) {
		child.reset(node.pattern, node.handlers)
		child.nameNode = node.nameNode
		child.children[newChild.pattern[0]] = newChild
		child.childrenNum++
		return child
	}

	// child and node has same prefix
	child.reset(pre, nil)
	child.children[newChild.pattern[0]] = newChild
	child.childrenNum++
	node.pattern = node.pattern[pos:]
	child.children[node.pattern[0]] = node
	child.childrenNum++
	return node
}

// resetPattern reset route pattern and alpha
func (l *leaf) reset(pattern string, handlers []HandlerFunc) {
	l.pattern = pattern
	l.children = make([]*leaf, 128)
	l.childrenNum = 0
	l.paramChild = nil
	l.wideChild = nil
	l.nameNode = nil
	l.param = ""
	l.handlers = handlers
}

// hasPrefixString returns the same prefix position, if none return 0
func (l *leaf) hasPrefixString(s string) int {
	var i, j int
	j = len(l.pattern)
	if len(s) < j {
		j = len(s)
	}
	for i = 0; i < j && s[i] == l.pattern[i]; i++ {
	}
	return i
}

// String returns full pattern of leaf
func (l *leaf) String() string {
	s := l.pattern
	if l.kind == leafKindParam {
		s += l.param
	}
	if l.parent != nil {
		s = l.parent.String() + s
	}
	return s
}

// Name set name of route
func (n *Node) Name(name string) {
	if name == "" {
		return
	}
	p := 0
	f := make([]byte, 0, len(n.pattern))
	for i := 0; i < len(n.pattern); i++ {
		if n.pattern[i] != ':' {
			f = append(f, n.pattern[i])
			continue
		}
		f = append(f, '%')
		f = append(f, 'v')
		p++
		for i = i + 1; i < len(n.pattern); i++ {
			if n.pattern[i] == '/' {
				i--
				break
			}
		}
	}
	n.format = string(f)
	n.paramNum = p
	n.name = name
	n.root.nameNodes[name] = n
}
