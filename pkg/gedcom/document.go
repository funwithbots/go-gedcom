package gedcom

// Document defines a generic interface for interacting with a gedcom document
type Document interface {
	String() string
	Len() int
	ValidateNode(n Node) error
	AddRecord(n *Node)
	Records() []*Node
}
