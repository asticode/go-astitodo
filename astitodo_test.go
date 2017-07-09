package astitodo_test

import (
	"testing"

	"bytes"

	"github.com/asticode/go-astitodo"
	"github.com/stretchr/testify/assert"
)

func TestExtract(t *testing.T) {
	expected := astitodo.TODOs{
		{
			Line:     5,
			Message:  []string{"Here is a", "multi line todo"},
			Filename: "testdata/level1/level2.go",
		},
		{
			Line:     11,
			Assignee: "my.weird-email_address+1@email.com",
			Message:  []string{"This is a named TODO"},
			Filename: "testdata/level1/level2.go",
		},
		{
			Line:     16,
			Assignee: "quentin renard",
			Message:  []string{"Here is another", "multi line todo"},
			Filename: "testdata/level1/level2.go",
		},
		{
			Line:     8,
			Message:  []string{"Is it really your second function?"},
			Filename: "testdata/level1.go",
		},
		{
			Line:     11,
			Message:  []string{"This is a tabbed TODO"},
			Filename: "testdata/level1.go",
		},
		{
			Line:     12,
			Message:  []string{"this a second todo in the same comment group"},
			Filename: "testdata/level1.go",
		},
		{
			Line:     19,
			Message:  []string{"Please delete me!"},
			Filename: "testdata/level1.go",
		},
	}

	todos, err := astitodo.Extract("testdata", "testdata/excluded.go")
	assert.NoError(t, err)
	assert.Len(t, todos, 7)
	assert.Equal(t, expected, todos)
}

func mockTODOs() astitodo.TODOs {
	return astitodo.TODOs{
		{Assignee: "1", Line: 1, Message: []string{"multi", "line"}, Filename: "filename-1"},
		{Line: 2, Message: []string{"no-assignee"}, Filename: "filename-1"},
		{Assignee: "2", Line: 3, Message: []string{"message-1"}, Filename: "filename-2"},
	}
}

func TestTODOs_AssignedTo(t *testing.T) {
	todos := mockTODOs()
	filteredTODOs := todos.AssignedTo("1")
	assert.Equal(t, astitodo.TODOs{{Assignee: "1", Line: 1, Message: []string{"multi", "line"}, Filename: "filename-1"}}, filteredTODOs)
}

func TestTODOs_WriteCSV(t *testing.T) {
	todos := mockTODOs()
	buf := &bytes.Buffer{}
	err := todos.WriteCSV(buf)
	assert.NoError(t, err)
	assert.NoError(t, err)
	assert.Equal(t, "Filename,Line,Assignee,Message\nfilename-1,1,1,\"multi\nline\"\nfilename-1,2,,no-assignee\nfilename-2,3,2,message-1\n", buf.String())
}

func TestTODOs_WriteText(t *testing.T) {
	todos := mockTODOs()
	buf := &bytes.Buffer{}
	err := todos.WriteText(buf)
	assert.NoError(t, err)
	assert.NoError(t, err)
	assert.Equal(t, "Assignee: 1\nMessage: multi\nline\nFile:filename-1:1\n\nMessage: no-assignee\nFile:filename-1:2\n\nAssignee: 2\nMessage: message-1\nFile:filename-2:3\n\n", buf.String())
}
