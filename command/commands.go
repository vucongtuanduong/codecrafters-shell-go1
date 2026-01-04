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
	var prevReader io.ReadCloser
	var builtInCmds []chan struct{}
	for i, cmdArgs := range cmds {
		if len(cmdArgs) == 0 {
			//invalid command slice
			cleanUpProcesses(processes, prevReader)
			return
		}
		name := cmdArgs[0]
		//Builtin handling
		if handler, ok := BuiltinRegistry[name]; ok {
			// If builtin is not the last command, connect it to the next via io.Pipe()
			if i < len(cmds)-1 {
				pr, pw := io.Pipe()
				done := make(chan struct{})
				//Run builtin in goroutine writing to pw
				go func(h CommandHandler, args []string, out *io.PipeWriter, done chan struct{}) {
					h(args, out)
					out.Close()
					close(done)
				}(handler, cmdArgs[1:], pw, done)
				// Parent doesn't hold writer; the reader becomes prevReader for next command.
				if prevReader != nil {
					prevReader.Close()
				}
				prevReader = pr
				builtInCmds = append(builtInCmds, done)
			} else {
				// Last builtin — run it synchronously, writing to final stdout.
				// Close any prior reader (we're not wiring it into the builtin's stdin here).
				if prevReader != nil {
					prevReader.Close()
					prevReader = nil
				}
				handler(cmdArgs[1:], stdout)
			}
			continue
		}
		//external ocmmand handling
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
		var r io.ReadCloser
		var w io.WriteCloser
		var err error
		if i < len(cmds)-1 {
			rf, wf, e := os.Pipe()
			if e != nil {
				//cleanup
				cleanUpProcesses(processes, prevReader)
				return
			}
			r = rf
			w = wf
			cmd.Stdout = w
		} else {
			cmd.Stdout = stdout
		}
		//stderr
		cmd.Stderr = stderr
		if err = cmd.Start(); err != nil {
			//cleanup
			if w != nil {
				w.Close()
			}
			if r != nil {
				r.Close()
			}
			cleanUpProcesses(processes, prevReader)
			return
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
	// Wait in reverse order so consumers exit first and producers receive EOF/SIGPIPE.
	for i := len(processes) - 1; i >= 0; i-- {
		processes[i].Wait()
	}
	// Wait for builtin goroutines to finish
	for _, d := range builtInCmds {
		<-d
	}
	//ensure any leftover reader is closed
	if prevReader != nil {
		prevReader.Close()
	}
}
func cleanUpProcesses(processes []*exec.Cmd, prevReader io.ReadCloser) {
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
