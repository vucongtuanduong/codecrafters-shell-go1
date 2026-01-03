package command

import (
	"fmt"
	"os"
)

func CdCommandHandling(cmd string, args []string) {
	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, "cd command requires only one parameter")
		return
	}
	directory := args[0]
	info, err := os.Stat(directory)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "%s: %s: No such file or directory\n", cmd, directory)
			return
		}
		if os.IsPermission(err) {
			fmt.Fprintf(os.Stderr, "%s: %s: permission denied\n", cmd, directory)
			return
		}
	}
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "%s: %s: is not a directory\n", cmd, directory)
		return
	}
	if err = os.Chdir(directory); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s: failed to change directory\n", cmd, directory)
		return
	}
}
