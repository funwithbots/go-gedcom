package gedcom

import "fmt"

// Node represents a tree of subnodes
// Rather than being an interface, the Node struct can be used with any tree of subnodes.

type (
	Node struct {
		value    interface{}
		subnodes []*Node
		parent   *Node
	}

	// node struct {
	// 	value  interface{}
	// 	parent *Node
	// }
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
func (r *Node) AddSubnode(in interface{}) *Node {
	n := NewNode(in)
	r.subnodes = append(r.subnodes, n)

	return n
}

// RemoveSubnode removes a subnode from the current node.
// No changes are made to the subnode's children if the subnode does not exist.
func (r *Node) RemoveSubnode(in interface{}) {
	for i, v := range r.subnodes {
		if v.value == in {
			r.subnodes = append(r.subnodes[:i], r.subnodes[i+1:]...)
			return
		}
	}
}

// UpdateSubnode replaces the contents of a subnode with a new value.
// It doesn't change the tree structure.
func (r *Node) UpdateSubnode(old, new interface{}) error {
	for _, v := range r.subnodes {
		if v.value == old {
			v.value = new
			return nil
		}
	}

	return fmt.Errorf("subnode not found")
}

// GetNode returns the current node.
func (r *Node) GetNode() interface{} {
	return r.value
}

func (r *Node) GetParent() interface{} {
	return r.parent
}

func (r *Node) GetSubnodes() []*Node {
	return r.subnodes
}
