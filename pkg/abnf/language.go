package abnf

// These operators are derived from BCP 47 / RFC-5646 as defined in
// https://www.rfc-editor.org/rfc/pdfrfc/rfc5646.txt.pdf
// and by extension ISO 15924, ISO 3166-1, UN M.49, ISO 639

import "github.com/funwithbots/go-abnf/operators"

// A-Za-z
func Alpha() operators.Operator {
	return operators.Alts("ascii letters",
		operators.Range("lc", []byte{65}, []byte{90}),
		operators.Range("lc", []byte{97}, []byte{122}),
	)
}

func Alphanum() operators.Operator {
	return operators.Alts("alphanum",
		Digit(),
		Alpha(),
	)
}

func LangPrivateUse() operators.Operator {
	return operators.Concat("privateuse",
		operators.String("x", "x"),
		operators.Repeat1Inf("",
			operators.Concat("",
				operators.String("-", "-"),
				operators.Repeat("1*8alphanum", 1, 8, Alphanum()),
			),
		),
	)
}

func LangSingleton() operators.Operator {
	return operators.Alts("singleton",
		Digit(),
		operators.Range("%x41-57", []byte{65}, []byte{87}),
		operators.Range("%x59-5A", []byte{89}, []byte{90}),
		operators.Range("%x61-77", []byte{97}, []byte{119}),
		operators.Range("%x79-7A", []byte{121}, []byte{122}),
	)
}

func LangExtension() operators.Operator {
	return operators.Concat("extension",
		LangSingleton(),
		operators.Repeat1Inf("1*(\"-\" (2*8alphanum))",
			operators.Concat("\"-\" (2*8alphanum)",
				operators.String("-", "-"),
				operators.Repeat("2*8alphnum", 2, 8, Alphanum()),
			),
		),
	)
}

func LangVariant() operators.Operator {
	return operators.Alts("variant",
		operators.Repeat("5*8alphanum", 5, 8, Alphanum()),
		operators.Concat("(DIGIT 3alphnum",
			Digit(),
			operators.RepeatN("3alphanum", 3, Alphanum()),
		),
	)
}

func LangRegion() operators.Operator {
	return operators.Alts("region",
		operators.RepeatN("2ALPHA", 2, Alpha()),
		operators.RepeatN("3DIGIT", 3, Digit()),
	)
}

func LangScript() operators.Operator {
	return operators.RepeatN("4ALPHA", 4, Alpha())
}

func ExtLang() operators.Operator {
	return operators.Concat("extlang",
		operators.RepeatN("3ALPHA", 3, Alpha()),
		operators.Repeat0Inf("*2(\"-\" 3ALPHA)",
			operators.Concat("(\"-\" 3ALPHA)",
				operators.String("-", "-"),
				operators.RepeatN("3ALPHA", 3, Alpha()),
			),
		),
	)
}

func Language() operators.Operator {
	return operators.Alts("language",
		operators.Concat("2*3ALPHA [\"-\" extlang]",
			operators.Repeat("2*3ALPHA", 2, 3, Alpha()),
			operators.Optional("[\"-\" extlang]",
				operators.Concat("\"-\" extlang",
					operators.String("-", "-"),
					ExtLang(),
				),
			),
		),
		operators.RepeatN("4ALPHA", 4, Alpha()),
		operators.Repeat("5*8ALPHA", 5, 8, Alpha()),
	)
}

func LangTag() operators.Operator {
	return operators.Concat("langtag",
		Language(),
		operators.Optional("[\"-\" script]",
			operators.Concat("\"-\" script",
				operators.String("-", "-"),
				LangScript(),
			),
		),
		operators.Optional("[\"-\" region]",
			operators.Concat("\"-\" region",
				operators.String("-", "-"),
				LangRegion(),
			),
		),
		operators.Repeat0Inf("*(\"-\" variant)",
			operators.Concat("(\"-\" variant)",
				operators.String("-", "-"),
				LangVariant(),
			),
		),
		operators.Repeat0Inf("*(\"-\" extension)",
			operators.Concat("(\"-\" extension)",
				operators.String("-", "-"),
				LangExtension(),
			),
		),
		operators.Optional("[\"-\" privateuse]",
			operators.Concat("\"-\" privateuse",
				operators.String("-", "-"),
				LangPrivateUse(),
			),
		),
	)
}

// Grandfathered supports non-redundant tags registered during the RFC 3066 era
func Grandfathered() operators.Operator {
	return operators.Alts("grandfathered",
		LangIrregular(),
		LangRegular(),
	)
}

// LangIrregular supports non-redundant tags registered during the RFC 3066 era
func LangIrregular() operators.Operator {
	return operators.Alts("irregular",
		operators.String("en-GB-oed", "en-GB-oed"),
		operators.String("i-ami", "i-ami"),
		operators.String("i-bnn", "i-bnn"),
		operators.String("i-default", "i-default"),
		operators.String("i-enochian", "i-enochian"),
		operators.String("i-hak", "i-hak"),
		operators.String("i-klingon", "i-klingon"),
		operators.String("i-lux", "i-lux"),
		operators.String("i-mingo", "i-mingo"),
		operators.String("i-navajo", "i-navajo"),
		operators.String("i-pwn", "i-pwn"),
		operators.String("i-tao", "i-tao"),
		operators.String("i-tay", "i-tay"),
		operators.String("i-tsu", "i-tsu"),
		operators.String("sgn-BE-FRE", "sgn-BE-FRE"),
		operators.String("sgn-BE-NL", "sgn-BE-NL"),
		operators.String("sgn-CH-DE", "sgn-CH-DE"),
	)
}

// LangRegular supports non-redundant tags registered during the RFC 3066 era
func LangRegular() operators.Operator {
	return operators.Alts("regular",
		operators.String("art-lojban", "art-lojban"),
		operators.String("cel-gaulish", "cel-gaulish"),
		operators.String("no-bok", "no-bok"),
		operators.String("no-nyn", "no-nyn"),
		operators.String("zh-guoyu", "zh-guoyu"),
		operators.String("zh-hakka", "zh-hakka"),
		operators.String("zh-min", "zh-min"),
		operators.String("zh-min-nan", "zh-min-nan"),
		operators.String("zh-xiang", "zh-xiang"),
	)
}
