package gedcom7

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

/*
Line = Level D [Xref D] Tag [D LineVal] EOL
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
	voidXref     = "@VOID@"
	g7Banned     = `\x{0}-\x{8}\x{B}-\x{C}\x{E}-\x{1F}\x{7F}\x{80}-\x{9F}\x{D800}-\x{DFFF}\x{FFFE}-\x{FFFF}`

	setFixAtsignStartLineVal = true
	setFixAtsignAllLineVal   = false // no longer standard. @@ should not be used
)

var (
	regXref   = regexp.MustCompile(fmt.Sprintf("%s[%s%s%s]+%s", g7Atsign, g7Underscore, g7UCLetter, g7Digit, g7Atsign))
	regLevel  = regexp.MustCompile(fmt.Sprintf("[0%s]+", g7Nonzero))
	regTag    = regexp.MustCompile(fmt.Sprintf("^%s?[%s%s]{1,}$", g7Underscore, g7UCLetter, g7Digit))
	regBanned = regexp.MustCompile(fmt.Sprintf("[%s]+", g7Banned))
)

// node defines the structure of a v7 gedcom line
// Special case if Tag is CONT, then the payload is a continuation of the previous line.
type node struct {
	Level   int
	Xref    string
	Tag     string
	Payload string

	// Deleted reflects if a node is flagged to be deleted.
	Deleted bool
	//
	// Parent   *node
	// Subnodes []*node
}

type GTime node

// String creates a text row compatible with the Gedcom 7.x specification.
// If line breaks exist in the payload, continuation lines are also generated.
func (n *node) String() string {
	if n.Tag == "" || n.Level < 0 {
		return ""
	}
	lines := strings.Split(n.Payload, "\n")
	out := fmt.Sprintf("%s", strconv.Itoa(n.Level))

	if n.Xref != "" {
		out += " " + n.Xref
	}

	if n.Tag != "" {
		out += " " + string(n.Tag)
	}

	if strings.TrimSpace(n.Payload) != "" {
		out += " " + lines[0]
	}

	if len(lines) > 1 {
		level := n.Level + 1
		for _, v := range lines[1:] {
			out += fmt.Sprintf("\n%d CONT %s", level, v)
		}
	}

	return out
}

// ToNode converts a gedcom file line to a Node structure and returns it.
func ToNode(s string) (*node, error) {
	var (
		N      node
		err    error
		marker = 1
	)

	// TODO Need to return a warning and still process the line. Let caller decide to throw it away.
	if regBanned.MatchString(s) {
		return nil, errors.New("line contains banned characters")
	}

	tokens := strings.Split(strings.TrimSpace(s), " ")
	maxIndex := len(tokens) - 1
	if maxIndex < 1 {
		return nil, errors.New("not enough tokens to parse line")
	}

	N.Level, err = strconv.Atoi(tokens[0])
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
		N.Xref = t
		marker++
	case regTag.MatchString(t):
		N.Tag = t
		marker++
	default:
		return nil, errors.New("missing or malformed xref/tag in second position")
	}

	if N.Tag == "" {
		if !regTag.MatchString(tokens[marker]) {
			return nil, errors.New("missing tag")
		}
		N.Tag = tokens[marker]
		marker++
	}

	skip := len(N.Xref) + 2
	n := strings.Index(s[skip:], string(N.Tag)) + len(N.Tag) + 1 + skip
	if len(s) > n {
		lv := s[n:]
		if setFixAtsignStartLineVal && lv[0:1] == "@" && !isXref(lv) {
			if len(lv) > 1 && lv[:2] != "@@" {
				lv = fmt.Sprintf("@%s", lv)
			}
		}
		N.Payload = lv
	}

	return &N, nil
}

// isXref validates that string segment is a validly formatted XRef.
func isXref(v string) bool {
	return regXref.MatchString(v)
}
