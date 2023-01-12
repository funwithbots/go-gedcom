package gedcom7

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/*
Text = Level D [Xref D] Tag [D LineVal] EOL
Level = "0" / nonzero *digit
D = %x20 ; space
Xref = atsign 1*tagchar atsign ; but not "@VOID@"
Tag = stdTag / extTag
LineVal = pointer / lineStr
EOL = %x0D [%x0A] / %x0A ; CR-LF, CR, or LF
stdTag = ucletter *tagchar
extTag = underscore 1*tagchar
tagchar = ucletter / digit / underscore
pointer = voidPtr / Xref
voidPtr = %s"@VOID@"
nonAt = %x09 / %x20-3F / %x41-10FFFF ; non-EOL, non-@
nonEOL = %x09 / %x20-10FFFF ; non-EOL
lineStr = (nonAt / atsign atsign) *nonEOL ; leading @ doubled
*/

const (
	g7UCLetter   = "A-Z"
	g7Digit      = "0-9"
	g7Nonzero    = "1-9"
	g7Underscore = "_"
	g7Atsign     = "@"
	g7Banned     = `\x{0}-\x{8}\x{B}-\x{C}\x{E}-\x{1F}\x{7F}\x{80}-\x{9F}\x{D800}-\x{DFFF}\x{FFFE}-\x{FFFF}`

	voidXref = "@VOID@"

	setFixAtsignStartLineVal = true
	setFixAtsignAllLineVal   = false // no longer standard. @@ should not be used
)

var (
	regXref = regexp.MustCompile(fmt.Sprintf("%s[%s%s%s]+%s", g7Atsign, g7Underscore, g7UCLetter, g7Digit, g7Atsign))
	// regLevel  = regexp.MustCompile(fmt.Sprintf("[0%s]+", g7Nonzero))
	regTag    = regexp.MustCompile(fmt.Sprintf("^%s?[%s%s]{1,}$", g7Underscore, g7UCLetter, g7Digit))
	regBanned = regexp.MustCompile(fmt.Sprintf("[%s]+", g7Banned))
)

// Line defines the structure of a v7 gedcom line
// Special case if Tag is CONT, then the payload is a continuation of the previous line.
type Line struct {
	Level   int
	Xref    string
	Tag     string
	Payload string

	// Deleted reflects if a Line is flagged to be deleted.
	Deleted bool

	// Text is the original line of text extracted from the gedcom file.
	Text string
}

type GTime Line

type Lines []Line

// String creates a text row compatible with the Gedcom 7.x specification.
// If line breaks exist in the payload, continuation lines are also generated.
func (l *Line) String() string {
	if l.Tag == "" || l.Level < 0 {
		return ""
	}
	lines := strings.Split(l.Payload, "\n")
	out := strconv.Itoa(l.Level)

	if l.Xref != "" {
		out += " " + l.Xref
	}

	if l.Tag != "" {
		out += " " + l.Tag
	}

	if strings.TrimSpace(l.Payload) != "" {
		out += " " + lines[0]
	}

	if len(lines) > 1 {
		level := l.Level + 1
		for _, v := range lines[1:] {
			out += fmt.Sprintf("\n%d CONT %s", level, v)
		}
	}

	return out
}

// Matches evaluates each line.Text using fn() with `pattern` and returns a slice of matching lines.
// e.g. res := src.Matches(" INDI ", strings.Contains())
func (ls *Lines) Matches(pattern string, fn func(string, string) bool) *Lines {
	var out Lines
	for _, v := range *ls {
		if fn(v.Text, pattern) {
			out = append(out, v)
		}
	}
	return &out
}

// ToLine converts a gedcom file line to a Line structure and returns it.
func ToLine(s string) (*Line, error) {
	var (
		node   Line
		err    error
		marker = 1
	)

	node.Text = s

	// TODO Need to return a warning and still process the line. Let caller decide to throw it away.
	if regBanned.MatchString(s) {
		return nil, errors.New("line contains banned characters")
	}

	tokens := strings.Split(strings.TrimSpace(s), " ")
	maxIndex := len(tokens) - 1
	if maxIndex < 1 {
		return nil, errors.New("not enough tokens to parse line")
	}

	node.Level, err = strconv.Atoi(tokens[0])
	if err != nil {
		return nil, err
	}

	// must be Xref or Tag
	t := tokens[marker]
	switch {
	case regXref.MatchString(t):
		if t == voidXref {
			return nil, errors.New("xref cannot be void")
		}
		node.Xref = t
		marker++
	case regTag.MatchString(t):
		node.Tag = t
		marker++
	default:
		return nil, errors.New("missing or malformed xref/tag in second position")
	}

	if node.Tag == "" {
		if !regTag.MatchString(tokens[marker]) {
			return nil, errors.New("missing tag")
		}
		node.Tag = tokens[marker]
	}

	skip := len(node.Xref) + 2
	n := strings.Index(s[skip:], string(node.Tag)) + len(node.Tag) + 1 + skip
	if len(s) > n {
		lv := s[n:]
		if setFixAtsignStartLineVal && lv[0:1] == "@" && !isXref(lv) {
			if len(lv) > 1 && lv[:2] != "@@" {
				lv = fmt.Sprintf("@%s", lv)
			}
		}
		node.Payload = lv
	}

	return &node, nil
}

// NewLinesFromFile creates a slice of Line structures from a file.
func NewLinesFromFile(name string, options ...DocOptions) (*Lines, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	return NewLines(s), nil
}

// NewLines accepts a buffer and converts it into a slice of parsed Lines.
// Other than parsing the input string, no validation is performed.
func NewLines(s *bufio.Scanner) *Lines {
	s.Split(scanLines)

	var lines Lines

	first := true
	for s.Scan() {
		v := s.Text()
		// Trim BOM from first line
		if first && strings.HasPrefix(v, bom) {
			v = strings.TrimLeft(v, bom)
			first = false
		}
		line, err := ToLine(v) // return a pointer to new Line
		if err != nil {
			line = &Line{
				Level:   -1,
				Tag:     "ERROR",
				Deleted: true,
				Text:    "v",
			}
		}

		lines = append(lines, *line)
	}

	return &lines
}

// isXref validates that string segment is a validly formatted XRef.
func isXref(v string) bool {
	return regXref.MatchString(v)
}
