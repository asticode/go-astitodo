# About

I use inline TODOs a lot in my code and needed a way to print all TODOs easily, so I created `gotodo`

# Install

Run 

    go get github.com/asticode/gotodo && go install
    
# Usage

    Usage of gotodo:
        -a string
            Only TODOs assigned to this username will be displayed
        -no-skip
            If true, no directories are skipped
        -v    If true, then verbose
        
# Formatting

A todo is formatted this way:

    // TODO<line 1>
    // <line 2>
    // ...
       
You can also add an assignee:

    // TODO(this is the assignee)<message>
        
# Examples
## Basic

Assume the following file:

    package mypackage
    
    // TODO Damn this package seems useless
    
    // Here is a dummy comment
    // TODO(asticode) This variable should be dropped
    var myvariable int
    
    // TODO Damn this function should be rewritten
    // Or maybe it should be dropped as well
    func UselessFunction() {
    	var a = 1
    	a++
    }
    
Running

    gotodo <path to file>
    
will give

    Message: Damn this package seems useless
    File: mypackage/main.go:3
    
    Assignee: asticode
    Message: This variable should be dropped
    File: mypackage/main.go:6
    
    Message: Damn this function should be rewritten
    Or maybe it should be dropped  as well
    File: mypackage/main.go:9
    
## Filter by assignee

Running

    gotodo -a asticode <path to file>
    
will give

    Assignee: asticode
    Message: This variable should be dropped
    File: mypackage/main.go:6
    
# Contributions

You know you want to =D