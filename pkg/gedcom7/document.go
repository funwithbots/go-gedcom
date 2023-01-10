package gedcom7

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"go-gedcom/pkg/gedcom"
	"go-gedcom/pkg/gedcom7/gc70val"
	"go-gedcom/pkg/stack"
)

const (
	bom         = "\xef\xbb\xbf" // byte order mark
	abnfDocPath = "data/"
)

type document struct {
	bom     string // Byte Order Mark
	header  *gedcom.Node
	records []*gedcom.Node
	trailer *Node

	Warnings  []Warning
	XRefCache *sync.Map
	Validator *gc70val.Specs

	options docOptions
}

func NewDocumentFromFile(name string, options ...DocOptions) (*document, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	doc := NewDocument(s, WithMaxDeprecatedTags("5.5.1"))
	return doc, nil
}

// NewDocument accepts a buffer and converts it into a structured gedcom document.
func NewDocument(s *bufio.Scanner, options ...DocOptions) *document {
	s.Split(scanLines)

	nodeStack := stack.New()

	doc := &document{
		records:   make([]*gedcom.Node, 0),
		Warnings:  make([]Warning, 0),
		XRefCache: &sync.Map{},
		Validator: gc70val.New(),
		options: docOptions{
			docPath:                 "./",
			maxDeprecatedTagVersion: "ZZZ",
		},
	}

	opts := docOptions{}
	opts.withOpts(options)
	if opts.maxDeprecatedTagVersion != "" {
		doc.Validator.SetDeprecatedTags(opts.maxDeprecatedTagVersion)
	}

	doc.XRefCache.Store(voidXref, &Node{})
	var i int

	for s.Scan() {
		i++
		v := s.Text()
		// Trim BOM from first line
		if i == 1 && strings.HasPrefix(v, bom) {
			v = strings.TrimLeft(v, bom)
			doc.bom = bom
		}
		nod, err := ToNode(v) // return a pointer to new Node
		rec := gedcom.NewNode(nod)
		if err != nil {
			doc.AddWarning(v, fmt.Sprintf("unable to parse input row %d: %s", i, v))
			continue
		}
		if !doc.Validator.IsValidTag(nod.Tag) {
			doc.AddWarning(nod, fmt.Sprintf("invalid tag %s on line %d", nod.Tag, i))
		}
		switch nod.Level {
		case 0:
			switch nod.Tag {
			case "HEAD":
				doc.header = rec
				doc.AddRecord(rec)
			case "TRLR":
				doc.trailer = nod
			default:
				doc.AddRecord(rec)
			}
			nodeStack.Push(rec)
			// parent = nod
		default:
			last := nodeStack.Peek().(*gedcom.Node)
			parent := last.GetNode().(*Node)
			if nod.Level > parent.Level+1 {
				doc.AddWarning(v, fmt.Sprintf("Level jumped from %d to %d.", parent.Level, nod.Level))
			}
			switch {
			// same level. Just add subnode to parent
			case parent.Level == nod.Level:
				sn := last.AddSubnode(nod)
				nodeStack.Push(sn)
			// deeper level. Add subnode to last subnode
			// Special case. append to payload if tag is 'CONT'
			case parent.Level < nod.Level:
				if nod.Tag == "CONT" {
					parent.Payload += "\n" + nod.Payload
				} else {
					sn := last.AddSubnode(nod)
					nodeStack.Push(sn)
				}
			// get to matching level on stack then add node.
			case parent.Level > nod.Level:
				for parent.Level >= nod.Level {
					nodeStack.Pop()
					last = nodeStack.Peek().(*gedcom.Node)
					parent = last.GetNode().(*Node)
				}
				sn := last.AddSubnode(nod)
				nodeStack.Push(sn)
			}
		}
	}

	return doc
}

// AddRecord adds a primary record to the document.
func (d *document) AddRecord(n *gedcom.Node) {
	d.records = append(d.records, n)
	nod := n.GetNode().(*Node)
	if nod.Xref != "" {
		d.XRefCache.Store(nod.Xref, n)
	}
}

// AddWarning adds a warning to the document.
// It is generally populated by the parser.
// Batch operations should check for warnings after processing.
func (d *document) AddWarning(src interface{}, msg string) {
	var row string
	var n *Node
	switch data := src.(type) {
	case string:
		row = data
	case *Node:
		n = data
		row = n.String()
	default:
		row = fmt.Sprintf("unknown src type %T", src)
	}
	d.Warnings = append(d.Warnings, Warning{
		Node:    n,
		Line:    row,
		Message: msg,
	})
}

// GetWarnings returns a slice of warnings.
func (d *document) GetWarnings() []Warning {
	return d.Warnings
}

// GetFamily returns a family Node by xref
func (d *document) GetFamily(xref string) *Node {
	return nil
}

// GetXRef returns a Node by its xref value.
func (d *document) GetXRef(xref string) *Node {
	if n, ok := d.XRefCache.Load(xref); ok {
		return n.(*gedcom.Node).GetNode().(*Node)
	}

	return nil
}

// FindDuplicateIndividuals identifies potential duplicates in a document.
// Birth and death dates can be inexact +/- the provided durations.
// Name can be varied to allow looser matches. e.g. Soundex, different given names, etc.
// Each set of matches are returned in a slice of slices.
func (d *document) FindDuplicateIndividuals(person Node, birthDate, deathDate time.Duration, nameVariations bool) ([][]*Node, error) {
	return nil, fmt.Errorf("not implemented")
}

func (d *document) ExportToFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	if err = d.exportGedcom7(w); err != nil {
		return err
	}

	w.Flush()
	return nil
}

// exportNodeTree recursively writes out a Node and its subnodes and returns them as a multi-line string.
func exportNodeTree(n *gedcom.Node) string {
	out := n.GetNode().(*Node).String() + "\n"

	for _, v := range n.GetSubnodes() {
		out += exportNodeTree(v)
	}

	return out
}

// exportGedcom7 exports a document to a writer in GEDCOM 7 format.
func (d *document) exportGedcom7(w io.Writer) error {
	_, err := w.Write([]byte(d.bom))
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(exportNodeTree(d.header)))
	if err != nil {
		return err
	}

	for _, v := range d.records {
		n := v.GetNode().(*Node)
		if n.Tag == "TRLR" || n.Tag == "HEAD" {
			continue
		}
		_, err = w.Write([]byte(exportNodeTree(v)))
		if err != nil {
			return err
		}
	}

	_, err = w.Write([]byte(d.trailer.String() + "\n"))
	if err != nil {
		return err
	}

	return nil
}

func (d *document) String() string {
	buf := new(bytes.Buffer)
	if err := d.exportGedcom7(buf); err != nil {
		return ""
	}

	return buf.String()
}

func (d *document) Validate() ([]string, error) {
	return nil, fmt.Errorf("not implemented")
}

func (d *document) Len() int {
	return len(d.records)
}

func (d *document) Records() []*gedcom.Node {
	return d.records
}
