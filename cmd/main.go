package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/asticode/gotodo"
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
	exclude  = flagArray{}
)

func main() {
	// Parse flags
	flag.Var(&exclude, "e", "Path that will be excluded from the process")
	flag.Parse()

	// Loop through paths
	for _, path := range flag.Args() {
		// Process path
		var todos todo.TODOs
		var err error
		if todos, err = todo.Extract(path, exclude...); err != nil {
			log.Fatal(err)
		}

		// Display results
		for _, t := range todos {
			if *assignee == "" || *assignee == t.Assignee {
				if t.Assignee != "" {
					fmt.Printf("Assignee: %s\n", t.Assignee)
				}
				fmt.Printf("Message: %s\n", strings.Join(t.Message, "\n"))
				fmt.Printf("File: %s:%d\n\n", t.Filename, t.Line)
			}
		}
	}
}
