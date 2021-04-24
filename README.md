# golf

The *Go Language Facilitator* tool is used to process in-file code while copying a source to a destination file. It is basically helpful when different files are required for other build destinations, but templating is not provided out-of-the-box. This gives you the possibility to have different content in files for e.g. local and server builds.

![](https://github.com/jurjevic/golf/blob/main/example/simple/tty.gif)

##  Install
```sh
#  Version: 0.8 (work-in-progress)
# Required: Go 1.16+

go install github.com/jurjevic/golf@latest
```

## Usage
```
golf source_file target_file [-i=<file>] [-v] [-- {initialzation}]
```
- `source_file` is the input file which will be processed.
- `target_file` will be created or overwritten with the processed input.
- `-i` is optional and includes the file for evaluation before the processing is started.
- `-v` is optional and produces verbose outputs.
- `--` is optional and separates the initialization code, which is interpreted before the processing starts.

The source and destination files may be the same, so the source file will be overwritten.

```sh
# if you have cloned the golf repository, you can test it with ...
golf README.md README.tmp -v
```

### Details
**golf** supports the Go language for processing text lines. 
Functions are introduced by ```###``` with a whitespace before and after this token and are usually placed inside a comment in the source file. Depending on the source file type, comments may start if different charters and sequences. **golf** is not aware of any file types, so comments are agnostic to it and the developer has to take care of it.

The `"strings"` and `"fmt"` packages are dot imported by default. That means you don't have to use the package prefix to access their global functions. Other Go packages can be used to, but you would have to import them manually.

If the last statement process has a `string` type result, its value is used to replace the line. This can be easily avoided by using a different type at the end of the statement.

#### Built-in variables
```go
// line represents the current line containing the string on the left side of '###'.
var line string

// token is the sliced line of strings splitted by whitespaces.
var token []string

// repeat will execute the following line statements also for the next lines,
// until a new statement is defined.
var repeat bool

// next will execute the statement for the current and all next defined amount of lines. It will also stop, when a new statement is defined.
var next uint

// arg are optional key values pairs, which can be assign from
// outside at initialization time.
var arg map[string]interface{}
```
#### Built-in functions
```go
// isSet returns true if the key is found in the arg map
func isSet(key string) bool
```

### Include files
Go files can be included before the processing is executed. With the given example, the statement `# ### commentIf(true, line)` can be used.
```go
package sh

import "strings"

func commentIf(cond bool, line string) string {
	if cond {
		line = "# " + line
	}
	return line
}

func uncommentIf(cond bool, line string) string {
	if cond {
		line = strings.TrimLeft(line, " ")
		if strings.HasPrefix(line, "#") {
			line = line[1:]
			line = strings.TrimLeft(line, " ")
		}
	}
	return line
}
```

## Examples

### Simple line examples
```go
world // ### "hello " + line
```
```go
hello world //
```

### Simple token examples
```go
hello my friend // ### token[0] + " world"
```
```go
hello world
```

### Replace token examples
```go
drive the rabbit home // ### token[2] = "cow"; Join(token, " ")
```
```go
drive the cow home
```

### Advanced token examples
```go
hello my friend // ### "BE " + ToUpper(Join(token[1:], " "))
```
```go
BE MY FRIEND //
```

### Advanced processing examples
```go
hello my friend // ### otherToken := token
let it be // ### anotherToken := token
goodbye dad  // ### token[0] + " and " + Join(anotherToken[:3], " ") + " " + Join(otherToken[1:], " ")
```
```go
hello my friend // ### otherToken := token
let it be // ### anotherToken := token
goodbye and let it be my friend //
```

### Ignore processing examples
```go
what is wrong // ### s := Sprintf("nothing %s", Join(token[1:], " ")); false
what was wrong // ### s
```
```go
what is wrong // ### s := Sprintf("nothing %s", Join(token[1:], " ")); false
nothing is wrong //
```

### Repeat examples
```go
Hello number one ### repeat = true; "# " + line
Hello number two
Hello number three
Hello number four
Hello number five ### line + " and five"
Hello number six
```
```go
# Hello number one
# Hello number two
# Hello number three
# Hello number four
Hello number five and five
Hello number six
```

### Next examples
```go
// ### next = 2; "// " + line
red
blue
green
yellow
```
```go
// //
// red
// blue
green
yellow
```

### Custom processing examples
```go
// ### func modExample(line string) string { return "_" + line + "_" }
Spaceship // ### modExample(token[0])
```
```go
// ### func modExample(line string) string { return "_" + line + "_" }
_Spaceship_
```

### Initialize examples
```go
// the output will change depending whether you use an initialization:
// golf README.md README.tmp -v -- 'arg["second"] = true'
// ### if arg["second"] == true { line = "Loud out" }; line
```
```go
// the output will change depending whether you use an initialization:
// golf README.md README.tmp -v -- 'arg["second"] = true'
//
```

**Bash examples:**
```bash
echo "hello my friend $1" # ### Join(token[:4], " ") + " $2\""
```
```bash
echo "hello my friend $2" #
```

### Bad examples

**Undefined processing token:**
```bash
echo "hello" #### Join(token[1:], " ")
echo "world" ## Join(token[1:], " ")
```
```bash
echo "hello" #### Join(token[1:], " ")
echo "world" ## Join(token[1:], " ")
```
The processing is ignored, because the token `###` is missing (including before and after a whitespace).

### Testable simple example
* [red.sh](https://github.com/jurjevic/golf/blob/main/example/simple/red.sh) Input file with simple example.
```sh
# run without processing
chmod +x red.sh && ./red.sh

# perform processing
golf red.sh blue.sh

# run processed example
chmod +x blue.sh && ./blue.sh
```

### Testable advanced example
* [gen.sh](https://github.com/jurjevic/golf/blob/main/example/advanced/gen.sh) Call example for `golf` command.
* [red.sh](https://github.com/jurjevic/golf/blob/main/example/advanced/red.sh) Input file with advanced examples.
* [sh.go](https://github.com/jurjevic/golf/blob/main/example/advanced/sh.go) Helper functions for shell processing.
```sh
# run without processing
chmod +x red.sh && ./red.sh

# perform processing
./gen.sh

# run processed example
chmod +x blue.sh && ./blue.sh

# perform another processing
./gen.sh green

# run green processed example
chmod +x green.sh && ./green.sh
```

## Credits
*golf* is uses [github.com/traefik/yaegi](https://github.com/traefik/yaegi) as the Go language interpreter.

## License
[MIT](https://github.com/jurjevic/golf/blob/main/LICENSE)

