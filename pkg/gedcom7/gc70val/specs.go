package gc70val

import (
	"fmt"

	"github.com/funwithbots/go-abnf/operators"
)

type Specs struct {
	Tags      map[string]TagDef
	Calendars map[string]calDef
	Types     map[string]typeDef

	// deprecated tags from previous gedcom versions
	depTags map[string]bool
}

func New() *Specs {
	return &Specs{
		Tags:      baseline.tags,
		Calendars: baseline.calendars,
		Types:     baseline.types,
		depTags:   make(map[string]bool),
	}
}

// AddExtTag adds a tag to the list of extension tags
func (s *Specs) AddExtTag(in []byte) error {
	tag, err := loadTag(in)
	if err != nil {
		return err
	}
	if _, ok := s.Tags[tag.Tag]; ok {
		return fmt.Errorf("tag already exists: %s", tag.Tag)
	}
	tag.Type = tagTypeExtended
	s.Tags[tag.Tag] = tag

	return nil
}

// SetDeprecatedTags sets the list of deprecated tags based on the provided GEDCOM version.
func (s *Specs) SetDeprecatedTags(v string) {
	for tag, ver := range deprecatedTags {
		if ver <= v {
			s.depTags[tag] = true
		}
	}
}

// GetRule gets the rule definition for this tag.
func (s *Specs) GetRule(k string) operators.Operator {
	if tag, ok := s.Tags[k]; ok {
		return tag.Rule
	}

	return nil
}

// GetEnums gets the list of Enums for this tag.
func (s *Specs) GetEnums(k string) []string {
	if tag, ok := s.Tags[k]; ok {
		return tag.EnumSet.Values
	}

	return nil
}
