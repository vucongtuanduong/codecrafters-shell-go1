package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/autocompleter"
	"github.com/codecrafters-io/shell-starter-go/command"
	"github.com/codecrafters-io/shell-starter-go/parser"
	"github.com/codecrafters-io/shell-starter-go/redirection"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

const PROMPT = "$ "

func main() {
	completer := autocompleter.FinalCompleter()
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       PROMPT,
		AutoComplete: completer,
	})

	if err != nil {
		panic(err)
	}
	completer.SetInstance(rl)
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

		comarr := parser.ParseInput(input)
		// Parse redirection.
		args, stdoutFilePath, isStdoutAppend, stderrFilePath, isStderrAppend := redirection.ParseAndSetupRedirection(comarr)

		// Set writers.
		stdout, closeStdout, err := redirection.OpenRedirectFile(stdoutFilePath, isStdoutAppend, os.Stdout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "err: %v\n", err)
			continue
		}
		stderr, closeStderr, err := redirection.OpenRedirectFile(stderrFilePath, isStderrAppend, os.Stderr)
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
		//Split args into pipline segments by '|'
		var segments [][]string
		start := 0
		for i, arg := range args {
			if arg == "|" {
				segments = append(segments, args[start:i])
				start = i + 1
			}
		}
		segments = append(segments, args[start:])
		//validate segments(no empty commands)
		invalid := false
		for _, seg := range segments {
			if len(seg) == 0 {
				invalid = true
				break
			}
		}
		if invalid {
			fmt.Fprintln(os.Stderr, "invalid pipeline")
			continue
		}
		if len(segments) == 1 {
			command.ExecuteCommand(segments[0], stdout, stderr)
		} else {
			command.ExecutePipeline(segments, stdout, stderr)
		}

	}
}
