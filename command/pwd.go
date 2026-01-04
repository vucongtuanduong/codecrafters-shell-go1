package command

import (
	"fmt"
	"io"
	"os"
)

func PwdCommand(args []string, stdout io.Writer) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "pwd: %v\n", err)
		return
	}
	fmt.Fprintln(stdout, dir)
}
