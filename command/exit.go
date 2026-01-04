package command

import (
	"io"
	"os"
)

func ExitCommand(args []string, stdout io.Writer) {
	os.Exit(0)
}
