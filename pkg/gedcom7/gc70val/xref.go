package gc70val

import (
	"fmt"
	"regexp"
)

var regXref = regexp.MustCompile(fmt.Sprintf("%s[%s%s%s]+%s", G7Atsign, G7Underscore, G7UCLetter, G7Digit, G7Atsign))

// IsXref validates that a string is a validly formatted XRef.
func IsXref(v string) bool {
	return regXref.MatchString(v)
}
