package command

import (
	"fmt"
	"io"
	"os"
)

func CdCommand(args []string, stdout io.Writer) {
	if len(args) > 1 {
		fmt.Errorf("cd command requires only one parameter")
	}
	directory := args[0]
	if directory == "~" {
		cdHomeDirectory()
	}
	info, err := os.Stat(directory)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Errorf("%s: No such file or directory", directory)
		}
		if os.IsPermission(err) {
			fmt.Errorf("%s: permission denied", directory)
		}
	}
	if !info.IsDir() {
		fmt.Errorf("%s: is not a directory", directory)
	}
	if err = os.Chdir(directory); err != nil {
		fmt.Errorf("%s: failed to change directory", directory)
	}
}

func cdHomeDirectory() error {
	homeDir := os.Getenv("HOME")
	if err := os.Chdir(homeDir); err != nil {
		return fmt.Errorf("failed to change directory")
	}
	return nil
}
