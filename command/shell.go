package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func DefaultShellCommand(cmd string, args []string) {
	fmt.Printf("%s: command not found\n", cmd)
}
func runCommand(cmd string, args []string) {
	str := strings.Join(args, " ")
	lines := cmd + " " + str
	path, ok := searchPath(lines)
	if !ok {
		fmt.Printf("%s: not found\n", cmd)
		return
	}
	runCmd := exec.Command(cmd, args...)
	runCmd.Path = path
	runCmd.Stdin = os.Stdin
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	_ = runCmd.Run()
}
