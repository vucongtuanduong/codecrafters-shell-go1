package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/command"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		input = strings.TrimSpace(input)
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
