package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/codecrafters-io/shell-starter-go/command"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	reader := bufio.NewReader(os.Stdin)
	reg := command.NewCommandRegistry()
	repl("$ ", reader, reg)

}
func repl(prompt string, reader *bufio.Reader, reg *command.Registry) {
	for {
		fmt.Print(prompt)
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println()
				return
			}
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}
		args := command.ParseLine(line)
		if len(args) == 0 {
			continue
		}
		cmdName := args[0]
		cmdArgs := args[1:]

		if cmd, ok := reg.Get(cmdName); ok {
			if err := cmd.Execute(cmdArgs); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		} else {
			// Fallback to external command
			cmd, err := reg.Get("external")
			if !err {
				fmt.Fprintln(os.Stderr, err)
			}
			cmd.Execute(args)
		}
	}
}
