package main

import (
	"os"
	"strings"
)

func printHelp() {
	println("Syntax: golf source_file target_file [-v] [-- 'initialization']")
	println()
}

func printTitle() {
	println()
	println("Go Language Facilitator v" + version)
	println()
}

func checkFail(err error, message string) {
	if err != nil {
		println("Error: " + err.Error())
		fail(message)
	}
}

func fail(message string) {
	println(message)
	os.Exit(1)
}

func argumentError(arg string, verbose bool) {
	if !verbose {
		printTitle()
	}
	printHelp()
	fail("Unsupported argument found: " + arg)
}

func processArguments(args []string) (source string, dest string, initialize string, verbose bool, include []string) {
	initMode := false
	for i, arg := range args {
		if initMode {
			initialize += arg
			if i != len(args)-1 {
				initialize += " "
			}
		} else if arg == "--" {
			initMode = true
		} else if strings.HasPrefix(arg, "-") {
			if arg[1:] == "verbose" || arg[1:] == "v" {
				verbose = true
			} else if strings.HasPrefix(arg[1:], "include=") || strings.HasPrefix(arg[1:], "i=") {
				include = append(include, getParseIncludeFilename(arg[1:]))
			} else {
				argumentError(arg, verbose)
			}
		} else {
			if len(source) == 0 {
				source = arg
			} else if len(dest) == 0 {
				dest = arg
			} else {
				argumentError(arg, verbose)
			}
		}
	}
	if len(source) == 0 || len(dest) == 0 {
		if !verbose {
			printTitle()
		}
		printHelp()
		fail("Please provide a source and target file.")
	}
	return
}

func getParseIncludeFilename(arg string) string {
	if strings.HasPrefix(arg, "i=") {
		arg = arg[len("i="):]
	} else if strings.HasPrefix(arg, "include=") {
		arg = arg[len("include="):]
	} else {
		fail("Invalid include argument provided")
	}
	return arg
}
