package command

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type CommandHandler func(args []string, stdout io.Writer)

var Builtins = map[string]struct{}{
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
	_, ok := Builtins[cmd]
	return ok
}
func GetBuiltinCompletions(prefix string) []string {
	var matches []string
	for builtin := range BuiltinRegistry {
		if strings.HasPrefix(builtin, prefix) {
			matches = append(matches, builtin+" ")
		}
	}
	return matches
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
func ExecutePipeline(cmds [][]string, stdout io.Writer, stderr io.Writer) {
	if len(cmds) == 0 {
		return
	}
	var processes []*exec.Cmd
	var prevReader *os.File
	for i, cmdArgs := range cmds {
		if len(cmdArgs) == 0 {
			//invalid command slice
			return
		}
		name := cmdArgs[0]
		path, ok := FindInPath(name)
		if !ok {
			fmt.Fprintf(stdout, "%s: command not found\n", name)
			//cleanup: kill started process
			cleanUpProcesses(processes, prevReader)
			return
		}
		cmd := exec.Command(path, cmdArgs[1:]...)
		//stdin
		if prevReader != nil {
			cmd.Stdin = prevReader
		}
		//stdout
		var r *os.File
		var w *os.File
		var err error
		if i < len(cmds)-1 {
			r, w, err = os.Pipe()
			if err != nil {
				//cleanup
				cleanUpProcesses(processes, prevReader)
				return
			}
			cmd.Stdout = w
		} else {
			cmd.Stdout = stdout
		}
		//stderr
		cmd.Stderr = stderr
		if err := cmd.Start(); err != nil {
			//cleanup
			if w != nil {
				w.Close()
			}
			if r != nil {
				r.Close()
			}
			cleanUpProcesses(processes, prevReader)
		}
		// Parent should close its copy of previous read end and the writer we created.
		if prevReader != nil {
			prevReader.Close()
		}
		if w != nil {
			// close parent's copy of writer; child keeps its copy
			w.Close()
		}

		// Next command will use r as its stdin
		prevReader = r

		processes = append(processes, cmd)
	}
}
func cleanUpProcesses(processes []*exec.Cmd, prevReader *os.File) {
	for _, process := range processes {
		if process.Process != nil {
			process.Process.Kill()
			process.Wait()
		}
	}
	if prevReader != nil {
		prevReader.Close()
	}
}
