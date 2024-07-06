package gedcom

// Document defines a generic interface for interacting with a gedcom document
type Document interface {
	Validate() ([]string, error)
	String() string
	Warnings() []string
	Len() int
	Records() []*Node
}
