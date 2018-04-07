package tests

func MyFirstFunction() {
	var a = 1
	a++
}

// TODO Is it really your second function?
func MySecondFunction() {
	// This is a dummy comment
	// TODO This is a tabbed TODO
	// TODO this a second todo in the same comment group

	// This is another dummy comment
	var b = 2
	b++
}

// FIXME Please delete me!
var DeleteMe = 5

// TODO(asticode) I should be false
var Oops = true

// TODO(astitodo) Something else comes here
var SomethingElse = "Something else"

// TODO: I can use colons to signal the todo.
var WithColons = "Something else"

// TODO(astitodo): It also works with assignee.
var WithColons2 = "Something else"
