package gedcom

import "github.com/funwithbots/go-abnf/operators"

type Validator interface {
	AddExtTag(in []byte) error
	SetDeprecatedTags(ver string)
	GetRule(tag string) operators.Operator
	GetEnums(tag string) []string
	IsValidTag(tag string) bool
}
