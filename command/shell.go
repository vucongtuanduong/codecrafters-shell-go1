package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

func DefaultShellCommand(cmd string, args []string) {
	runCommand(cmd, args)
}
func runCommand(cmd string, args []string) {
	path, ok := searchPath(cmd)
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
func SplitArgsLine(line string) []string {
	var args []string
	var b strings.Builder
	inSingle := false
	inDouble := false
	for _, r := range line {
		if inSingle {
			if r == '\'' {
				inSingle = false
				continue
			}
			// everything literal inside single quotes
			b.WriteRune(r)
			continue
		}
		if inDouble {
			if r == '"' {
				inDouble = false
				continue
			}
			b.WriteRune(r)
			continue
		}
		switch r {
		case '\'':
			inSingle = true
		case '"':
			inDouble = true
		default:
			if unicode.IsSpace(r) {
				if b.Len() > 0 {
					args = append(args, b.String())
					b.Reset()
				}
			} else {
				b.WriteRune(r)
			}
		}
	}

	// flush if ended while not at rune-end (covers non-rune-indexed case)
	if b.Len() > 0 {
		args = append(args, b.String())
	}

	return args
}
