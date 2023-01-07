package gc70val

import (
	"go-gedcom/pkg/gedcom"
)

func (g *Specs) ValidateNode(*gedcom.Node) (bool, error) {
	return true, nil
}

func (g *Specs) ValidateDocument() (bool, error) {
	return true, nil
}

// IsValidTag returns true if the tag is a valid standard or extended tag
func (g *Specs) IsValidTag(tag string) bool {
	if _, ok := g.StdTags[tag]; ok {
		return true
	}
	if _, ok := g.StdTags["record-"+tag]; ok {
		return true
	}
	if _, ok := g.ExtTags[tag]; ok {
		return true
	}
	if tag[0] == '_' {
		return true
	}
	if v, ok := g.StdTags["FAM-"+tag]; ok {
		if v.Tag == tag {
			return true
		}
	}
	if v, ok := g.StdTags["INDI-"+tag]; ok {
		if v.Tag == tag {
			return true
		}
	}

	// weird special cases
	if tag == "STAT" {
		return true
	}

	if _, ok := g.depTags[tag]; ok {
		return true
	}

	return false
}
