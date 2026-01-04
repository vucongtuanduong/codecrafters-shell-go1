package command

import (
	"fmt"
	"io"
	"os"
)

func TypeCommand(args []string, stdout io.Writer) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "type: missing operand")
		return
	}
	cmdLookup := args[0]
	if IsBuiltin(cmdLookup) {
		fmt.Fprintf(stdout, "%s is a shell builtin\n", cmdLookup)
	} else if path, ok := FindInPath(cmdLookup); ok {
		fmt.Fprintf(stdout, "%s is %s\n", cmdLookup, path)
	} else {
		fmt.Fprintf(stdout, "%s: not found\n", cmdLookup)
	}
}
