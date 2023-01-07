// This file is generated - do not edit.
// Except where needed. Generator does NOT support Proposed RFC 7405 for handling case-insensitive strings.
// To process, remove %s from abnf definition and then replace operators.String() with operators.String()
// as needed.

package abnf

import "github.com/funwithbots/go-abnf/operators"

// Age validates Age = [ageBound D] ageDuration
func Age() operators.Operator {
	return operators.Concat(
		"Age",
		operators.Optional("[ageBound D]", operators.Concat(
			"ageBound D",
			AgeBound(),
			D(),
		)),
		AgeDuration(),
	)
}

// D validates D = %x20
func D() operators.Operator {
	return operators.Terminal("D", []byte{32})
}

// DateExact validates DateExact = day D month D year
func DateExact() operators.Operator {
	return operators.Concat(
		"DateExact",
		Day(),
		D(),
		Month(),
		D(),
		Year(),
	)
}

// DatePeriod validates DatePeriod = [ "TO" D date ] / "FROM" D date [ D "TO" D date ]
func DatePeriod() operators.Operator {
	return operators.Alts(
		"DatePeriod",
		operators.Optional("[ \"TO\" D date ]", operators.Concat(
			"\"TO\" D date",
			operators.String("TO", "TO"),
			D(),
			Date(),
		)),
		operators.Concat(
			"\"FROM\" D date [ D \"TO\" D date ]",
			operators.String("FROM", "FROM"),
			D(),
			Date(),
			operators.Optional("[ D \"TO\" D date ]", operators.Concat(
				"D \"TO\" D date",
				D(),
				operators.String("TO", "TO"),
				D(),
				Date(),
			)),
		),
	)
}

// DateValue validates DateValue = [ date / DatePeriod / dateRange / dateApprox ]
func DateValue() operators.Operator {
	return operators.Optional("DateValue", operators.Alts(
		"date / DatePeriod / dateRange / dateApprox",
		Date(),
		DatePeriod(),
		DateRange(),
		DateApprox(),
	))
}

// EOL validates EOL = %x0D [%x0A] / %x0A
func EOL() operators.Operator {
	return operators.Alts(
		"EOL",
		operators.Concat(
			"%x0D [%x0A]",
			operators.Terminal("%x0D", []byte{0x0D}),
			operators.Optional("[%x0A]", operators.Terminal("%x0A", []byte{0x0A})),
		),
		operators.Terminal("%x0A", []byte{0x0A}),
	)
}

// Enum validates Enum = stdEnum / extTag
func Enum() operators.Operator {
	return operators.Alts(
		"Enum",
		StdEnum(),
		ExtTag(),
	)
}

// Integer validates Integer = 1*digit
func Integer() operators.Operator {
	return operators.Repeat1Inf("Integer", Digit())
}

// Level validates Level = "0" / nonzero *digit
func Level() operators.Operator {
	return operators.Alts(
		"Level",
		operators.String("0", "0"),
		operators.Concat(
			"nonzero *digit",
			Nonzero(),
			operators.Repeat0Inf("*digit", Digit()),
		),
	)
}

// Line validates Line = Level D [Xref D] Tag [D LineVal] EOL
func Line() operators.Operator {
	return operators.Concat(
		"Line",
		Level(),
		D(),
		operators.Optional("[Xref D]", operators.Concat(
			"Xref D",
			Xref(),
			D(),
		)),
		Tag(),
		operators.Optional("[D LineVal]", operators.Concat(
			"D LineVal",
			D(),
			LineVal(),
		)),
		EOL(),
	)
}

// LineVal validates LineVal = pointer / lineStr
func LineVal() operators.Operator {
	return operators.Alts(
		"LineVal",
		Pointer(),
		LineStr(),
	)
}

// ListEnum validates List-Enum = Enum *(listDelim Enum)
func ListEnum() operators.Operator {
	return operators.Concat(
		"List-Enum",
		Enum(),
		operators.Repeat0Inf("*(listDelim Enum)", operators.Concat(
			"listDelim Enum",
			ListDelim(),
			Enum(),
		)),
	)
}

// ListText validates List-Text = list
func ListText() operators.Operator {
	return List()
}

// PersonalName validates PersonalName = nameStr / [nameStr] "/" [nameStr] "/" [nameStr]
func PersonalName() operators.Operator {
	return operators.Alts(
		"PersonalName",
		NameStr(),
		operators.Concat(
			"[nameStr] \"/\" [nameStr] \"/\" [nameStr]",
			operators.Optional("[nameStr]", NameStr()),
			operators.String("/", "/"),
			operators.Optional("[nameStr]", NameStr()),
			operators.String("/", "/"),
			operators.Optional("[nameStr]", NameStr()),
		),
	)
}

// Special validates Special = Text
func Special() operators.Operator {
	return Text()
}

// Tag validates Tag = stdTag / extTag
func Tag() operators.Operator {
	return operators.Alts(
		"Tag",
		StdTag(),
		ExtTag(),
	)
}

// Text validates Text = *anychar
func Text() operators.Operator {
	return operators.Repeat0Inf("Text", Anychar())
}

// Time validates Time = hour ":" minute [":" second ["." fraction]] ["Z"]
func Time() operators.Operator {
	return operators.Concat(
		"Time",
		Hour(),
		operators.String(":", ":"),
		Minute(),
		operators.Optional("[\":\" second [\".\" fraction]]", operators.Concat(
			"\":\" second [\".\" fraction]",
			operators.String(":", ":"),
			Second(),
			operators.Optional("[\".\" fraction]", operators.Concat(
				"\".\" fraction",
				operators.String(":", ":"),
				Fraction(),
			)),
		)),
		operators.Optional("[\"Z\"]", operators.String("Z", "Z")),
	)
}

// Xref validates Xref = atsign 1*tagchar atsign
func Xref() operators.Operator {
	return operators.Concat(
		"Xref",
		Atsign(),
		operators.Repeat1Inf("1*tagchar", Tagchar()),
		Atsign(),
	)
}

// AgeBound validates ageBound = "<" / ">"
func AgeBound() operators.Operator {
	return operators.Alts(
		"ageBound",
		operators.String("<", "<"),
		operators.String(">", ">"),
	)
}

// AgeDuration validates ageDuration = years [D months] [D weeks] [D days] / months [D weeks] [D days] / weeks [D days] / days
func AgeDuration() operators.Operator {
	return operators.Alts(
		"ageDuration",
		operators.Concat(
			"years [D months] [D weeks] [D days]",
			Years(),
			operators.Optional("[D months]", operators.Concat(
				"D months",
				D(),
				Months(),
			)),
			operators.Optional("[D weeks]", operators.Concat(
				"D weeks",
				D(),
				Weeks(),
			)),
			operators.Optional("[D days]", operators.Concat(
				"D days",
				D(),
				Days(),
			)),
		),
		operators.Concat(
			"months [D weeks] [D days]",
			Months(),
			operators.Optional("[D weeks]", operators.Concat(
				"D weeks",
				D(),
				Weeks(),
			)),
			operators.Optional("[D days]", operators.Concat(
				"D days",
				D(),
				Days(),
			)),
		),
		operators.Concat(
			"weeks [D days]",
			Weeks(),
			operators.Optional("[D days]", operators.Concat(
				"D days",
				D(),
				Days(),
			)),
		),
		Days(),
	)
}

// Anychar validates anychar = %x09-10FFFF
func Anychar() operators.Operator {
	return operators.Range("anychar", []byte{9}, []byte{0x10, 0xFF, 0xFF})
}

// Atsign validates atsign = %x40
func Atsign() operators.Operator {
	return operators.Terminal("atsign", []byte{64})
}

// Banned validates banned = %x00-08 / %x0B-0C / %x0E-1F ; C0 other than LF CR and Tab / %x7F ; DEL / %x80-9F ; C1 / %xD800-DFFF ; Surrogates / %xFFFE-FFFF
func Banned() operators.Operator {
	return operators.Alts(
		"banned",
		operators.Range("%x00-08", []byte{0}, []byte{8}),
		operators.Range("%x0B-0C", []byte{0x0B}, []byte{0x0C}),
		operators.Range("%x0E-1F", []byte{0x0E}, []byte{0x1F}),
		operators.Terminal("%x7F", []byte{0x7F}),
		operators.Range("%x80-9F", []byte{0x80}, []byte{0x9F}),
		operators.Range("%xD800-DFFF", []byte{0xD8, 0}, []byte{0xDF, 0xFF}),
		operators.Range("%xFFFE-FFFF", []byte{0xFF, 0xFE}, []byte{0xFF, 0xFF}),
	)
}

// Calendar validates calendar = "GREGORIAN" / "JULIAN" / "FRENCH_R" / "HEBREW" / extTag
func Calendar() operators.Operator {
	return operators.Alts(
		"calendar",
		operators.String("GREGORIAN", "GREGORIAN"),
		operators.String("JULIAN", "JULIAN"),
		operators.String("FRENCH_R", "FRENCH_R"),
		operators.String("HEBREW", "HEBREW"),
		ExtTag(),
	)
}

// Date validates date = [calendar D] [[day D] month D] year [D epoch]
func Date() operators.Operator {
	return operators.Concat(
		"date",
		operators.Optional("[calendar D]", operators.Concat(
			"calendar D",
			Calendar(),
			D(),
		)),
		operators.Optional("[[day D] month D]", operators.Concat(
			"[day D] month D",
			operators.Optional("[day D]", operators.Concat(
				"day D",
				Day(),
				D(),
			)),
			Month(),
			D(),
		)),
		Year(),
		operators.Optional("[D epoch]", operators.Concat(
			"D epoch",
			D(),
			Epoch(),
		)),
	)
}

// DateApprox validates dateApprox = ("ABT" / "CAL" / "EST") D date
func DateApprox() operators.Operator {
	return operators.Concat(
		"dateApprox",
		operators.Alts(
			"\"ABT\" / \"CAL\" / \"EST\"",
			operators.String("ABT", "ABT"),
			operators.String("CAL", "CAL"),
			operators.String("EST", "EST"),
		),
		D(),
		Date(),
	)
}

// DateRange validates dateRange = "BET" D date D "AND" D date / "AFT" D date / "BEF" D date
func DateRange() operators.Operator {
	return operators.Alts(
		"dateRange",
		operators.Concat(
			"\"BET\" D date D \"AND\" D date",
			operators.String("BET", "BET"),
			D(),
			Date(),
			D(),
			operators.String("AND", "AND"),
			D(),
			Date(),
		),
		operators.Concat(
			"\"AFT\" D date",
			operators.String("AFT", "AFT"),
			D(),
			Date(),
		),
		operators.Concat(
			"\"BEF\" D date",
			operators.String("BEF", "BEF"),
			D(),
			Date(),
		),
	)
}

// DateRestrict validates dateRestrict = "FROM" / "TO" / "BET" / "AND" / "BEF" / "AFT" / "ABT" / "CAL" / "EST" / "BCE"
func DateRestrict() operators.Operator {
	return operators.Alts(
		"dateRestrict",
		operators.String("FROM", "FROM"),
		operators.String("TO", "TO"),
		operators.String("BET", "BET"),
		operators.String("AND", "AND"),
		operators.String("BEF", "BEF"),
		operators.String("AFT", "AFT"),
		operators.String("ABT", "ABT"),
		operators.String("CAL", "CAL"),
		operators.String("EST", "EST"),
		operators.String("BCE", "BCE"),
	)
}

// Day validates day = Integer
func Day() operators.Operator {
	return Integer()
}

// Days validates days = Integer %x64
func Days() operators.Operator {
	return operators.Concat(
		"days",
		Integer(),
		operators.Terminal("%x64", []byte{0x64}), // d
	)
}

// Digit validates digit = %x30-39
func Digit() operators.Operator {
	return operators.Range("digit", []byte{48}, []byte{57})
}

// Epoch validates epoch = "BCE" / extTag
func Epoch() operators.Operator {
	return operators.Alts(
		"epoch",
		operators.String("BCE", "BCE"),
		ExtTag(),
	)
}

// ExtTag validates extTag = underscore 1*tagchar
func ExtTag() operators.Operator {
	return operators.Concat(
		"extTag",
		Underscore(),
		operators.Repeat1Inf("1*tagchar", Tagchar()),
	)
}

// Fraction validates fraction = 1*digit
func Fraction() operators.Operator {
	return operators.Repeat1Inf("fraction", Digit())
}

// Hour validates hour = digit / ("0" / "1") digit / "2" ("0" / "1" / "2" / "3")
func Hour() operators.Operator {
	return operators.Alts(
		"hour",
		Digit(),
		operators.Concat(
			"(\"0\" / \"1\") digit",
			operators.Alts(
				"\"0\" / \"1\"",
				operators.String("0", "0"),
				operators.String("1", "1"),
			),
			Digit(),
		),
		operators.Concat(
			"\"2\" (\"0\" / \"1\" / \"2\" / \"3\")",
			operators.String("2", "2"),
			operators.Alts(
				"\"0\" / \"1\" / \"2\" / \"3\"",
				operators.String("0", "0"),
				operators.String("1", "1"),
				operators.String("2", "2"),
				operators.String("3", "3"),
			),
		),
	)
}

// LineStr validates lineStr = (nonAt / atsign atsign) *nonEOL
func LineStr() operators.Operator {
	return operators.Concat(
		"lineStr",
		operators.Alts(
			"nonAt / atsign atsign",
			NonAt(),
			operators.Concat(
				"atsign atsign",
				Atsign(),
				Atsign(),
			),
		),
		operators.Repeat0Inf("*nonEOL", NonEOL()),
	)
}

// List validates list = listItem *(listDelim listItem)
func List() operators.Operator {
	return operators.Concat(
		"list",
		ListItem(),
		operators.Repeat0Inf("*(listDelim listItem)", operators.Concat(
			"listDelim listItem",
			ListDelim(),
			ListItem(),
		)),
	)
}

// ListDelim validates listDelim = *D "," *D
func ListDelim() operators.Operator {
	return operators.Concat(
		"listDelim",
		operators.Repeat0Inf("*D", D()),
		operators.String(",", ","),
		operators.Repeat0Inf("*D", D()),
	)
}

// ListItem validates listItem = [ nocommasp / nocommasp *nocomma nocommasp ]
func ListItem() operators.Operator {
	return operators.Optional("listItem", operators.Alts(
		"nocommasp / nocommasp *nocomma nocommasp",
		Nocommasp(),
		operators.Concat(
			"nocommasp *nocomma nocommasp",
			Nocommasp(),
			operators.Repeat0Inf("*nocomma", Nocomma()),
			Nocommasp(),
		),
	))
}

// Minute validates minute = ("0" / "1" / "2" / "3" / "4" / "5") digit
func Minute() operators.Operator {
	return operators.Concat(
		"minute",
		operators.Alts(
			"\"0\" / \"1\" / \"2\" / \"3\" / \"4\" / \"5\"",
			operators.String("0", "0"),
			operators.String("1", "1"),
			operators.String("2", "2"),
			operators.String("3", "3"),
			operators.String("4", "4"),
			operators.String("5", "5"),
		),
		Digit(),
	)
}

// Month validates month = stdTag / extTag
func Month() operators.Operator {
	return operators.Alts(
		"month",
		StdTag(),
		ExtTag(),
	)
}

// Months validates months = Integer %x6D
func Months() operators.Operator {
	return operators.Concat(
		"months",
		Integer(),
		operators.Terminal("%x6D", []byte{0x6D}), // m
	)
}

// MediaType defines the structure of a media type as `type + "/" + subtype`
func MediaType() operators.Operator {
	return operators.Concat("type / subtype",
		MtType(),
		operators.String("/", "/"),
		MtSubtype(),
		operators.Repeat0Inf("parameters", MtParameter()),
	)
}

// MtAttribute validates mt-attribute = mt-token
func MtAttribute() operators.Operator {
	return MtName()
}

// MtName validates mt-name = mt-name-first *126(mt-name-chars)
func MtName() operators.Operator {
	return operators.Concat("mt-name-first *126(mt-name-chars)",
		MtFirstChar(),
		operators.Repeat("*126(mt-name-chars)", 0, 126, MtChars()),
	)
}

// MtFirstChar validates mt-name-first = mt-alpha / mt-digit
func MtFirstChar() operators.Operator {
	return operators.Alts(
		"mt-name-first",
		operators.Range("%x41-5A", []byte{0x41}, []byte{0x5A}),
		Digit(),
	)
}

// MtChar validates mt-name-chars = mt-alpha / mt-digit / "!" / "#" / "$" / "&" / "-" / "^" / "_" / "." / "+"
func MtChars() operators.Operator {
	return operators.Alts(
		"mt-char",
		operators.Range("%x41-5a", []byte{0x41}, []byte{0x5A}),
		Digit(),
		operators.String("!", "!"),
		operators.String("#", "#"),
		operators.String("$", "$"),
		operators.String("&", "&"),
		operators.String("-", "-"),
		operators.String("^", "^"),
		operators.String("_", "_"),
		operators.String(".", "."),
		operators.String("+", "+"),
	)
}

// MtParameter validates mt-parameter = mt-attribute "=" mt-value
func MtParameter() operators.Operator {
	return operators.Concat(
		"mt-parameter",
		operators.Optional("OWS", OWS()),
		operators.String(";", ";"),
		operators.Optional("OWS", OWS()),
		MtAttribute(),
		operators.String("=", "="),
		MtName(),
	)
}

// MtType validates mt-type = mt-token
func MtType() operators.Operator {
	return MtName()
}

// MtSubtype validates mt-subtype = mt-token
func MtSubtype() operators.Operator {
	return MtName()
}

// NameChar validates nameChar = %x20-2E / %x30-10FFFF
func NameChar() operators.Operator {
	return operators.Alts(
		"nameChar",
		operators.Range("%x20-2E", []byte{32}, []byte{46}),
		operators.Range("%x30-10FFFF", []byte{48}, []byte{16, 255, 255}),
	)
}

// NameStr validates nameStr = 1*nameChar
func NameStr() operators.Operator {
	return operators.Repeat1Inf("nameStr", NameChar())
}

func Name() operators.Operator {
	return operators.Concat(
		"name",
		operators.Repeat1Inf("nameStr", NameStr()),
		operators.Optional("slash", operators.String("/", "/")),
		NameStr(),
		operators.Optional("slash", operators.String("/", "/")),
		NameStr(),
	)
}

// Nocomma validates nocomma = %x09-2B / %x2D-10FFFF
func Nocomma() operators.Operator {
	return operators.Alts(
		"nocomma",
		operators.Range("%x09-2B", []byte{9}, []byte{43}),
		operators.Range("%x2D-10FFFF", []byte{45}, []byte{16, 255, 255}),
	)
}

// Nocommasp validates nocommasp = %x09-1D / %x21-2B / %x2D-10FFFF
func Nocommasp() operators.Operator {
	return operators.Alts(
		"nocommasp",
		operators.Range("%x09-1D", []byte{9}, []byte{29}),
		operators.Range("%x21-2B", []byte{33}, []byte{43}),
		operators.Range("%x2D-10FFFF", []byte{45}, []byte{16, 255, 255}),
	)
}

// NonAt validates nonAt = %x09 / %x20-3F / %x41-10FFFF
func NonAt() operators.Operator {
	return operators.Alts(
		"nonAt",
		operators.Terminal("%x09", []byte{9}),
		operators.Range("%x20-3F", []byte{32}, []byte{63}),
		operators.Range("%x41-10FFFF", []byte{65}, []byte{16, 255, 255}),
	)
}

// NonEOL validates nonEOL = %x09 / %x20-10FFFF
func NonEOL() operators.Operator {
	return operators.Alts(
		"nonEOL",
		operators.Terminal("%x09", []byte{9}),
		operators.Range("%x20-10FFFF", []byte{32}, []byte{16, 255, 255}),
	)
}

// Nonzero validates nonzero = %x31-39
func Nonzero() operators.Operator {
	return operators.Range("nonzero", []byte{49}, []byte{57})
}

// Pointer validates pointer = voidPtr / Xref
func Pointer() operators.Operator {
	return operators.Alts(
		"pointer",
		VoidPtr(),
		Xref(),
	)
}

// Second validates second = ("0" / "1" / "2" / "3" / "4" / "5") digit
func Second() operators.Operator {
	return operators.Concat(
		"second",
		operators.Alts(
			"\"0\" / \"1\" / \"2\" / \"3\" / \"4\" / \"5\"",
			operators.String("0", "0"),
			operators.String("1", "1"),
			operators.String("2", "2"),
			operators.String("3", "3"),
			operators.String("4", "4"),
			operators.String("5", "5"),
		),
		Digit(),
	)
}

// StdEnum validates stdEnum = stdTag / Integer
func StdEnum() operators.Operator {
	return operators.Alts(
		"stdEnum",
		StdTag(),
		Integer(),
	)
}

// StdTag validates stdTag = ucletter *tagchar
func StdTag() operators.Operator {
	return operators.Concat(
		"stdTag",
		Ucletter(),
		operators.Repeat0Inf("*tagchar", Tagchar()),
	)
}

// Tagchar validates tagchar = ucletter / digit / underscore
func Tagchar() operators.Operator {
	return operators.Alts(
		"tagchar",
		Ucletter(),
		Digit(),
		Underscore(),
	)
}

// Ucletter validates ucletter = %x41-5A
func Ucletter() operators.Operator {
	return operators.Range("ucletter", []byte{65}, []byte{90})
}

// Underscore validates underscore = %x5F
func Underscore() operators.Operator {
	return operators.Terminal("underscore", []byte{95})
}

// VoidPtr validates voidPtr =
func VoidPtr() operators.Operator {
	return operators.String("voidPtr", "@VOID@")
}

// Weeks validates weeks = Integer %x77
func Weeks() operators.Operator {
	return operators.Concat(
		"weeks",
		Integer(),
		operators.Terminal("%x77", []byte{0x77}), // w
	)
}

// Year validates year = Integer
func Year() operators.Operator {
	return Integer()
}

// Years validates years = Integer %x79
func Years() operators.Operator {
	return operators.Concat(
		"years",
		Integer(),
		operators.Terminal("%x79", []byte{0x79}), // y
	)
}
