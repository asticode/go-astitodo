package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Constants
const (
	printFormat = `Message: %s
File: %s:%d

`
)

// Flags
var (
	verbose = flag.Bool("v", false, "If true, then verbose")
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
			fmt.Printf(printFormat, strings.Join(t.Message, "\n"), t.Path, t.Line)
		}
	}
}

// TODO represents a todo
type TODO struct {
	Line    int
	Message []string
	Path    string
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
		// Blacklist some directories
		if file.Name() == "vendor" || file.Name()[0] == '.' {
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
			TODOFound = true
			todo = &TODO{
				Line:    lineCount,
				Message: []string{strings.TrimSpace(line[7:])},
				Path:    path,
			}
			todos = append(todos, todo)
		} else if TODOFound && len(line) >= 4 && line[:3] == "// " {
			todo.Message = append(todo.Message, strings.TrimSpace(line[3:]))
		} else {
			TODOFound = false
		}
	}
	return
}
