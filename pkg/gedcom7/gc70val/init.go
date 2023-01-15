package gc70val

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	abnfDir = "data/abnf"
	logFN   = "gedcom7.log"
)

var (
	logFile  *os.File
	baseline = struct {
		tags      map[string]TagDef
		calendars map[string]calDef
		types     map[string]typeDef
		enumSets  map[string]enumSet
	}{}
)

//go:embed data/abnf/*
var abnfFS embed.FS

func init() {
	var (
		tags      = pseudoTags
		types     = make(map[string]typeDef)
		calendars = make(map[string]calDef)
		enumSets  = make(map[string]enumSet)

		err error
	)
	logFile, err = os.OpenFile(logFN, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	logger := log.New(logFile, "GEDCOM Logger", log.LstdFlags)
	logger.Println("Importing Gedcom 7 configs.")

	AddValidTag(TagHEAD)
	AddValidTag(TagTRLR)
	AddValidTag(TagCONT)
	logger.Println("Added", TagHEAD, TagTRLR, "and", TagCONT, "tags.")

	files, err := abnfFS.ReadDir(abnfDir)
	if err != nil {
		log.Fatal("unable to open abnf folder", abnfDir, err)
	}

	for _, fn := range files {
		data, err := abnfFS.ReadFile(abnfDir + "/" + fn.Name())
		if err != nil {
			logger.Println("Error accessing ", fn.Name(), err.Error())
			continue
		}
		logger.Printf("Processing %s.\n", fn.Name())

		name := strings.Split(fn.Name(), "-")[0]
		switch name {
		case "enum":
			// These values are extracted from enumSets.
			continue
		case "month":
			// These values are extracted from calendars.
			continue
		case "enumset":
			if es, err := loadEnumSet(data); err != nil {
				logger.Printf("Error parsing %s as enumSet\n%s\n", fn.Name(), err.Error())
			} else {
				enumSets[es.URI] = es
			}
		case "cal":
			cm, err := loadCal(data)
			if err != nil {
				logger.Printf("Error parsing %s as calendar: %s\n", fn.Name(), err.Error())
			} else {
				calendars[cm.Cal] = cm
			}
			logger.Printf("Added calendar %s.\n", cm.Cal)
		case "type":
			tm, err := loadType(data)
			if err != nil {
				logger.Printf("Error parsing %s as type: %s\n", fn.Name(), err.Error())
			} else {
				types[tm.Type] = tm
			}
			logger.Printf("Added type %s.\n", tm.Type)
		default:
			t, err := loadTag(data)
			if err != nil {
				logger.Printf("Error parsing %s as default: %s\n", fn.Name(), err.Error())
			} else {
				if name == "ord" {
					// special case for LDS Ordinance tags
					t.FullTag = "ord-" + t.FullTag
				}
				if name == "record" {
					// record types have no superstructures and distinct substructure lists.
					// They apply only if level is 0.
					t.FullTag = "record-" + t.FullTag
				}
				tags[t.FullTag] = t
			}
			logger.Printf("Loaded tag %s.\n", t.FullTag)
		}
	}

	for key, tag := range tags {
		if tag.EnumSetName != "" {
			if es, ok := enumSets[tag.EnumSetName]; !ok {
				logger.Printf("No matching tag for %s to %s.\n", key, tag.EnumSetName)
			} else {
				tag.EnumSet = es
				tags[key] = tag
				logger.Printf("Added enumset %s.\n", es.FullTag)
			}
		}
	}

	baseline.tags = tags
	baseline.calendars = calendars
	baseline.types = types
	baseline.enumSets = enumSets
}

// deserializeYAML populates v with the contents of the first document in the YAML text.
// It skips everything prior to the first document delimiter.
func deserializeYAML[T any](data []byte, v *T) error {
	pos := bytes.Index(data, []byte("---"))
	if pos == -1 {
		pos = 0
	}

	decoder := yaml.NewDecoder(bytes.NewBuffer(data[pos:]))
	for {
		if err := decoder.Decode(v); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("Document decode failed: %w", err)
		}
	}
	return nil
}
