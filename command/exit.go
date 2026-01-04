package command

import "os"

type ExitCommand struct{}

func (e *ExitCommand) Execute(args []string) error {
	os.Exit(0)
	return nil
}
