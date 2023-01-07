package gc70val

import (
	"strings"
)

type enumSet struct {
	Lang   string   `yaml:"lang"`
	Type   string   `yaml:"type"`
	URI    string   `yaml:"uri"`
	Values []string `yaml:"enumeration values"`

	StandardTag string
	FullTag     string
}

func loadEnumSet(in []byte) (enumSet, error) {
	es := enumSet{}
	if err := deserializeYAML(in, &es); err != nil {
		return es, err
	}

	tag := extractFullTag(es.URI)
	es.FullTag = tag
	t := strings.Split(tag, "-")
	es.StandardTag = t[len(t)-1]
	for i, v := range es.Values {
		es.Values[i] = extractFullTag(v)
	}

	return es, nil
}
