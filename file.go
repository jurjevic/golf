package main

import (
	"os"
)

func checkFile(filename string) {
	if !FileExists(filename) {
		fail("File not found: " + filename)
	}
}

func saveFile(dst string, content string) {
	out, err := os.Create(dst)
	checkFail(err, "File could not be created: "+dst)
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	_, err = out.WriteString(content)
	checkFail(err, "File could not be written: "+dst)
	err = out.Sync()
	checkFail(err, "File could not be persisted: "+dst)
	return
}

// Exists reports whether the named file or directory exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
