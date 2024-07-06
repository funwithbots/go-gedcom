package gc70val

import (
	_ "embed"
	"errors"
	"fmt"
	"strings"

	"github.com/funwithbots/go-abnf/operators"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/funwithbots/go-gedcom/pkg/abnf"
)

const (
	TagCONT = "CONT"
	TagHEAD = "HEAD"
	TagTRLR = "TRLR"
)

var validTags = make(map[string]interface{})

// AddValidTag adds a tag to the list of valid tags.
// If it already exists, return false, otherwise true.
func AddValidTag(tag string) bool {
	if IsValidTag(tag) {
		return false
	}
	validTags[tag] = true
	return true
}

func IsValidTag(tag string) bool {
	_, ok := validTags[tag]
	return ok
}

// TODO Implement this with validation of gedcom documents
// type cardinality int
//
// const (
//
//	CardinalityZeroToOne cardinality = iota
//	CardinalityZeroToMany
//	CardinalityOneToOne
//	CardinalityOneToMany
//
// )
//
//	var cardinalityMap = map[string]cardinality{
//		"{0:1}": CardinalityZeroToOne,
//		"{0:M}": CardinalityZeroToMany,
//		"{1:1}": CardinalityOneToOne,
//		"{1:M}": CardinalityOneToMany,
//	}

var (
	URIPrefixes = map[string]string{
		"g7":   "https://gedcom.io/terms/v7/",
		"xsd":  "http://www.w3.org/2001/XMLSchema#",
		"dcat": "http://www.w3.org/ns/decat#",
	}
)

var pseudoTags = map[string]TagDef{
	"HEAD": {
		Lang:    "en-US",
		Type:    tagTypePseudostructure,
		URI:     "/HEAD",
		Tag:     "HEAD",
		FullTag: "HEAD",
		Specification: []string{
			"The header pseudo-structure provides metadata about the entire dataset.",
		},
		Substructures: map[string]string{
			"GEDC":  "{1:1}",
			"SCHMA": "{0:1}",
			"SOUR":  "{0:1}",
			"DEST":  "{0:1}",
			"DATE":  "{0:1}",
			"SUBM":  "{0:1}",
			"COPR":  "{0:1}",
			"LANG":  "{0:1}",
			"NOTE":  "{0:1}",
		},
		Rule: abnf.Null(),
	},
	"TRLR": {
		Lang:    "en-US",
		Type:    tagTypePseudostructure,
		URI:     "/TRLR",
		Tag:     "TRLR",
		FullTag: "TRLR",
		Specification: []string{
			"The trailer resembles a record, comes last in each document, and cannot contain substructures.",
		},
		Rule: abnf.Null(),
	},
	"CONT": {
		Lang:    "en-US",
		Type:    tagTypePseudostructure,
		URI:     "/CONT",
		Tag:     "CONT",
		FullTag: "CONT",
		Specification: []string{
			"A line continuation resembles a substructure, comes before any other substruc‐",
			"tures, is used to encode multi-line payloads, and cannot contain substructures.",
		},
		Payload: "http://www.w3.org/2001/XMLSchema#string",
		Rule:    abnf.Validation["String"],
	},
}

type TagDef struct {
	Lang          string `yaml:"lang"`
	Type          string `yaml:"type"`
	URI           string `yaml:"uri"`
	Tag           string `yaml:"standard tag"`
	FullTag       string
	Specification []string `yaml:"specification"`
	Payload       string   `yaml:"payload"`

	Substructures   map[string]string `yaml:"substructures"`
	Superstructures map[string]string `yaml:"superstructures"`
	ValueOf         []string          `yaml:"value of"`

	// Validation func based on ABNF spec
	Rule operators.Operator

	EnumSetName string `yaml:"enumeration set"`
	EnumSet     enumSet
}

func (t *TagDef) SetRule(rule operators.Operator) {
	t.Rule = rule
}

// InferRule extracts the tag's validation rule from the ABNF Validation map.
func (t *TagDef) InferRule() {
	if t.URI != "" {
		switch {
		case len(t.Payload) > 0 && t.Payload[0] == '@':
			t.Rule = abnf.Validation["Xref"]
		case t.Payload == "null", t.Payload == "":
			t.Rule = abnf.Validation["Null"]
		case len(t.Payload) > 2 && t.Payload[0] == '@' && t.Payload[len(t.Payload)-1] == '@':
			t.Rule = abnf.Validation["Tag"]
		case strings.HasPrefix(t.Payload, URIPrefixes["g7"]):
			k := t.Payload[len(URIPrefixes["g7"]):]
			if k[:5] == "type-" {
				k = k[5:]
			}
			k = toUpperCamel(k)
			t.Rule = abnf.Validation[k]
		case strings.HasPrefix(t.Payload, URIPrefixes["xsd"]),
			strings.HasPrefix(t.Payload, URIPrefixes["w3"]):
			parts := strings.Split(t.Payload, "#")
			uc := toUpperCamel(parts[len(parts)-1])
			t.Rule = abnf.Validation[uc]
		default:
			fmt.Printf(t.Tag)
		}
	}
	if t.Rule == nil {
		fmt.Printf(t.Tag)
	}
}

// ValidatePayload checks if the tag's payload conforms to the ABNF definition.
func (t *TagDef) ValidatePayload(str string) bool {
	// If it's an enum, check that first
	if len(t.EnumSet.Values) > 0 {
		for _, v := range t.EnumSet.Values {
			if v == str {
				return true
			}
		}
		return false
	}

	// If no validation rule exists, return true
	if t.Rule == nil {
		return true
	}
	v := t.Rule([]byte(str))

	// Check if it's an Xref.
	if len(v) == 0 {
		return IsXref(str)
	}

	return str == v[0].String()
}

// loadTag returns a TagDef for the provided yaml.
func loadTag(in []byte) (TagDef, error) {
	td := TagDef{}

	if err := deserializeYAML(in, &td); err != nil {
		return td, err
	}
	if td.Tag == "" {
		ss := strings.Split(td.URI, "/")
		if len(ss) == 0 {
			return TagDef{}, errors.New("no standard tag")
		}
		td.Tag = ss[len(ss)-1]
	}
	td.FullTag = extractFullTag(td.URI)
	AddValidTag(td.FullTag)
	td.InferRule()
	td.Superstructures = remapStructures(td.Superstructures)
	td.Substructures = remapStructures(td.Substructures)

	return td, nil
}

// remapStructures removes the URI prefix from the structure keys
func remapStructures(structs map[string]string) map[string]string {
	s := make(map[string]string)
	for k, v := range structs {
		key := extractFullTag(k)
		s[key] = v
	}

	return s
}

func extractFullTag(u string) string {
	parts := strings.Split(u, "/")
	key := parts[len(parts)-1]
	last := strings.Split(key, "-")
	if len(last) >= 1 && strings.ToLower(last[0]) == last[0] {
		key = strings.Join(last[1:], "-")
	}

	return key
}

// toUpperCamel converts strings with hyphens and underscores to UpperCamelCase
func toUpperCamel(in string) string {
	caser := cases.Title(language.English)
	var out string
	mid := make([]string, 0)
	parts := strings.Split(in, "-")
	for _, part := range parts {
		mid = append(mid, strings.Split(part, "_")...)
	}
	for _, v := range mid {
		out += caser.String(v)
	}

	return out
}
