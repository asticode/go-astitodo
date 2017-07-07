package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/asticode/go-astitodo"
)

type flagArray []string

func (f *flagArray) String() string {
	return strings.Join(*f, ",")
}

func (f *flagArray) Set(i string) error {
	*f = append(*f, i)
	return nil
}

var myFlags flagArray

// Flags
var (
	assignee = flag.String("a", "", "Only TODOs assigned to this username will be displayed")
	format   = flag.String("f", "text", "Format to use when outputting TODOs (supported formats: text, csv)")
	exclude  = flagArray{}
)

func main() {
	// Parse flags
	flag.Var(&exclude, "e", "Path that will be excluded from the process")
	flag.Parse()

	// Loop through paths
	for _, path := range flag.Args() {
		// Process path
		var todos astitodo.TODOs
		var err error
		if todos, err = astitodo.Extract(path, exclude...); err != nil {
			log.Fatal(err)
		}

		// Filter results for assignee
		var filteredTODOs = todos
		if *assignee != "" {
			filteredTODOs = astitodo.TODOs{}

			for _, t := range todos {
				if *assignee == t.Assignee {
					filteredTODOs = append(filteredTODOs, t)
				}
			}
		}

		// Handle selected format
		switch *format {
		case "text":
			fmt.Print(formatText(filteredTODOs))
		case "csv":
			fmt.Print(formatCSV(filteredTODOs))
		default:
			log.Fatalf("unsupported format: %s", *format)
		}
	}
}
