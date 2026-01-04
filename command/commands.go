package command

import (
	"os"
	"path/filepath"
)

type Command interface {
	Execute(args []string) error
}
type Registry struct {
	commands map[string]Command
}

func (r *Registry) Register(name string, cmd Command) {
	r.commands[name] = cmd
}
func (r *Registry) Get(name string) (Command, bool) {
	cmd, ok := r.commands[name]
	return cmd, ok
}
func (r *Registry) IsBuiltin(name string) bool {
	_, ok := r.commands[name]
	return ok && name != "external"
}
func NewCommandRegistry() *Registry {
	reg := &Registry{
		commands: make(map[string]Command),
	}
	// Register built-ins
	reg.Register("exit", &ExitCommand{})
	reg.Register("echo", &EchoCommand{})
	reg.Register("type", &TypeCommand{})
	reg.Register("pwd", &PwdCommand{})
	reg.Register("cd", &CdCommand{})
	// External commands handled by a single executor
	reg.Register("external", &ExternalCommand{})
	return reg
}

// Builtins Shared helpers (moved from type.go)
var Builtins = map[string]struct{}{
	"exit": {},
	"echo": {},
	"type": {},
	"pwd":  {},
	"cd":   {},
}

func searchPath(cmd string) (string, bool) {
	dirs := filepath.SplitList(os.Getenv("PATH"))
	for _, dir := range dirs {
		files, _ := os.ReadDir(dir)
		for _, file := range files {
			if file.Name() == cmd {
				fileInfo, _ := file.Info()
				if isExecutable(fileInfo.Mode()) {
					fullPath := filepath.Join(dir, file.Name())
					return fullPath, true
				}
			}
		}
	}
	return "", false
}

func isExecutable(mode os.FileMode) bool {
	return mode&0111 != 0
}
