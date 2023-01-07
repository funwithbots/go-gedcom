package gc70val

import (
	"errors"
	"regexp"
	"strconv"
)

// Data types are defined in chapter 2 of https://gedcom.io/specifications/FamilySearchGEDCOMv7.pdf.

// dataType defines a generic interface for gedcom data types available to Line.LineVal
type dataType interface {
	// validate determines if a line is properly formatted and meets content restrictions
	// Hard errors will fail and an error
	// Soft errors will pass but still return an error. Data would be available if the value is parsed.
	// Callers must check the error level regardless
	validate(string) (ok bool, err error)
}

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

type dtText string

func (tx dtText) validate(str string) (bool, error) {
	reg := regexp.MustCompile(`[\x{9}-\x{10ffff}]+`)

	ok := reg.MatchString(str)

	if !ok {
		return false, errors.New("invalid characters in text")
	}
	return true, nil
}

type dtInteger int

func (int dtInteger) validate(str string) (bool, error) {
	if len(str) == 0 {
		return false, errors.New("empty value")
	}

	i := 0
	var err error
	if i, err = strconv.Atoi(str); err != nil {
		return false, err
	}
	if i < 0 {
		return false, errors.New("non-negative integer:")
	}

	if str[0] == 0 {
		return true, errors.New("warn leading zeroes should be omitted")
	}

	return true, nil
}
