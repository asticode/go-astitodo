package main_test

import (
	"testing"

	main "github.com/asticode/gotodo"
	"github.com/stretchr/testify/assert"
)

func TestProcessPath(t *testing.T) {
	todos, err := main.ProcessPath("./tests")
	assert.NoError(t, err)
	assert.Len(t, todos, 3)
	assert.Equal(t, &main.TODO{Line: 10, Message: []string{"Rewrite this entirely", "because it kinda sucks"}, Path: "./tests/level1/level2.go"}, todos[0])
	assert.Equal(t, &main.TODO{Line: 8, Message: []string{"Is it really your second function?"}, Path: "./tests/level1.go"}, todos[1])
	assert.Equal(t, &main.TODO{Line: 11, Message: []string{"This is a tabbed TODO"}, Path: "./tests/level1.go"}, todos[2])
}
