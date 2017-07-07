package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"

	"github.com/asticode/go-astitodo"
)

func formatText(todos astitodo.TODOs) string {
	var b bytes.Buffer

	// Append each TODO to the buffer, separated by two newlines
	for _, t := range todos {
		if t.Assignee != "" {
			b.WriteString(fmt.Sprintf("Assignee: %s\n", t.Assignee))
		}

		b.WriteString(fmt.Sprintf("Message: %s\nFile: %s:%d\n\n", strings.Join(t.Message, "\n"), t.Filename, t.Line))
	}

	return b.String()
}

func formatCSV(todos astitodo.TODOs) string {
	var b bytes.Buffer
	var w = csv.NewWriter(&b)

	// Write the headings for the document
	w.Write([]string{"Filename", "Line", "Assignee", "Message"})

	for _, t := range todos {
		w.Write([]string{
			t.Filename,
			strconv.Itoa(t.Line),
			t.Assignee,
			strings.Join(t.Message, "\n"),
		})
	}

	w.Flush()

	return b.String()
}
