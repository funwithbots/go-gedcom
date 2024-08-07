package main

/********************************************************************************
This, or any other application, should only read existing GEDCOM files. If any
actions are taken to modify the input file, the output should be written to a
new file. This will allow the original file to be preserved in case of errors.
********************************************************************************/

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/funwithbots/go-gedcom/pkg/gedcom"
	"github.com/funwithbots/go-gedcom/pkg/gedcom7"
)

func main() {
	// docPath is overridden by the environment variable GEDCOM_DOC_PATH
	docPath := "./"

	if p, ok := os.LookupEnv("GEDCOM_DOC_PATH"); ok {
		docPath = p
	}

	inspectCmd := flag.NewFlagSet("inspect", flag.ExitOnError)
	in := inspectCmd.String("in", "", "input file")

	recordsCmd := flag.NewFlagSet("records", flag.ExitOnError)
	rec := recordsCmd.String("rec", "", "input file")

	if len(os.Args) < 2 {
		exitWithHelp()
	}

	switch os.Args[1] {
	case "inspect":
		if err := inspectCmd.Parse(os.Args[2:]); err != nil {
			log.Fatal(err)
		}
		inspect(docPath, *in)
	case "records":
		if err := recordsCmd.Parse(os.Args[2:]); err != nil {
			log.Fatal(err)
		}
		recordCount(docPath, *rec)
	default:
		exitWithHelp()
	}
	flag.Parse()

}

func exitWithHelp() {
	fmt.Println("expected command")
	fmt.Println("Currently supported")
	fmt.Println("\tinspect -in {gedfilename}")
	os.Exit(1)
}

func inspect(docPath, fn string) {
	fn = docPath + fn
	fmt.Printf("Processing file '%s'\n", fn)
	doc, err := loadGedcom7FromFile(fn)
	if err != nil {
		fmt.Println("Error accessing ", fn)
		log.Fatal(err)
	}

	warnings := doc.Warnings()
	fmt.Printf("Processed %d records with %d warnings.\n", doc.Len(), len(warnings))
	for _, v := range warnings {
		fmt.Println(v)
	}
}

func recordCount(docPath, fn string) {
	fn = docPath + fn
	fmt.Printf("Processing file '%s'\n", fn)
	exitWithHelp()
}

func loadGedcom7FromFile(file string) (gedcom.Document, error) {
	fmt.Println("Loading file ", file)
	return gedcom7.NewDocumentFromFile(file, gedcom7.WithMaxDeprecatedTags("5.5.1"))
}
