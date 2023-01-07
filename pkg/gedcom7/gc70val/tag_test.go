package gc70val

import (
	"log"
	"strings"
	"testing"

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

func TestListTags(t *testing.T) {
	var tags string
	for k := range validTags {
		tags += string(k) + " "
	}
	t.Logf("standard %s", tags)

	for k := range customTags {
		tags += string(k) + " "
	}
	t.Logf("custom %s", tags)
}

func Test_tagMeta_InferRule(t *testing.T) {
	for k := range validTags {
		if tag, ok := baseline.tags[string(k)]; ok {
			tag.InferRule()
			if tag.Rule == nil {
				t.Errorf("Unknown validation rule for %s:%s:%s.\n", tag.Type, tag.Payload, tag.URI)
				debugInferRule(baseline.tags[string(k)])
			}
		}
	}
}

func debugInferRule(tm tagDef) {
	if tm.URI != "" {
		switch {
		case tm.Payload == "null", tm.Payload == "":
			log.Printf("%v: Null payload.\n", tm.Tag)
			tm.Rule = abnf.Validation["Null"]
		case len(tm.Payload) > 2 && tm.Payload[0] == '@' && tm.Payload[len(tm.Payload)-1] == '@':
			log.Printf("%v: Tag via @*@ pattern.\n", tm.Tag)
			tm.Rule = abnf.Validation["Tag"]
		case strings.HasPrefix(tm.Payload, URIPrefixes["g7"]):
			k := tm.Payload[len(URIPrefixes["g7"]):]
			if k[:5] == "type-" {
				k = k[5:]
			}
			tm.Rule = abnf.Validation[k]
			log.Printf("%v: g7 prefix. %s '%s'\n", tm.Tag, tm.Payload, k)
		case strings.HasPrefix(tm.Payload, URIPrefixes["xsd"]),
			strings.HasPrefix(tm.Payload, URIPrefixes["w3"]):
			parts := strings.Split(tm.Payload, "#")
			uc := toUpperCamel(parts[len(parts)-1])
			tm.Rule = abnf.Validation[uc]
			log.Printf("%v: xsd or w3 prefix. %s '%s'\n", tm.Tag, tm.Payload, uc)
		default:
			log.Printf("%v: No matching case.\n", tm.Tag)

		}
	} else {
		log.Printf("%v: No URI.\n", tm.Tag)
	}
	log.Println("=======================================")
}

// func Test_tagMeta_Validate(t *testing.T) {
// 	fn := DocPath + "shaw.ged"
//
// 	s, closer, err := readFile(fn)
// 	if err != nil {
// 		t.Fatalf("Error reading file %s: %v", fn, err.Error())
// 	}
// 	defer closer()
//
// 	nodes, err := Load(s)
// 	if err != nil {
// 		t.Fatalf("Error loading nodes from file %s: %v", fn, err.Error())
// 	}
//
// 	var errorCount int
// 	for i, tt := range nodes {
// 		tm, ok := tagList[tt.Tag]
// 		if !ok {
// 			t.Errorf("Unknown tag %s at index %d\nLine: %s\n", tt.Tag, i, tt.String())
// 			continue
// 		}
// 		payloadVal := tm.Rule([]byte(tt.Payload))
// 		if payloadVal == nil {
// 			continue
// 		}
// 		if len(payloadVal[0].Value) != len(tt.Payload) {
// 			errorCount++
// 			t.Errorf("Invalid payload for tag %s at index %d\n:Payload: '%s' : '%s'\nLine: %s\n", tt.Tag, i, tt.Payload, payloadVal[0].Value, tt.String())
// 		}
// 		if errorCount > 10 {
// 			break
// 		}
// 	}
// }

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
