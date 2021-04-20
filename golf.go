package main

import (
	"bufio"
	"fmt"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

type golf struct {
	interpreter *interp.Interpreter
}

const (
	funcProcessingPrefix = " ### "
)

func NewGolf(gopath string) *golf {
	interpreter := interp.New(interp.Options{
		GoPath: gopath,
	})
	interpreter.Use(stdlib.Symbols)
	// default imports
	doDotImport(interpreter, "fmt")
	doDotImport(interpreter, "strings")
	doImport(interpreter, "os")
	builtIns(interpreter)
	return &golf{interpreter: interpreter}
}

func builtIns(interpreter *interp.Interpreter) {
	statements := []string{
		"var arg = make(map[string]interface{})",
		"func isSet(key string) bool { _, ok := arg[key]; return ok }",
	}
	for _, statement := range statements {
		_, err := interpreter.Eval(statement)
		checkFail(err, "Failed to initialize the built-in variables and function. Detail: "+statement)
	}
}

func doImport(interpreter *interp.Interpreter, pkg string) {
	_, err := interpreter.Eval(`import "` + pkg + `"`)
	checkFail(err, "Package not found: "+pkg)
}

func doDotImport(interpreter *interp.Interpreter, pkg string) {
	_, err := interpreter.Eval(`import . "` + pkg + `"`)
	checkFail(err, "Package not found: "+pkg)
}

func (g *golf) eval(call string, line string) (result string, ok bool, repeat bool, next int64, err error) {
	// disable the repeat instruction
	_, err = g.interpreter.Eval("repeat := false")
	if err != nil {
		return "", false, false, 0, err
	}
	_, err = g.interpreter.Eval("next := 0")
	if err != nil {
		return "", false, false, 0, err
	}
	_, err = g.interpreter.Eval("line := `" + line + "`")
	if err != nil {
		return "", false, false, 0, err
	}
	_, err = g.interpreter.Eval("token := Split(line, ` `)")
	if err != nil {
		return "", false, false, 0, err
	}
	resultValue, err := g.interpreter.Eval(call)
	result = resultValue.String()
	ok = resultValue.Kind() == reflect.String
	if err != nil {
		return "", false, false, 0, err
	}
	repeatValue, err := g.interpreter.Eval("repeat")
	repeat = repeatValue.Bool()
	nextValue, err := g.interpreter.Eval("next")
	next = nextValue.Int()
	return
}

func (g *golf) processInitialize(initialize string, verbose bool) {
	if verbose && len(initialize) > 0 {
		fmt.Printf("Initialize: '%s'\n", initialize)
	}
	_, err := g.interpreter.Eval(initialize)
	checkFail(err, "Initialization syntax error: '"+initialize+"'")
}

func (g *golf) processFile(filename string, verbose bool) string {
	output := ""
	file, err := os.Open(filename)
	checkFail(err, "File "+filename+" not openable!")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	// start processing
	lineCounter := 0
	nextLines := 0
	repeatFunction := ""
	if verbose {
		fmt.Printf("File: %s\n", filename)
	}
	for scanner.Scan() {
		line := scanner.Text()
		lineCounter++
		function, arg := getGoLineFunction(line)
		if len(function) == 0 && len(repeatFunction) > 0 {
			function = repeatFunction
			arg = line
			if nextLines > 0 {
				nextLines -= 1
			}
		} else {
			nextLines = 0
		}
		if len(function) > 0 {
			result, ok, repeat, next, err := g.eval(function, arg)
			if err == nil {
				if verbose {
					fmt.Printf("[%4.0d] %s\n", lineCounter, line)
				}
			} else {
				processFail(filename, lineCounter, line, function, verbose, err)
			}
			if !repeat && nextLines == 0 && repeatFunction == "" && next > 0 {
				nextLines = int(next)
			}
			if repeat || nextLines > 0 {
				repeatFunction = function
			} else {
				nextLines = 0
				repeatFunction = ""
			}
			if ok {
				line = processResult(result, verbose, line)
			}
		}
		output += line + "\n"
	}
	checkFail(scanner.Err(), "File "+filename+" not readable!")
	return output
}

func (g *golf) include(inc string, verbose bool) {
	if verbose {
		fmt.Printf("Including: %s", inc)
	}
	checkFile(inc)
	content, err := ioutil.ReadFile(inc)
	checkFail(err, "File could not be loaded: " + inc)
	source := string(content)
	g.interpreter.Eval(source)
}

func getGoLineFunction(line string) (function string, arg string) {
	function = ""
	prefix := funcProcessingPrefix
	functionIdx := strings.Index(line, prefix)
	if functionIdx >= 0 {
		function = line[functionIdx+len(prefix):]
		function = strings.TrimSpace(function)
		arg = line[0:functionIdx]
	}
	return
}

func processResult(result string, verbose bool, line string) string {
	processed := fmt.Sprintf("%s", result)
	if verbose {
		fmt.Printf("     ➥ %s\n", processed)
	}
	line = fmt.Sprintf("%s", processed)
	return line
}

func processFail(filename string, lineCounter int, line string, function string, verbose bool, err error) {
	if !verbose {
		printTitle()
		fmt.Printf("File: %s\n", filename)
	}
	fmt.Printf("[%4.0d] %s\n", lineCounter, line)
	println()
	fmt.Printf(" Code: %s\n", function)
	checkFail(err, fmt.Sprintf("Check the golf syntax in %s:%d", filename, lineCounter))
}
