package sh

import . "strings"

func commentIf(cond bool, line string) string {
	if cond {
		line = "# " + line
	}
	return line
}

func uncommentIf(cond bool, line string) string {
		if cond {
			line = TrimLeft(line, " ")
			if HasPrefix(line, "#") {
				line = line[1:]
				line = Trim(line, " ")
			}
		}
		return line
}
