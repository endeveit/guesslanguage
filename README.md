# guesslanguage [![Build Status](https://travis-ci.org/endeveit/guesslanguage.svg?branch=master)](https://travis-ci.org/endeveit/guesslanguage)

This is a Go version of python [guess-language](http://code.google.com/p/guess-language>).

guesslanguage provides a simple way to detect the natural language of unicode string and detects over 60 languages listed in the [models](https://github.com/endeveit/guesslanguage/tree/master/models) directory.

## Supported Go versions

guesslanguage is regularly tested against Go 1.1, 1.2, 1.3 and tip.

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
