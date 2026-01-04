package command

import (
	"fmt"
	"io"
	"os"
)

func CdCommand(args []string, stdout io.Writer) {
	if len(args) > 1 {
		fmt.Fprintf(os.Stderr, "cd command requires only one parameter\n")
		return
	}
	directory := args[0]
	if directory == "~" {
		directory = os.Getenv("HOME")
	}
	info, err := os.Stat(directory)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "%s: No such file or directory\n", directory)
			return
		}
		if os.IsPermission(err) {
			fmt.Fprintf(os.Stderr, "%s: permission denied\n", directory)
			return
		}
	}
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "%s: is not a directory\n", directory)
		return
	}
	if err = os.Chdir(directory); err != nil {
		fmt.Fprintf(os.Stderr, "%s: failed to change directory\n", directory)
		return
	}
}
