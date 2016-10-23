# About

I use inline TODOs a lot in my code and needed a way to print all TODOs easily, so I created `gotodo`

It parses the comments from the AST and extract their TODOs. It can provide valuable information such as the TODO's assignee which can be filtered afterwards

# Install

Run 

    go get github.com/asticode/gotodo && go install github.com/asticode/gotodo
    
# Usage

    Usage of gotodo:
        -a string
            Only TODOs assigned to this username will be displayed
        -e
            Path that will be excluded from the process
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

    gotodo <paths to files or dirs>
    
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

    gotodo -a asticode <paths to files or dirs>
    
will give

    Assignee: asticode
    Message: This variable should be dropped
    File: mypackage/main.go:6
 
## Exclude paths
    
You can exclude paths running

    gotodo -e path/to/exclude/1 -e path/to/exclude/2 <paths to files or dirs>
    
# Contributions

You know you want to =D