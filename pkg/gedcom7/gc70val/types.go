package gc70val

// Data types are defined in chapter 2 of https://gedcom.io/specifications/FamilySearchGEDCOMv7.pdf.

type typeDef struct {
	Lang          string   `yaml:"lang"`
	Type          string   `yaml:"type"`
	URI           string   `yaml:"uri"`
	Specification []string `yaml:"specification"`
}

func loadType(in []byte) (typeDef, error) {
	tm := typeDef{}
	if err := deserializeYAML(in, &tm); err != nil {
		return tm, err
	}
	tm.Type = extractFullTag(tm.URI)

	return tm, nil
}
