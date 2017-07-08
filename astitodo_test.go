package astitodo_test

import (
	"testing"

	"github.com/asticode/go-astitodo"
	"github.com/stretchr/testify/assert"
)

func TestProcessPath(t *testing.T) {
	expected := astitodo.TODOs{
		{
			Line:     5,
			Message:  []string{"Here is a", "multi line todo"},
			Filename: "tests/level1/level2.go",
		},
		{
			Line:     11,
			Assignee: "my.weird-email_address+1@email.com",
			Message:  []string{"This is a named TODO"},
			Filename: "tests/level1/level2.go",
		},
		{
			Line:     16,
			Assignee: "quentin renard",
			Message:  []string{"Here is another", "multi line todo"},
			Filename: "tests/level1/level2.go",
		},
		{
			Line:     8,
			Message:  []string{"Is it really your second function?"},
			Filename: "tests/level1.go",
		},
		{
			Line:     11,
			Message:  []string{"This is a tabbed TODO"},
			Filename: "tests/level1.go",
		},
		{
			Line:     12,
			Message:  []string{"this a second todo in the same comment group"},
			Filename: "tests/level1.go",
		},
		{
			Line:     19,
			Message:  []string{"Please delete me!"},
			Filename: "tests/level1.go",
		},
	}

	todos, err := astitodo.Extract("tests", "tests/excluded.go")
	assert.NoError(t, err)
	assert.Len(t, todos, 7)
	assert.Equal(t, expected, todos)
}
