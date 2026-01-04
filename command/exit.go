package command

import (
	"io"
	"os"
)

func ExitCommand(args []string, stdout io.Writer) {
	WriteHistoryToFileWhenExit(stdout)
	os.Exit(0)
}
