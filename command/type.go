package command

import (
	"fmt"
	"os"
	"path/filepath"
)

var Builtins = map[string]struct{}{
	"exit": {},
	"echo": {},
	"type": {},
	"pwd":  {},
}

func TypeCommandHandling(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "type: missing operand")
		return
	}
	cmdLookup := args[0]
	if _, ok := Builtins[cmdLookup]; ok {
		fmt.Printf("%s is a shell builtin\n", cmdLookup)
	} else if path, ok := searchPath(cmdLookup); ok {
		fmt.Printf("%s is %s\n", cmdLookup, path)
	} else {
		fmt.Printf("%s: not found\n", cmdLookup)
	}
}
func searchPath(cmd string) (string, bool) {
	dirs := filepath.SplitList(os.Getenv("PATH"))
	for _, dir := range dirs {
		files, _ := os.ReadDir(dir)
		for _, file := range files {
			if file.Name() == cmd {
				fileInfo, _ := file.Info()
				if isExecutable(fileInfo.Mode()) {
					fullPath := filepath.Join(dir, file.Name())
					return fullPath, true
				}
			}
		}
	}
	return "", false
}
func isExecutable(mode os.FileMode) bool {
	return mode&0111 != 0
}
