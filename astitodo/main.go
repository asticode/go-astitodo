package main

import (
	"flag"
	"io"
	"log"
	"os"
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
	output   = flag.String("o", "stdout", "Destination for output (can be stdout, stderr or a file)")
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
		if *assignee != "" {
			todos = todos.AssignedTo(*assignee)
		}

		var writer io.Writer

		// Convert selected output into a writer
		switch *output {
		case "stdout":
			writer = os.Stdout
		case "stderr":
			writer = os.Stderr
		default:
			if writer, err = os.OpenFile(*output, os.O_WRONLY|os.O_CREATE, 0666); err != nil {
				log.Fatal(err)
			}

			defer writer.(*os.File).Close()
		}

		// Handle selected format
		switch *format {
		case "text":
			if err = todos.WriteText(writer); err != nil {
				log.Fatal(err)
			}
		case "csv":
			if err = todos.WriteCSV(writer); err != nil {
				log.Fatal(err)
			}
		default:
			log.Fatalf("unsupported format: %s", *format)
		}
	}
}
