package gedcom7

import (
	"fmt"

	"github.com/funwithbots/go-gedcom/pkg/gedcom"
)

type Warning struct {
	Level   gedcom.Level
	Line    *Line
	Message string
}

func (w Warning) String() string {
	return fmt.Sprintf("%d: Line %s (%s)", w.Level, w.Line, w.Message)
}
