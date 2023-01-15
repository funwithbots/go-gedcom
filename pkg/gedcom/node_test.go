package gedcom_test

import (
	"strconv"
	"sync"
	"testing"

	"go-gedcom/pkg/gedcom"
	"go-gedcom/pkg/gedcom7"
)

const (
	testFile = "data/maximal70-lds-mod.ged"
)

func TestNode_ProcessSubnodes(t *testing.T) {
	// TODO Need to refactor when ready to support additional GEDCOM versions.

	// Return the text of the node.
	fn := func(n *gedcom.Node) interface{} {
		if nn, ok := n.GetValue().(*gedcom7.Line); ok {
			return nn.String()
		}
		return nil
	}

	doc, err := gedcom7.NewDocumentFromFile(testFile)
	if err != nil {
		t.Fatalf("Error loading document %s: %v", testFile, err.Error())
	}

	for _, tt := range doc.Records() {
		if _, ok := tt.GetValue().(*gedcom7.Line); !ok {
			t.Errorf("Line value is not a gedcom7.Line")
			continue
		}
		line := tt.GetValue().(*gedcom7.Line).Text
		wg := new(sync.WaitGroup)
		wg.Add(1)
		if got := tt.ProcessTree(fn, wg); got == nil {
			t.Errorf("%s Line.ProcessSubnodes() returned nil", line)
		} else {
			v := got[0].(string)
			if v != line {
				t.Errorf("ProcessSubnodes() = %v, want %v", got, line)
			}
			for _, g := range got {
				gg := g.(string)
				if len(gg) > 7 && gg[:6] == "1 NOTE" {
					count, err := strconv.Atoi(gg[7:])
					if err != nil {
						t.Errorf("ProcessSubnodes() %s payload isn't an int.", gg)
						break
					}
					if count != len(got) {
						t.Errorf("ProcessSubnodes() line count wrong. Want %d, got %d.", count, len(got))
						break
					}
					break
				}
			}
		}
	}
}
