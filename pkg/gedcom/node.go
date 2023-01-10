package gedcom

import (
	"fmt"
	"sync"
)

// Node represents a tree of subnodes
// Rather than being an interface, the Node struct can be used with any tree of subnodes.

type (
	Node struct {
		value    interface{}
		subnodes []*Node
		parent   *Node
	}
)

// NewNode creates a new node
func NewNode(in interface{}) *Node {
	n := &Node{
		value:    in,
		subnodes: make([]*Node, 0),
	}
	return n
}

// AddSubnode appends a new subnode to the current node.
func (n *Node) AddSubnode(in interface{}) *Node {
	nod := NewNode(in)
	n.subnodes = append(n.subnodes, nod)

	return n
}

// RemoveSubnode removes a subnode from the current node.
// No changes are made to the subnode's children if the subnode does not exist.
func (n *Node) RemoveSubnode(in interface{}) {
	for i, v := range n.subnodes {
		if v.value == in {
			n.subnodes = append(n.subnodes[:i], n.subnodes[i+1:]...)
			return
		}
	}
}

// UpdateSubnode replaces the contents of a subnode with a new value.
// It doesn't change the tree structure.
func (n *Node) UpdateSubnode(old, new interface{}) error {
	for _, v := range n.subnodes {
		if v.value == old {
			v.value = new
			return nil
		}
	}

	return fmt.Errorf("subnode not found")
}

// GetNode returns the current node.
func (n *Node) GetNode() interface{} {
	return n.value
}

func (n *Node) GetParent() interface{} {
	return n.parent
}

func (n *Node) GetSubnodes() []*Node {
	return n.subnodes
}

// ProcessTree applies fn to each node in the tree starting with the current node.
func (n *Node) ProcessTree(fn func(*Node) interface{}, wg *sync.WaitGroup) []interface{} {
	defer wg.Done()
	if n == nil {
		return nil
	}
	out := make([]interface{}, 0)
	out = append(out, fn(n))
	for _, v := range n.subnodes {
		out = append(out, v.ProcessSubnodes(fn)...)
	}
	return out
}

func (n *Node) ProcessSubnodes(fn func(*Node) interface{}) []interface{} {
	if n == nil {
		return nil
	}
	out := make([]interface{}, 0)
	out = append(out, fn(n))
	for _, v := range n.subnodes {
		if vv := v.ProcessSubnodes(fn); vv != nil {
			out = append(out, vv...)
		}
	}
	return out
}

func (n *Node) GetValue() interface{} {
	return n.value
}
