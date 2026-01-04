package main

import (
	"bufio"
	"fmt"
	"io"
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
		args, redirectPath := command.ParseAndSetupRedirection(comarr)

		// Set stdout writer.
		var stdout io.Writer = os.Stdout
		if redirectPath != "" {
			file, err := os.OpenFile(redirectPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "err: %v\n", err)
				continue
			}
			defer file.Close()
			stdout = file
		}

		command.ExecuteCommand(args, stdout)
	}
}
