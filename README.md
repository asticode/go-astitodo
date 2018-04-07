[![GoReportCard](http://goreportcard.com/badge/github.com/asticode/go-astitodo)](http://goreportcard.com/report/github.com/asticode/go-astitodo)
[![GoDoc](https://godoc.org/github.com/asticode/go-astitodo?status.svg)](https://godoc.org/github.com/asticode/go-astitodo)
[![GoCoverage](https://cover.run/go/github.com/asticode/go-astitodo.svg)](https://cover.run/go/github.com/asticode/go-astitodo)
[![Travis](https://travis-ci.org/asticode/go-astitodo.svg?branch=master)](https://travis-ci.org/asticode/go-astitodo#)

This is a Golang library and CLI to parse TODOs in your GO code.

It parses the comments from the AST and extract their TODOs. It can provide valuable information such as the TODO's assignee which can be filtered afterwards.

Most IDEs allow parsing TODOs but they usually have problems with multi line TODOs, can't parse assignees, etc.

This is also a good start for people who want to use AST.

# Installation

Run

    $ go get -u github.com/asticode/go-astitodo/...

# Usage

    Usage of go-astitodo:
        -a string
            Only TODOs assigned to this username(s) will be displayed.( user OR user,anotheruser)
        -e
            Path that will be excluded from the process
        -f string
            Format to use when outputting TODOs (supported formats: text, csv) (default "text")
        -o string
            Destination for output (can be stdout, stderr or a file) (default "stdout")
        -v  If true, then verbose

# Formatting

A todo is formatted this way:

```go
    // TODO<line 1>
    // <line 2>
    // ...
```

You can also add an assignee:

```go
    // TODO(this is the assignee)<message>
```

# Examples
## Basic

Assume the following file:

```go
    package mypackage

    // TODO Damn this package seems useless

    // Here is a dummy comment
    // TODO(asticode) This variable should be dropped
    var myvariable int

    // TODO(username) This should be renamed
    var oops bool

    // TODO Damn this function should be rewritten
    // Or maybe it should be dropped as well
    func UselessFunction() {
    	var a = 1
    	a++
    }
```

Running

    go-astitodo <paths to files or dirs>

will give

    Message: Damn this package seems useless
    File: mypackage/main.go:3

    Assignee: asticode
    Message: This variable should be dropped
    File: mypackage/main.go:6

    Assignee: username
    Message: This variable should be renamed
    File: mypackage/main.go:9

    Message: Damn this function should be rewritten
    Or maybe it should be dropped  as well
    File: mypackage/main.go:12

## Filter by assignee

Running

    go-astitodo -a asticode <paths to files or dirs>

will output

    Assignee: asticode
    Message: This variable should be dropped
    File: mypackage/main.go:6

### Filter by multiple asignees

Running

    astitodo -a user,anotheruser <paths to files or dirs>

will output

    Assignee: asticode
    Message: This variable should be dropped
    File: mypackage/main.go:6

    Assignee: username
    Message: This variable should be renamed
    File: mypackage/main.go:9

## Exclude paths

You can exclude paths by running

    go-astitodo -e path/to/exclude/1 -e path/to/exclude/2 <paths to files or dirs>

## Change output format

You can output CSV by running

    go-astitodo -f csv <path to files or dirs>

You can output JSON by running

  astitodo -f json testdata/ | jq '.[] | select(.Assignee=="") | "\(.Filename):\(.Line)"'
  "testdata/excluded.go:3"
  "testdata/level1/level2.go:5"
  ...


## Output to a file

You can output to a file by running

    go-astitodo -o <path to output file> <path to files or dirs>

# Contributions

You know you want to =D

### hacking

```sh
go install github.com/asticode/go-astitodo/astitodo
go test -v ./...
```
