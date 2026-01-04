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
	Completer  *readline.PrefixCompleter
	Readline   *readline.Instance
	TabCount   int
	LastPrefix string
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
	//Reset tab count when prefix changes
	if prefix != c.LastPrefix {
		c.TabCount = 0
		c.LastPrefix = prefix
	}
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
	if len(matches) == 0 {
		if c.TabCount == 1 {
			fmt.Printf("\a")
		}
		//nothing to show
		return nil
	}
	if len(matches) == 1 {
		//only 1 match -> complete and add trailing space
		c.TabCount = 0
		return []string{matches[0] + " "}
	}
	//multiple matches -> try longest common prefix(LCP)
	sort.Strings(matches)
	lcp := longestCommonPrefix(matches)
	//if lcp extends the typed prefix, return it to complete to that point
	if len(lcp) > len(prefix) {
		c.TabCount = 0
		return []string{lcp}
	}
	// lcp == prefix -> first tab rings bell, second tab prints matches
	if c.TabCount == 1 {
		fmt.Printf("\a")
		return nil
	}
	fmt.Println()
	fmt.Println(strings.Join(matches, " "))
	c.Readline.Refresh()
	c.TabCount = 0
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
func longestCommonPrefix(args []string) string {
	if len(args) == 0 {
		return ""
	}
	sort.Strings(args)
	first := args[0]
	last := args[len(args)-1]
	i := 0
	for i < len(first) && i < len(last) && first[i] == last[i] {
		i++
	}
	return first[:i]
}
