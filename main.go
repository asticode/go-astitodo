package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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

// Vars
var (
	// Flags
	assignee = flag.String("a", "", "Only TODOs assigned to this username will be displayed")
	exclude  = flagArray{}
	verbose  = flag.Bool("v", false, "If true, then verbose")

	// Others
	regexpAssignee = regexp.MustCompile("^\\([\\w \\._\\+\\-@]+\\)")
)

func main() {
	// Parse flags
	flag.Var(&exclude, "e", "Path that will be excluded from the process")
	flag.Parse()

	// Loop through paths
	var todos []*TODO
	var err error
	for _, path := range flag.Args() {
		// Process path
		if todos, err = ProcessPath(path); err != nil {
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

// TODO represents a todo
type TODO struct {
	Assignee string
	Filename string
	Line     int
	Message  []string
}

// ProcessPath processes a path
func ProcessPath(path string) (todos []*TODO, err error) {
	// Walk the path
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// Log
		if *verbose {
			log.Printf("Processing path %s\n", path)
		}

		// Check exclude list
		for _, p := range exclude {
			if p == path {
				if *verbose {
					log.Printf("Skipping path %s\n", path)
				}
				if info.IsDir() {
					return filepath.SkipDir
				} else {
					return nil
				}
			}
		}

		// Check whether file is a dir
		if info.IsDir() {
			// Skip vendor and all directories beginning with a .
			if info.Name() == "vendor" || (len(info.Name()) > 1 && info.Name()[0] == '.') {
				if *verbose {
					log.Printf("Skipping path %s\n", path)
				}
				return filepath.SkipDir
			}
		} else {
			// Only process go files
			if filepath.Ext(path) != ".go" {
				if *verbose {
					log.Printf("Skipping path %s\n", path)
				}
				return nil
			}

			// Process file and add the todos
			var t []*TODO
			if t, err = ProcessFile(path); err != nil {
				return err
			}
			todos = append(todos, t...)
		}
		return nil
	})
	return
}

// ProcessFile processes a file and extract its TODOs
func ProcessFile(path string) (todos []*TODO, err error) {
	// Parse file and create the AST
	var fset = token.NewFileSet()
	var f *ast.File
	if f, err = parser.ParseFile(fset, path, nil, parser.ParseComments); err != nil {
		return
	}

	// Loop in comment groups
	for _, cg := range f.Comments {
		// Loop in comments
		var todo *TODO
		var TODOFound bool
		for _, c := range cg.List {
			// Loop in lines
			for i, l := range strings.Split(c.Text, "\n") {
				// Init text
				var t = strings.TrimSpace(l)
				if len(t) >= 2 && (t[:2] == "//" || t[:2] == "/*" || t[:2] == "*/") {
					t = strings.TrimSpace(t[2:])
				}

				// To do found
				if len(t) >= 4 && strings.ToLower(t[:4]) == "todo" {
					// Init to do
					todo = &TODO{Filename: path, Line: fset.Position(c.Slash).Line + i}
					t = strings.TrimSpace(t[4:])

					// Look for assignee
					if todo.Assignee = regexpAssignee.FindString(t); todo.Assignee != "" {
						t = strings.TrimSpace(t[len(todo.Assignee):])
						todo.Assignee = todo.Assignee[1 : len(todo.Assignee)-1]
					}

					// Append text
					todo.Message = append(todo.Message, t)
					todos = append(todos, todo)
					TODOFound = true
				} else if TODOFound && len(t) > 0 {
					todo.Message = append(todo.Message, t)
				} else {
					TODOFound = false
				}
			}
		}
	}
	return
}
