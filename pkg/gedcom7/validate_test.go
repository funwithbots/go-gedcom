package gedcom7

import (
	"testing"

	"github.com/funwithbots/go-gedcom/pkg/gedcom"
)

func Test_document_ValidateNode(t *testing.T) {
	name := "data/test/valid.ged"
	nodeCh := make(chan *gedcom.Node)

	d, err := NewDocumentFromFile(name)
	doc := d.(*document)
	if err != nil {
		t.Fatalf("couldn't load document %s", err.Error())
	}

	go func() {
		for _, node := range doc.records {
			if node == nil {
				continue
			}
			loadNodes(node, nodeCh)
		}
		close(nodeCh)
	}()

	i := 0
	for node := range nodeCh {
		i++
		err := doc.ValidateNode(*node)
		if err != nil {
			t.Errorf("Validate() error = %v", err)
		}
	}
	t.Logf("Tested %d nodes", i)
}

func loadNodes(n *gedcom.Node, ch chan *gedcom.Node) {
	ch <- n
	for _, v := range n.GetSubnodes() {
		loadNodes(v, ch)
	}
}
