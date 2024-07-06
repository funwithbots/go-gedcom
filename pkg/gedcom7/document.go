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

	"github.com/funwithbots/go-gedcom/pkg/gedcom"
	"github.com/funwithbots/go-gedcom/pkg/gedcom7/gc70val"
	"github.com/funwithbots/go-gedcom/pkg/stack"
)

const (
	bom = "\xef\xbb\xbf" // byte order mark
)

type document struct {
	bom     string // Byte Order Mark
	header  *gedcom.Node
	records []*gedcom.Node
	trailer *Line

	warnings  []warning
	XRefCache *sync.Map
	Validator *gc70val.Specs

	options docOptions
}

func NewDocumentFromFile(name string, options ...DocOptions) (gedcom.Document, error) {
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
func NewDocument(s *bufio.Scanner, options ...DocOptions) gedcom.Document {
	s.Split(scanLines)

	nodeStack := stack.New()

	doc := &document{
		records:   make([]*gedcom.Node, 0),
		warnings:  make([]warning, 0),
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

	doc.XRefCache.Store(voidXref, &Line{})
	var i int

	for s.Scan() {
		i++
		v := s.Text()
		// Trim BOM from first line
		if i == 1 && strings.HasPrefix(v, bom) {
			v = strings.TrimLeft(v, bom)
			doc.bom = bom
		}
		line, err := ToLine(v) // return a pointer to new Line
		if err != nil {
			doc.AddWarning(v, fmt.Sprintf("Line %d: unable to parse input row '%s'. Error %s", i, v, err.Error()))
			continue
		}
		if !doc.Validator.IsValidTag(line.Tag) {
			doc.AddWarning(line, fmt.Sprintf("Line %d: invalid tag %s", i, line.Tag))
		}

		if errors := line.Validate(); len(errors) > 0 {
			for _, e := range errors {
				doc.AddWarning(e, fmt.Sprintf("Line %d: %s", i, line.Payload))
			}
		}

		switch line.Level {
		case 0:
			node := gedcom.NewNode(nil, line)
			switch line.Tag {
			case "HEAD":
				if i != 1 {
					doc.AddWarning(line, fmt.Sprintf("Line %d: HEAD tag must be first tag in file", i))
				}
				// only add the first header to the document
				if doc.header == nil {
					doc.header = node
				}
				doc.AddRecord(node)
			case "TRLR":
				doc.trailer = line
			default:
				doc.AddRecord(node)
			}
			nodeStack = stack.New()
			nodeStack.Push(node)
		default:
			top := nodeStack.Peek().(*gedcom.Node)
			tLine := top.GetValue().(*Line)
			if line.Level > tLine.Level+1 {
				doc.AddWarning(v, fmt.Sprintf("Level jumped from %d to %d.", tLine.Level, line.Level))
			}
			switch {
			// same level. Just add subnode to parent
			case tLine.Level == line.Level:
				nodeStack.Pop()
				top = nodeStack.Peek().(*gedcom.Node)
				nodeStack.Push(top.AddSubnode(line))
			// deeper level.
			case tLine.Level < line.Level:
				if line.Tag == "CONT" {
					// Special case. append to payload if tag is 'CONT'
					tLine.Payload += "\n" + line.Payload
				} else {
					// Add subnode to most recent subnode
					nodeStack.Push(top.AddSubnode(line))
				}
			// get to matching level on stack then add node.
			case tLine.Level > line.Level:
				for tLine.Level >= line.Level {
					nodeStack.Pop()
					tLine = nodeStack.Peek().(*gedcom.Node).GetValue().(*Line)
				}
				top = nodeStack.Peek().(*gedcom.Node)
				nodeStack.Push(top.AddSubnode(line))
			}
		}
	}

	return doc
}

// AddRecord adds a primary record to the document.
func (d *document) AddRecord(n *gedcom.Node) {
	d.records = append(d.records, n)
	line := n.GetValue().(*Line)
	if line.Xref != "" {
		d.XRefCache.Store(line.Xref, n)
	}
}

// AddWarning adds a warning to the document.
// It is generally populated by the parser.
// Batch operations should check for warnings after processing.
func (d *document) AddWarning(src interface{}, msg string) {
	var row string
	var line *Line
	switch data := src.(type) {
	case string:
		row = data
	case *Line:
		line = data
		row = line.String()
	default:
		row = fmt.Sprintf("unknown src type %T", src)
	}
	d.warnings = append(d.warnings, warning{
		Node:    line,
		Line:    row,
		Message: msg,
	})
}

// GetFamily returns a family Line by xref
func (d *document) GetFamily(xref string) *Line {
	return nil
}

// GetXRef returns a Line by its xref value.
func (d *document) GetXRef(xref string) *Line {
	if line, ok := d.XRefCache.Load(xref); ok {
		return line.(*gedcom.Node).GetValue().(*Line)
	}

	return nil
}

// FindDuplicateIndividuals identifies potential duplicates in a document.
// Birth and death dates can be inexact +/- the provided durations.
// Name can be varied to allow looser matches. e.g. Soundex, different given names, etc.
// Each set of matches are returned in a slice of slices.
func (d *document) FindDuplicateIndividuals(person Line, birthDate, deathDate time.Duration, nameVariations bool) ([][]*Line, error) {
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

// exportNodeTree recursively writes out a Line and its subnodes and returns them as a multi-line string.
func exportNodeTree(n *gedcom.Node) string {
	out := n.GetValue().(*Line).String() + "\n"

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
		line := v.GetValue().(*Line)
		if line.Tag == "TRLR" || line.Tag == "HEAD" {
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

func (d *document) Len() int {
	return len(d.records)
}

func (d *document) Records() []*gedcom.Node {
	return d.records
}
