package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"go-gedcom/pkg/gedcom7"
)

type logData struct {
	lineNo int
	line   string
	err    string
}

func main() {
	docPath := "./"

	if p, ok := os.LookupEnv("G7DOCPATH"); ok {
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
		inspectCmd.Parse(os.Args[2:])
		inspect(docPath, *in)
	case "records":
		recordsCmd.Parse(os.Args[2:])
		recordCount(docPath, *rec)
	default:
		exitWithHelp()
	}
	flag.Parse()

}

func exitWithHelp() {
	fmt.Println("expected command")
	fmt.Println("Currently supported")
	fmt.Println("\tinspect -input {gedfilename}")
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
