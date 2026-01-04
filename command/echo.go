package command

import (
	"fmt"
	"io"
	"strings"
)

func EchoCommand(args []string, stdout io.Writer) {
	fmt.Fprintln(stdout, strings.Join(args, " "))
}
