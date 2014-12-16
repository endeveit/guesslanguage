# Guess the natural language of a text

This is a Go version of python [guess-language](http://code.google.com/p/guess-language>).

## Usage

Install in your `${GOPATH}` using `go get -u github.com/endeveit/guesslanguage`

Then call it:
```go
package main

import (
	"fmt"
	"github.com/endeveit/guesslanguage"
)

func main() {
	lang, err := guesslanguage.Guess("This is a test of the language checker.")

	// Output:
	// en
	if err != nil {
		fmt.Println(lang)
	}
}
```