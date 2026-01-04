package command

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ExternalCommand(args []string, stdout io.Writer, stderr io.Writer) bool {
	name := args[0]
	cmdArgs := args[1:]
	if path, ok := FindInPath(name); ok {
		cmd := exec.Command(path, cmdArgs...)
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		cmd.Args[0] = name
		cmd.Run()
		return true
	}
	return false
}
func FindInPath(command string) (string, bool) {
	pathEnv := os.Getenv("PATH")
	directories := strings.Split(pathEnv, ":")

	for _, dir := range directories {
		fullPath := filepath.Join(dir, command)
		if info, err := os.Stat(fullPath); err == nil && !info.IsDir() && info.Mode()&0100 != 0 {
			return fullPath, true
		}
	}
	return "", false
}
