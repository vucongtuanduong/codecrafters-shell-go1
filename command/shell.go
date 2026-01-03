package command

import (
	"fmt"
	"os"
	"os/exec"
)

func DefaultShellCommand(cmd string, args []string) {
	fmt.Printf("%s: command not found\n", cmd)
}
func runCommand(cmd string, args []string) {
	_, ok := searchPath(cmd)
	if !ok {
		fmt.Printf("%s: not found\n", cmd)
	}
	runCmd := exec.Command(cmd, args...)
	runCmd.Stdin = os.Stdin
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	_ = runCmd.Run()
}
