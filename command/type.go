package command

import (
	"fmt"
)

type TypeCommand struct{}

func (c *TypeCommand) Execute(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("type: missing operand")
	}
	cmdLookup := args[0]
	if _, ok := Builtins[cmdLookup]; ok {
		fmt.Printf("%s is a shell builtin\n", cmdLookup)
	} else if path, ok := searchPath(cmdLookup); ok {
		fmt.Printf("%s is %s\n", cmdLookup, path)
	} else {
		return fmt.Errorf("%s: not found", cmdLookup)
	}
	return nil
}
