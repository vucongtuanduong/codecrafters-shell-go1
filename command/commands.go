package command

import (
	"fmt"
	"io"
)

type CommandHandler func(args []string, stdout io.Writer)

var builtins = map[string]struct{}{
	"exit": {},
	"echo": {},
	"pwd":  {},
	"cd":   {},
	"type": {},
}
var BuiltinRegistry = map[string]CommandHandler{
	"exit": ExitCommand,
	"echo": EchoCommand,
	"pwd":  PwdCommand,
	"cd":   CdCommand,
	"type": TypeCommand,
}

func IsBuiltin(cmd string) bool {
	_, ok := builtins[cmd]
	return ok
}
func ExecuteCommand(args []string, stdout io.Writer, stderr io.Writer) {
	if len(args) == 0 {
		return
	}
	name := args[0]
	cmdArgs := args[1:]

	if handler, ok := BuiltinRegistry[name]; ok {
		handler(cmdArgs, stdout)
		return
	}

	// External.
	if !ExternalCommand(args, stdout, stderr) {
		fmt.Fprintf(stdout, "%s: command not found\n", name)
	}
}
