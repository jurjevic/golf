package main

import (
	"os"
)

const (
	version = "0.9.0-alpha" // ###  ftoken[2] = "\"" + NewVersion + "\""; Join(ftoken, " ")
)

// golf is a line facilitator which acts like preprocessor based on the go language syntax.
func main() {
	args := os.Args
	if len(args) == 1 {
		printTitle()
		printHelp()
		os.Exit(1)
	}
	source, dest, initialize, verbose, include := processArguments(args[1:])
	if verbose {
		printTitle()
	}
	checkFile(source)
	golf := NewGolf(os.Getenv("GOPATH"))
	for _, inc := range include {
		golf.include(inc, verbose)
	}
	golf.processInitialize(initialize, verbose)
	content := golf.processFile(source, verbose)
	saveFile(dest, content)
}
