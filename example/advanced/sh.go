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
