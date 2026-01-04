package command

import (
	"fmt"
	"os"
	"os/exec"
)

type ExternalCommand struct{}

func (c *ExternalCommand) Execute(args []string) error {
	if len(args) == 0 {
		return nil
	}
	cmd := args[0]
	path, ok := searchPath(cmd)
	if !ok {
		return fmt.Errorf("%s: not found", cmd)
	}
	runCmd := exec.Command(cmd, args[1:]...)
	runCmd.Path = path
	runCmd.Stdin = os.Stdin
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	return runCmd.Run()
}
