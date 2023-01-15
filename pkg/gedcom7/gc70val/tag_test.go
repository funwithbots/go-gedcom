package gc70val

import (
	"log"
	"math/rand"
	"strings"
	"testing"
	"time"

	"go-gedcom/pkg/abnf"
)

func Test_loadTags(t *testing.T) {
	in := []byte(`%YAML 1.2
---
lang: en-US

type: structure

uri: https://gedcom.io/terms/v7/ABBR

standard tag: ABBR

specification:
  - Abbreviation
  - A short name of a title, description, or name used for sorting, filing, and
    retrieving records.

payload: http://www.w3.org/2001/XMLSchema#string

substructures: {}

superstructures:
  "https://gedcom.io/terms/v7/record-SOUR": "{0:1}"
...`)

	// scanner := bufio.NewScanner(bytes.NewBufferString(in))

	data, err := loadTag(in)
	if err != nil {
		t.Fatal("couldn't load yaml data", err.Error())
	}
	if data.Tag != "ABBR" {
		t.Fatalf("missing Tag. Wanted ABBR; got %s", data.Tag)
	}
	if len(data.Specification) != 2 {
		t.Fatalf("wrong number of specifications. wanted 2; got %d", len(data.Specification))
	}
}

func Test_tagMeta_InferRule(t *testing.T) {
	for k := range validTags {
		if tag, ok := baseline.tags[string(k)]; ok {
			tag.InferRule()
			if tag.Rule == nil {
				t.Errorf("Unknown validation rule for %s:%s:%s.\n", tag.Type, tag.Payload, tag.URI)
				debugInferRule(t, baseline.tags[string(k)])
			}
		}
	}
}

func debugInferRule(t *testing.T, tm TagDef) {
	if tm.URI != "" {
		switch {
		case tm.Payload == "null", tm.Payload == "":
			t.Logf("%v: Null payload.\n", tm.Tag)
			tm.Rule = abnf.Validation["Null"]
		case len(tm.Payload) > 2 && tm.Payload[0] == '@' && tm.Payload[len(tm.Payload)-1] == '@':
			t.Logf("%v: Tag via @*@ pattern.\n", tm.Tag)
			tm.Rule = abnf.Validation["Tag"]
		case strings.HasPrefix(tm.Payload, URIPrefixes["g7"]):
			k := tm.Payload[len(URIPrefixes["g7"]):]
			if k[:5] == "type-" {
				k = k[5:]
			}
			tm.Rule = abnf.Validation[k]
			t.Logf("%v: g7 prefix. %s '%s'\n", tm.Tag, tm.Payload, k)
		case strings.HasPrefix(tm.Payload, URIPrefixes["xsd"]),
			strings.HasPrefix(tm.Payload, URIPrefixes["w3"]):
			parts := strings.Split(tm.Payload, "#")
			uc := toUpperCamel(parts[len(parts)-1])
			tm.Rule = abnf.Validation[uc]
			t.Logf("%v: xsd or w3 prefix. %s '%s'\n", tm.Tag, tm.Payload, uc)
		default:
			t.Logf("%v: No matching case.\n", tm.Tag)

		}
	} else {
		log.Printf("%v: No URI.\n", tm.Tag)
	}
	log.Println("=======================================")
}

func Test_toUpperCamel(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "basic",
			in:   "basic",
			want: "Basic",
		},
		{
			name: "hyphen",
			in:   "basic-words",
			want: "BasicWords",
		},
		{
			name: "underscore",
			in:   "basic_words",
			want: "BasicWords",
		},
		{
			name: "harder",
			in:   "basic-words_1_one-two-three_four",
			want: "BasicWords1OneTwoThreeFour",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toUpperCamel(tt.in); got != tt.want {
				t.Errorf("toUpperCamel() = %s, want %s", got, tt.want)
			}
		})
	}
}

func Test_tagDef_ValidatePayload(t *testing.T) {
	defs := New()
	rand.Seed(time.Now().UnixNano())

	// TODO Add more strings to exercise the validation rules.
	for _, tt := range defs.Tags {
		t.Run(tt.FullTag, func(t *testing.T) {
			var str []string
			payload := strings.Split(tt.Payload, "/")[len(strings.Split(tt.Payload, "/"))-1]
			switch {
			case payload == "null", payload == "":
				str = []string{""}
			case payload == "Y|<NULL>":
				str = []string{"Y"}
			case payload == "Y|N":
				str = []string{"Y", "N"}
			case payload == "type-Date#period":
				str = []string{"TO 1992"}
			case payload == "type-Date#exact":
				str = []string{"12 DEC 1992"}
			case payload == "type-Date":
				str = []string{"TO 1992"}
			case strings.HasSuffix(payload, ">@"):
				str = []string{"@I1@"}
			case strings.HasSuffix(payload, "Enum"):
				l := len(tt.EnumSet.Values)
				if l != 0 {
					str = tt.EnumSet.Values
				} else {
					if set, ok := baseline.enumSets[tt.URI]; ok {
						str = set.Values
					} else {
						str = []string{"unknown"}
					}
				}
			case payload == "XMLSchema#Language":
				str = []string{"en-US"}
			case payload == "XMLSchema#nonNegativeInteger":
				str = []string{"14"}
			case payload == "type-Age":
				str = []string{"14y"}
			case payload == "dcat#mediaType":
				str = []string{"text/plain"}
			case payload == "type-Time":
				str = []string{"12:00"}
			default:
				str = []string{"test value"}
			}
			for _, s := range str {
				if got := tt.ValidatePayload(s); !got {
					t.Errorf("FAIL Validate(%s), tag: %s, payload: %s", s, tt.FullTag, tt.Payload)
				}
			}
		})
	}
}
