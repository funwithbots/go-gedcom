package gedcom7

import (
	"fmt"

	"github.com/funwithbots/go-abnf/operators"

	"go-gedcom/pkg/gedcom"
	"go-gedcom/pkg/gedcom7/gc70val"
)

// Validate checks the entire GEDCOM document for technical validity
func (d *document) Validate() error {
	// TODO Verify HEAD and TRLR are present

	// TODO Check for warnings

	// TODO Traverse records and validate nodes

	return fmt.Errorf("not implemented")
}

// ValidateNode checks the node and ensures it passes validation rules
// TODO Still need to verify the node is validly placed in the document tree.
func (d *document) ValidateNode(node gedcom.Node) error {
	line := node.GetValue().(*Line)
	if v := line.Validate(); len(v) > 0 {
		return fmt.Errorf("line validator logged %d errors", len(v))
	}
	enums := d.Validator.GetEnums(line.Tag)
	if enums != nil {
		for _, enum := range enums {
			if enum == line.Payload {
				return nil
			}
		}
		return fmt.Errorf("invalid enum value: %s", line.Payload)
	}

	rule := d.findRule(node)
	if rule == nil {
		return fmt.Errorf("no rule found for tag %s", line.Tag)
	}

	if x := rule([]byte(line.Payload)); x == nil {
		return fmt.Errorf("invalid payload for tag %s", line.Tag)
	}

	return nil
}

func (d *document) findRule(node gedcom.Node) operators.Operator {
	line := node.GetValue().(*Line)
	tag := line.Tag
	if line.Level == 0 && line.Tag != gc70val.TagHEAD && line.Tag != gc70val.TagTRLR {
		tag = "record-" + line.Tag
	}
	if rule := d.Validator.GetRule(tag); rule != nil {
		return rule
	}

	// Try parent-tag as key
	if parent := node.GetParent(); parent != nil {
		tag = fmt.Sprintf("%s-%s", parent.GetValue().(*Line).Tag, line.Tag)
	}

	return d.Validator.GetRule(tag)
}
