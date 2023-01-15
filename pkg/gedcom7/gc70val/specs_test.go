package gc70val

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		tag  string
		want bool
	}{
		{"INDI", true},
		{"FAM", true},
		{"record-SOUR", true},
		{"xxxyyy", false},
	}
	specs := New()
	if specs == nil {
		t.Fatal("New() returned nil")
	}

	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			if got := specs.IsValidTag(tt.tag); got != tt.want {
				t.Errorf("IsValidTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpecs_AddExtTag(t *testing.T) {
	tests := []struct {
		name    string
		in      []byte
		extTags TagDef
	}{
		{
			name: "EXTTAG",
			in: []byte(`%YAML 1.2
---
lang: en-US

type: structure

uri: https://gedcom.io/terms/v7/EXTTAG

standard tag: EXTTAG

specification:
  - Extended tag
  - Just some random words.

payload: http://www.w3.org/2001/XMLSchema#string

substructures: {}

superstructures:
  "https://gedcom.io/terms/v7/record-SOUR": "{0:1}"
...`),
		},
	}
	spec := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := spec.AddExtTag(tt.in); err != nil {
				t.Errorf("AddExtTag() unexpected error = %v", err)
			}
		})
	}
}
