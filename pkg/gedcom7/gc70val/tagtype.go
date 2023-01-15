package gc70val

const (
	tagTypePseudostructure = "pseudostructure"
	tagTypeStructure       = "structure"
	tagTypeExtended        = "extended"
	tagTypeDeprecated      = "deprecated"
)

var (
	tagTypes = map[string]interface{}{
		tagTypePseudostructure: nil,
		tagTypeStructure:       nil,
		tagTypeExtended:        nil,
		tagTypeDeprecated:      nil,
	}
)

func IsValidTagType(tp string) bool {
	_, ok := tagTypes[tp]
	return ok
}
