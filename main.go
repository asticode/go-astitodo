package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

// Vars
var (
	// Flags
	assignee = flag.String("a", "", "Only TODOs assigned to this username will be displayed")
	noSkip   = flag.Bool("no-skip", false, "If true, no directories are skipped")
	verbose  = flag.Bool("v", false, "If true, then verbose")

	// Others
	regexpAssignee = regexp.MustCompile("^\\([\\w ]+\\)")
)

func main() {
	// Parse flags
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
				fmt.Printf("File: %s:%d\n\n", t.Path, t.Line)
			}
		}
	}
}

// TODO represents a todo
type TODO struct {
	Assignee string
	Line     int
	Message  []string
	Path     string
}

// ProcessPath processes a path which can be either a directory or a file
func ProcessPath(path string) (todos []*TODO, err error) {
	// Log
	if *verbose {
		log.Printf("Processing path %s\n", path)
	}

	// Stat path
	var file os.FileInfo
	if file, err = os.Stat(path); err != nil {
		return
	}

	// Directory
	if file.IsDir() {
		// Skip some directories
		if !*noSkip && (file.Name() == "vendor" || (len(file.Name()) > 1 && file.Name()[0] == '.')) {
			if *verbose {
				log.Printf("Skipping directory %s\n", path)
			}
			return
		}

		// Read dir
		var files []os.FileInfo
		if files, err = ioutil.ReadDir(path); err != nil {
			return
		}

		// Process each file
		var fileTODOs []*TODO
		for _, file := range files {
			if fileTODOs, err = ProcessPath(path + string(os.PathSeparator) + file.Name()); err != nil {
				return
			}
			todos = append(todos, fileTODOs...)
		}
	} else {
		todos, err = ProcessFile(path)
	}
	return
}

// ProcessFile processes a file and extract its TODOs
func ProcessFile(path string) (todos []*TODO, err error) {
	// Open file
	var file *os.File
	if file, err = os.Open(path); err != nil {
		return
	}
	scanner := bufio.NewScanner(file)

	// Scan
	var line string
	var lineCount int
	var todo *TODO
	var TODOFound bool
	for scanner.Scan() {
		// Fetch line
		line = strings.TrimSpace(scanner.Text())
		lineCount++

		// To do found
		if len(line) >= 7 && line[:7] == "// TODO" {
			todo = &TODO{
				Line: lineCount,
				Path: path,
			}
			line = strings.TrimSpace(line[7:])
			if todo.Assignee = regexpAssignee.FindString(line); todo.Assignee != "" {
				line = strings.TrimSpace(line[len(todo.Assignee):])
				todo.Assignee = todo.Assignee[1 : len(todo.Assignee)-1]
			}
			todo.Message = append(todo.Message, line)
			todos = append(todos, todo)
			TODOFound = true
		} else if TODOFound && len(line) >= 4 && line[:3] == "// " {
			todo.Message = append(todo.Message, strings.TrimSpace(line[3:]))
		} else {
			TODOFound = false
		}
	}
	return
}
