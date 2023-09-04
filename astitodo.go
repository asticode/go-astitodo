package astitodo

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Vars
var (
	regexpAssignee = regexp.MustCompile("^\\([\\w \\._\\+\\-@]+\\)")

	todoIdentifiers = []string{"TODO", "FIXME"}
)

// TODOContainer represents a set of todos
type TODOContainer struct {
	Path  string
	TODOs []*TODO
}

// TODO represents a todo
type TODO struct {
	Assignee string
	Filename string
	Line     int
	Message  []string
}

// Extract walks through an input path and extracts TODOContainer from all files it encounters
func Extract(path string, excludedPaths ...string) (todos TODOContainer, err error) {
	err = todos.extract(path, excludedPaths...)
	return
}

func (todos *TODOContainer) extract(path string, excludedPaths ...string) error {
	todos.Path = path
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// Process error
		if err != nil {
			return err
		}

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

func (todos *TODOContainer) extractFile(filename string) (err error) {
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
				if length, isTodo := isTodoIdentifier(t); isTodo {
					// Init to do
					todo = &TODO{Filename: filename, Line: fset.Position(c.Slash).Line + i}
					t = strings.TrimSpace(t[length:])
					if strings.HasPrefix(t, ":") {
						t = strings.TrimLeft(t, ":")
						t = strings.TrimSpace(t)
					}

					// Look for assignee
					if todo.Assignee = regexpAssignee.FindString(t); todo.Assignee != "" {
						t = strings.TrimSpace(t[len(todo.Assignee):])
						if strings.HasPrefix(t, ":") {
							t = strings.TrimLeft(t, ":")
							t = strings.TrimSpace(t)
						}
						todo.Assignee = todo.Assignee[1 : len(todo.Assignee)-1]
					}

					// Append text
					todo.Message = append(todo.Message, t)
					todos.TODOs = append(todos.TODOs, todo)
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

func isTodoIdentifier(s string) (int, bool) {
	for _, indent := range todoIdentifiers {
		if strings.HasPrefix(strings.ToUpper(s), indent) {
			return len(indent), true
		}
	}
	return 0, false
}

// AssignedTo returns TODOContainer which are assigned to the specified assignee
func (todos TODOContainer) AssignedTo(assignees ...string) (filteredTODOs TODOContainer) {
	for _, t := range todos.TODOs {
		for _, assignee := range assignees {
			if assignee == t.Assignee {
				filteredTODOs.TODOs = append(filteredTODOs.TODOs, t)
			}
		}
	}

	return
}

// WriteHTML writes the TODOContainer markdown-formatted to the specified writer
func (todos TODOContainer) WriteHTML(w io.Writer) (err error) {

	var tocBuffer bytes.Buffer
	var contentBuffer bytes.Buffer

	_, err = io.WriteString(w, fmt.Sprintf("<h1>TODOs for %s</h1>\n\n", todos.Path))

	if err != nil {
		return err
	}

	if len(todos.TODOs) == 0 {
		_, err = io.WriteString(w, "<ul><li>NONE</li></ul>")
		return err
	}

	tocBuffer.WriteString("\n<ul id=\"toc\">\n")
	contentBuffer.WriteString("\n<ul id=\"content\">\n")
	i := 1
	for _, t := range todos.TODOs {

		tocBuffer.WriteString(fmt.Sprintf("<li><a href=\"#%d\">%s:%d</a></li>\n", i, t.Filename, t.Line))

		contentBuffer.WriteString("<li>")
		contentBuffer.WriteString(fmt.Sprintf("<h2><a id=\"%d\">%s:%d</a></h2>\n", i, t.Filename, t.Line))
		if t.Assignee != "" {
			contentBuffer.WriteString(fmt.Sprintf("<div class=\"assignee\">Assignee: %s</div>\n", t.Assignee))
		}
		contentBuffer.WriteString("<pre class=\"todo\">\n")
		for _, m := range t.Message {
			contentBuffer.WriteString(fmt.Sprintf("%s\n", m))
		}
		contentBuffer.WriteString("</pre>\n")
		contentBuffer.WriteString("</li>")
		i++
	}
	tocBuffer.WriteString("\n</ul>\n")
	contentBuffer.WriteString("\n</ul>\n")

	_, err = io.WriteString(w, fmt.Sprintf("<html><head><title>Todos for %s</title><link rel=\"stylesheet\" type=\"text/css\" href=\"todos.css\" /></head><body>%s<hr>%s</body></html>",
		todos.Path,
		tocBuffer.String(),
		contentBuffer.String()))

	return err
}

// WriteMarkdown writes the TODOContainer markdown-formatted to the specified writer
func (todos TODOContainer) WriteMarkdown(w io.Writer) (err error) {

	var tocBuffer bytes.Buffer
	var contentBuffer bytes.Buffer

	_, err = io.WriteString(w, fmt.Sprintf("# TODOs for %s\n\n", todos.Path))

	if err != nil {
		return err
	}

	if len(todos.TODOs) == 0 {
		_, err = io.WriteString(w, " - NONE")
		return err
	}

	for _, t := range todos.TODOs {
		header := fmt.Sprintf("%s:%d", t.Filename, t.Line)
		tocBuffer.WriteString(fmt.Sprintf(" - [%s:%d](#%s)\n", t.Filename, t.Line, header))

		contentBuffer.WriteString(fmt.Sprintf("## %s\n\n", header))
		if t.Assignee != "" {
			contentBuffer.WriteString(fmt.Sprintf("Assignee: `%s`\n", t.Assignee))
		}
		contentBuffer.WriteString("```\n")
		for _, m := range t.Message {
			contentBuffer.WriteString(fmt.Sprintf("%s\n", m))
		}
		contentBuffer.WriteString("```\n")
		contentBuffer.WriteString("\n---\n")
	}

	_, err = io.WriteString(w, fmt.Sprintf("%s\n\n---\n\n%s", tocBuffer.String(), contentBuffer.String()))

	return err
}

// WriteText writes the TODOContainer as text to the specified writer
func (todos TODOContainer) WriteText(w io.Writer) (err error) {
	for _, t := range todos.TODOs {
		if t.Assignee != "" {
			if _, err = io.WriteString(w, fmt.Sprintf("Assignee: %s\n", t.Assignee)); err != nil {
				return
			}
		}

		if _, err = io.WriteString(w, fmt.Sprintf("Message: %s\nFile:%s:%d\n\n", strings.Join(t.Message, "\n"), t.Filename, t.Line)); err != nil {
			return
		}
	}

	return
}

// WriteCSV writes the TODOContainer as CSV to the specified writer
// The columns are "Filename", "Line", "Assignee" and "Message" (which can contain newlines)
func (todos TODOContainer) WriteCSV(w io.Writer) (err error) {
	var c = csv.NewWriter(w)

	// Write the headings for the document
	if err = c.Write([]string{"Path", "Filename", "Line", "Assignee", "Message"}); err != nil {
		return
	}

	for _, t := range todos.TODOs {
		err = c.Write([]string{
			todos.Path,
			t.Filename,
			strconv.Itoa(t.Line),
			t.Assignee,
			strings.Join(t.Message, "\n"),
		})

		if err != nil {
			return
		}
	}

	c.Flush()

	return
}

// WriteJSON writes the TODOContainer as JSON to the specified writer
func (todos TODOContainer) WriteJSON(w io.Writer) (err error) {
	enc := json.NewEncoder(w)
	err = enc.Encode(todos)
	return
}
