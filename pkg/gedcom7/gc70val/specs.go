package gc70val

import (
	"fmt"

	"github.com/funwithbots/go-abnf/operators"
)

type Specs struct {
	StdTags   map[string]tagDef
	ExtTags   map[string]tagDef
	Calendars map[string]calDef
	Types     map[string]typeDef

	LineValidator operators.Operator

	// deprecated tags from previous gedcom versions
	depTags map[string]bool
}

func New() *Specs {
	return &Specs{
		StdTags:   baseline.tags,
		ExtTags:   make(map[string]tagDef),
		Calendars: baseline.calendars,
		Types:     baseline.types,
		depTags:   make(map[string]bool),
	}
}

// AddExtTag adds a tag to the list of extension tags
func (g *Specs) AddExtTag(in []byte) error {
	tag, err := loadTag(in)
	if err != nil {
		return err
	}
	if _, ok := g.ExtTags[tag.Tag]; ok {
		return fmt.Errorf("tag already exists as an extension tag")
	}
	if _, ok := g.StdTags[tag.Tag]; ok {
		return fmt.Errorf("tag already exists as standard tag")
	}
	g.ExtTags[tag.Tag] = tag

	return nil
}

func (g *Specs) SetDeprecatedTags(v string) {
	for tag, ver := range deprecatedTags {
		if ver <= v {
			g.depTags[tag] = true
		}
	}
}
