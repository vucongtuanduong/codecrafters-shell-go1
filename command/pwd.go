package command

import (
	"fmt"
	"os"
)

type PwdCommand struct{}

func (c *PwdCommand) Execute(args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error printing working directory")
	}
	fmt.Println(dir)
	return nil
}
