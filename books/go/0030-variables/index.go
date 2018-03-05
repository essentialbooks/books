package main

import (
	"bytes"
	"fmt"
	"io"
)

// :show start
// declaration of a single top-level variable
var topLevel int64 = 5

// grouping of multiple top-level declarations
var (
	intVal int            // value is initialized with zero-value
	str    string = "str" // assigning

	// functions are first-class values so can be assigned to variables
	// f is variable of type func(a int) string
	// it's uninitialized so is nil (zero-value for function variables)
	fn func(a int) string
)

func f() {
	// shorthand using local type inference
	// type of `i` is int and is infered from the value
	// note: this is not allowed at top-level
	i := 4

	// grouping inside a function
	var (
		i2 int
		s  string
	)

	// _ is like a variable whose value is discarded. It's called blank identifier.
	// Useful when we don't care about one of the values returned by a function
	_, err := io.Copy(dst, src) // don't care how many bytes were written
	// ...

	fmt.Printf("i: %d, i2: %d, s: %s, err: %v\n", i, i2, s, err)
}

// :show end

var (
	dst io.Writer = &bytes.Buffer{}
	src io.Reader = &bytes.Buffer{}
)

func main() {
	f()
}
