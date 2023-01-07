package abnf

import "github.com/funwithbots/go-abnf/operators"

// Null
func Null() operators.Operator {
	return operators.Terminal("Null", []byte{0})
}

// YOrNull validates [Y|<NULL>]
func YOrNull() operators.Operator {
	return operators.Optional("Y|Null",
		operators.String("Y", "Y"),
	)
}

// Lang is defined as BCP 47 or a short list of alternatives from
// https://gedcom.io/specifications/FamilySearchGEDCOMv7.pdf pp 77-78
func Lang() operators.Operator {
	return operators.Alts("LANG",
		operators.String("und", "und"),
		operators.String("mul", "mul"),
		operators.String("zxx", "zxx"),
		LangTag(),
	)
}

func OWS() operators.Operator {
	return operators.Alts("OWS",
		operators.Terminal("HTAB", []byte{9}),
		operators.Terminal("SP", []byte{20}),
	)
}
