package command

import (
	"fmt"
	"strings"
)

type EchoCommand struct{}

func (c *EchoCommand) Execute(args []string) error {
	fmt.Println(strings.Join(args, " "))
	return nil
}
