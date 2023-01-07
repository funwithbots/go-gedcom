package gedcom7

import (
	"bufio"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"go-gedcom/pkg/gedcom"
)

var DocPath = "data/"

func TestLoadDocument(t *testing.T) {
	tests := []struct {
		name      string
		file      string
		want      document
		xref      string
		rootCount int
		nodeCount int
		BOM       bool
	}{
		{
			name:      "shaw",
			file:      "../../" + DocPath + "shaw.ged",
			want:      document{},
			xref:      "@I1@",
			rootCount: 4885,
			nodeCount: 80604,
			BOM:       true,
		},
		{
			name:      "tgc551lf",
			file:      DocPath + "torture test/TGC551LF.ged",
			want:      document{},
			xref:      "@SUBMITTER@",
			rootCount: 64,
			nodeCount: 2161,
			BOM:       false,
		},
		{
			name:      "tgc551",
			file:      DocPath + "torture test/TGC551.ged",
			want:      document{},
			xref:      "@SUBMITTER@",
			rootCount: 64,
			nodeCount: 2161,
			BOM:       false,
		},
		{
			name:      "tgc55c",
			file:      DocPath + "torture test/TGC55C.ged",
			want:      document{},
			xref:      "@SUBMITTER@",
			rootCount: 66,
			nodeCount: 2197,
			BOM:       false,
		},
		{
			name:      "tgc55clf",
			file:      DocPath + "torture test/TGC55CLF.ged",
			want:      document{},
			xref:      "@SUBMITTER@",
			rootCount: 66,
			nodeCount: 2197,
			BOM:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, closer, err := readFile(tt.file)
			if err != nil {
				t.Fatalf("Error reading file %s: %v", tt.file, err.Error())
			}
			defer closer()

			doc := NewDocument(s, WithMaxDeprecatedTags("5.5.1"))
			if reflect.DeepEqual(doc.header, gedcom.NewNode(doc.header)) {
				t.Errorf("NewDocument() = missing header.")
			}
			if reflect.DeepEqual(doc.trailer, node{}) {
				t.Errorf("NewDocument() = missing header.")
			}
			if doc.Len() != tt.rootCount {
				t.Errorf("NewDocument() = wrong number root nodes. Wanted %d; got %d", tt.rootCount, doc.Len())
			}
			if nod := doc.GetXRef(tt.xref); nod == nil {
				t.Errorf("NewDocument() = missing xref nodes. Wanted a node; got nil")
			} else {
				if nod.Xref != tt.xref {
					t.Errorf("NewDocument() = missing xref node. Wanted %s; got %s", tt.xref, nod.Xref)
				}
			}

			if len(doc.GetWarnings()) != 0 {
				t.Errorf("NewDocument() flagged warnings. Got %d, wanted 0.", len(doc.GetWarnings()))
				for _, w := range doc.GetWarnings() {
					t.Logf("%s\n", w)
				}
			}

			doc.XRefCache.Range(func(k, v interface{}) bool {
				str := k.(string)
				if str[0] != '@' && str[len(str)-1] != '@' {
					t.Errorf("invalid xref. Should start and end with '@' %s", str)
				}
				return true
			})
			if tt.BOM && doc.bom != bom {
				t.Errorf("NewDocument() = %v; want %v", doc.bom, bom)
			}

			f, err := os.CreateTemp("", tt.name)
			if err != nil {
				t.Fatalf("Error creating file for %s: %s", tt.name, err.Error())
			}
			out := f.Name()
			defer os.Remove(out)
			defer f.Close()

			w := bufio.NewWriter(f)
			if err = doc.exportGedcom7(w); err != nil {
				t.Fatalf("Error writing file %s", out)
			}
			w.Flush()

			errCount, total := fileDiff(tt.file, out, 10)
			if errCount != 0 {
				t.Errorf("NewDocument() = %d rebuilding errors; want %d errors", errCount, 0)
			}
			if total != tt.nodeCount {
				t.Errorf("NewDocument() = %d rebuilding lines; want %d lines", total, tt.nodeCount)
			}

		})
	}
}

// readFile reads a file and returns a buffer
func readFile(fn string) (*bufio.Scanner, func(), error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, nil, err
	}
	closer := func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}

	s := bufio.NewScanner(f)

	return s, closer, nil
}

// fileDiff compares two files line by line and returns the count of differences and total lines.
func fileDiff(file1, file2 string, max int) (int, int) {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	f1, close1, err := readFile(file1)
	if err != nil {
		log.Fatal(err)
	}
	defer close1()

	f2, close2, err := readFile(file2)
	if err != nil {
		log.Fatal(err)
	}
	defer close2()

	var total, errorCount, x1 int
	f1.Split(scanLines) // deal with \r line endings
	for f1.Scan() {
		if errorCount > max {
			break
		}
		total++
		t1 := strings.TrimSpace(f1.Text())
		if !f2.Scan() {
			x1++
			continue
		}
		t2 := strings.TrimSpace(f2.Text())
		if t1 != t2 {
			errorCount++
			log.Printf("1 '%s'\n2 '%s'\nline %d, err # %d\n----\n", t1, t2, total, errorCount)
		}
	}
	if x1 != 0 {
		log.Printf("%d extra lines in %s\n", x1, file1)
	}

	var x2 int
	for f2.Scan() {
		x2++
	}
	if x2 != 0 {
		log.Printf("%d extra lines in %s\n", x2, file2)
	}
	return errorCount, total
}
