package astitodo

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Vars
var (
	regexpAssignee = regexp.MustCompile("^\\([\\w \\._\\+\\-@]+\\)")
)

// TODOs represents a set of todos
type TODOs []*TODO

// TODO represents a todo
type TODO struct {
	Assignee string
	Filename string
	Line     int
	Message  []string
}

// Extract walks through an input path and extracts TODOs from all files it encounters
func Extract(path string, excludedPaths ...string) (todos TODOs, err error) {
	err = todos.extract(path, excludedPaths...)
	return
}

func (todos *TODOs) extract(path string, excludedPaths ...string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// Skip excluded paths
		for _, p := range excludedPaths {
			if p == path {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// Skip vendor and all directories beginning with a .
		if info.IsDir() && (info.Name() == "vendor" || (len(info.Name()) > 1 && info.Name()[0] == '.')) {
			return filepath.SkipDir
		}

		// Only process go files
		if !info.IsDir() && filepath.Ext(path) != ".go" {
			return nil
		}

		// Everything is fine here, extract if path is a file
		if !info.IsDir() {
			if err = todos.extractFile(path); err != nil {
				return err
			}
		}
		return nil
	})
}

func (todos *TODOs) extractFile(filename string) (err error) {
	// Parse file and create the AST
	var fset = token.NewFileSet()
	var f *ast.File
	if f, err = parser.ParseFile(fset, filename, nil, parser.ParseComments); err != nil {
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
				if strings.HasPrefix(t, "//") || strings.HasPrefix(t, "/*") || strings.HasPrefix(t, "*/") {
					t = strings.TrimSpace(t[2:])
				}

				// To do found
				if strings.HasPrefix(strings.ToLower(t), "todo") {
					// Init to do
					todo = &TODO{Filename: filename, Line: fset.Position(c.Slash).Line + i}
					t = strings.TrimSpace(t[4:])

					// Look for assignee
					if todo.Assignee = regexpAssignee.FindString(t); todo.Assignee != "" {
						t = strings.TrimSpace(t[len(todo.Assignee):])
						todo.Assignee = todo.Assignee[1 : len(todo.Assignee)-1]
					}

					// Append text
					todo.Message = append(todo.Message, t)
					*todos = append(*todos, todo)
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
