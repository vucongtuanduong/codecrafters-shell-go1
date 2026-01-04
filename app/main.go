package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/command"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	completer := readline.NewPrefixCompleter(readline.PcItem("echo"),
		readline.PcItem("exit"))
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "$ ",
		AutoComplete: completer,
	})
	if err != nil {
		return
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		input := strings.TrimSpace(line)
		if input == "" {
			continue
		}

		comarr := command.ParseInput(input)
		// Parse redirection.
		args, stdoutFilePath, isStdoutAppend, stderrFilePath, isStderrAppend := command.ParseAndSetupRedirection(comarr)

		// Set writers.
		stdout, closeStdout, err := command.OpenRedirectFile(stdoutFilePath, isStdoutAppend, os.Stdout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "err: %v\n", err)
			continue
		}
		stderr, closeStderr, err := command.OpenRedirectFile(stderrFilePath, isStderrAppend, os.Stderr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "err: %v\n", err)
			if closeStdout != nil {
				closeStdout()
			}
			continue
		}

		// Defer closes.
		if closeStdout != nil {
			defer closeStdout()
		}
		if closeStderr != nil {
			defer closeStderr()
		}

		command.ExecuteCommand(args, stdout, stderr)
	}
}
