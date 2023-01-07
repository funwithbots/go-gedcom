package gedcom

// Document defines a generic interface for interacting with a gedcom document
type Document interface {
	Validate() ([]string, error)
	String() string
	Len() int
}
