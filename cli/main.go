package main

/********************************************************************************
This, or any other application, should only read existing GEDCOM files. If any
actions are taken to modify the input file, the output should be written to a
new file. This will allow the original file to be preserved in case of errors.
********************************************************************************/

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"go-gedcom/pkg/gedcom7"
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
	f, err := os.Open(fn)
	if err != nil {
		fmt.Println("Error accessing ", fn)
		log.Fatal(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	s := bufio.NewScanner(f)
	doc := gedcom7.NewDocument(s, gedcom7.WithMaxDeprecatedTags("5.5.1"))

	fmt.Printf("Processed %d records with %d warnings.\n", doc.Len(), len(doc.Warnings))
	for _, v := range doc.Warnings {
		fmt.Printf("%s\t%s\n", v.Line, v.Message)
	}
}

func recordCount(docPath, fn string) {
	fn = docPath + fn
	fmt.Printf("Processing file '%s'\n", fn)
	exitWithHelp()
}
