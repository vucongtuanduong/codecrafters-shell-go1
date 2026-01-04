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

	directories := GetPathEnvDirectories()

	for _, dir := range directories {
		fullPath := filepath.Join(dir, command)
		if info, err := os.Stat(fullPath); err == nil && !info.IsDir() && info.Mode()&0100 != 0 {
			return fullPath, true
		}
	}
	return "", false
}
func GetPathEnvDirectories() []string {
	pathEnv := os.Getenv("PATH")
	return strings.Split(pathEnv, string(os.PathListSeparator))
}
func GetExternalCommandNameInPath() []string {
	res := make([]string, 1)
	directories := GetPathEnvDirectories()
	for _, dir := range directories {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, file := range files {
			if !file.IsDir() {
				res = append(res, file.Name())
			}
		}
	}
	return res
}
