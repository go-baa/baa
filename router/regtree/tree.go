package regtree

import (
	"regexp"

	"github.com/go-baa/baa"
)

// Tree provlider router store for baa with regexp and radix tree
type Tree struct {
	static    bool
	key       string
	val       interface{}
	params    []string
	format    []byte
	re        *regexp.Regexp
	schildren []*Tree
	rchildren []*Tree
	root      *Tree
}

// NewTree create a new tree route node
func NewTree(key string, val interface{}) *Tree {
	return &Tree{
		static: true,
		key:    key,
		val:    val,
	}
}

// Get returns matched node and param values for key
func (t *Tree) Get(key string) (*Tree, []string) {
	matched := 0
	if t.static {
		for ; matched < len(key) && matched < len(t.key) && key[matched] == t.key[matched]; matched++ {
		}
		// no prefix
		if matched != len(t.key) {
			return nil, nil
		}
		// found
		if matched == len(key) {
			if t.val != nil {
				return t, nil
			}
		}
		// node is prefix
		key = key[matched:]
		// first, static rule
		if len(key) > 0 {
			for i := range t.schildren {
				if n, v := t.schildren[i].Get(key); n != nil {
					return n, v
				}
			}
		}

		// then, regexp rule
		for i := range t.rchildren {
			data := t.rchildren[i].re.FindStringSubmatch(key)
			if len(data) != len(t.rchildren[i].params)+1 || data[0] != key {
				continue
			}
			return t.rchildren[i], data[1:]
		}
		return nil, nil
	}

	// regexp node
	data := t.re.FindStringSubmatch(key)
	if len(data) != len(t.params)+1 || data[0] != key {
		return nil, nil
	}
	return t, data[1:]
}

// Add return new node with key and val
func (t *Tree) Add(key string, val interface{}) *Tree {
	// find the common prefix
	matched := 0
	for ; matched < len(key) && matched < len(t.key) && key[matched] == t.key[matched]; matched++ {
	}

	// no prefix
	if matched == 0 {
		return nil
	}

	if matched == len(t.key) {
		// the node key is the same as the key: make the current node as data node
		if matched == len(key) {
			if val != nil {
				if t.val != nil {
					panic("the route is be exists: " + t.String())
				}
				t.val = val
			}
			return t
		}

		// the node key is a prefix of the key: create a child node
		key = key[matched:]
		for _, child := range t.schildren {
			if node := child.Add(key, val); node != nil {
				return node
			}
		}

		// no child match, to be a new child
		return t.addChild(key, val)
	}

	// the key is a prefix of node key: create a new node instead of child
	if matched == len(key) {
		node := NewTree(t.key[matched:], t.val)
		node.schildren = t.schildren
		node.rchildren = t.rchildren
		node.root = t
		t.key = key
		t.val = val
		t.schildren = []*Tree{node}
		t.rchildren = nil
		return t
	}

	// the node key shares a partial prefix with the key: split the node key
	node := NewTree(t.key[matched:], t.val)
	node.schildren = t.schildren
	node.rchildren = t.rchildren
	node.root = t
	t.key = key[:matched]
	t.val = nil
	t.schildren = nil
	t.rchildren = nil
	t.schildren = append(t.schildren, node)
	return t.addChild(key[matched:], val)
}

func (t *Tree) addChild(key string, val interface{}) *Tree {
	// check it is a static route child or not
	var staticKey, param, rule string
	var params []string
	var newKey, format []byte
	var i, j, k int
	for i = 0; i < len(key); i++ {
		if key[i] == '*' {
			// set static prefix
			if len(staticKey) == 0 && len(params) == 0 && i > 0 {
				staticKey = key[:i]
			}
			rule = "(.*)"
			param = ""
			newKey = append(newKey, rule...)
			format = append(format, "%v"...)
			params = append(params, param)
			continue
		}
		if key[i] == ':' {
			for j = i + 1; j < len(key) && baa.IsParamChar(key[j]); j++ {
			}
			// set static prefix
			if len(staticKey) == 0 && len(params) == 0 && i > 0 {
				staticKey = key[:i]
			}
			param = key[i+1 : j]
			i = j - 1
			// check regexp rule
			rule = ""
			if j < len(key) && key[j] == '(' {
				for k = j + 1; k < len(key) && key[k] != ')'; k++ {
				}
				rule = key[j+1 : k]
				i = k
			}
			if rule == "" {
				rule = "([^\\/]+)"
			} else if rule == "int" {
				rule = "([\\d]+)"
			} else if rule == "string" {
				rule = "([\\w]+)"
			} else {
				rule = "(" + rule + ")"
			}
			newKey = append(newKey, rule...)
			format = append(format, "%v"...)
			params = append(params, param)
			continue
		}
		newKey = append(newKey, key[i])
		format = append(format, key[i])
	}

	if len(params) > 0 {
	}

	var reNode, staticNode *Tree
	var err error
	if len(params) > 0 {
		// key has regexp rule, new regexp rule
		reNode = NewTree(string(newKey[len(staticKey):]), val)
		reNode.static = false
		reNode.params = params
		reNode.format = format
		reNode.re, err = regexp.Compile(reNode.key + "$")
		if err != nil {
			panic("tree.addChild: " + err.Error())
		}
		// set key with static prefix
		key = staticKey
	}

	if len(key) > 0 {
		// key has static rule
		staticNode = NewTree(key, nil)
		staticNode.root = t
		if reNode != nil {
			reNode.root = staticNode
			staticNode.rchildren = append(staticNode.rchildren, reNode)
			t.schildren = append(t.schildren, staticNode)
			return reNode
		}
		staticNode.val = val
		staticNode.format = format
		t.schildren = append(t.schildren, staticNode)
		return staticNode
	}

	// key has regexp rule without static rule
	reNode.root = t
	for _, child := range t.rchildren {
		if child.key == reNode.key {
			panic("the route is be exists: " + child.String())
		}
	}
	t.rchildren = append(t.rchildren, reNode)
	return reNode
}

// String return full key
func (t *Tree) String() string {
	s := t.key
	if t.root != nil {
		s = t.root.String() + s
	}
	return s
}
