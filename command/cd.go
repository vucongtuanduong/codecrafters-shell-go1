package command

import (
	"fmt"
	"os"
)

type CdCommand struct{}

func (c *CdCommand) Execute(args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("cd command requires only one parameter")
	}
	directory := args[0]
	if directory == "~" {
		return cdHomeDirectory()
	}
	info, err := os.Stat(directory)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s: No such file or directory", directory)
		}
		if os.IsPermission(err) {
			return fmt.Errorf("%s: permission denied", directory)
		}
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s: is not a directory", directory)
	}
	if err = os.Chdir(directory); err != nil {
		return fmt.Errorf("%s: failed to change directory", directory)
	}
	return nil
}

func cdHomeDirectory() error {
	homeDir := os.Getenv("HOME")
	if err := os.Chdir(homeDir); err != nil {
		return fmt.Errorf("failed to change directory")
	}
	return nil
}
