package autocompleter

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/command"
)

type AutoCompleter struct {
	Completer *readline.PrefixCompleter
	Readline  *readline.Instance
	TabCount  int
}

func (c *AutoCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	newLine, length = c.Completer.Do(line, pos)
	if len(newLine) == 0 || length == 0 {
		fmt.Printf("\a")
	}
	return newLine, length

}

func (c *AutoCompleter) SetInstance(rl *readline.Instance) {
	c.Readline = rl
}
func (c *AutoCompleter) CompletePathExecutables(prefix string) []string {
	c.TabCount++
	dirs := command.GetPathEnvDirectories()
	var matches []string
	for _, dir := range dirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, file := range files {
			name := file.Name()
			fullPath := filepath.Join(dir, name)
			info, err := os.Stat(fullPath)
			if err != nil {
				continue
			}

			if strings.HasPrefix(name, prefix) && info.Mode()&0111 != 0 {
				if _, ok := command.Builtins[name]; !ok {
					matches = append(matches, name)
				}
			}
		}
	}
	if len(matches) == 1 {
		return matches
	}
	if c.TabCount == 1 {
		fmt.Printf("\a")
		return nil
	}
	if c.TabCount >= 2 && len(matches) > 0 {
		sort.Strings(matches)
		fmt.Println()
		fmt.Println(strings.Join(matches, "  "))
		c.Readline.Refresh()
		c.TabCount = 0
		return nil
	}
	return nil
}
func FinalCompleter() *AutoCompleter {
	completer := &AutoCompleter{}

	inner := readline.NewPrefixCompleter(
		readline.PcItem("echo"),
		readline.PcItem("cd"),
		readline.PcItem("exit"),
		readline.PcItem("type"),
		readline.PcItem("history"),
		readline.PcItem("pwd"),
		readline.PcItemDynamic(completer.CompletePathExecutables),
	)
	completer.Completer = inner
	return completer
}
