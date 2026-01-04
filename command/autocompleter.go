package command

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/chzyer/readline"
)

type AutoCompleter struct {
	Completer  readline.AutoCompleter
	Readline   *readline.Instance
	TabCount   int
	LastPrefix string
}

func (c *AutoCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	newLine, length = c.Completer.Do(line, pos)
	if len(newLine) > 0 && length > 0 {
		//Reset state when inner provided completions
		c.TabCount = 0
		c.LastPrefix = ""
		return newLine, length
	}
	// Determine current token (last "word" before cursor)
	prefixStr := string(line[:pos])
	lastSpace := strings.LastIndex(prefixStr, " ")
	var token string
	if lastSpace == -1 {
		token = prefixStr
	} else {
		token = prefixStr[lastSpace+1:]
	}
	//Reset counter if prefix changed
	if token != c.LastPrefix {
		c.TabCount = 0
		c.LastPrefix = token
	}
	// Find executable matches for the token
	matches := c.CompletePathExecutables(token)
	switch len(matches) {
	case 0:
		//No match -> bell, no change
		fmt.Print("\a")
		return nil, 0
	case 1:
		//single match -> complete with trailing space
		comp := []rune(matches[0] + " ")
		c.TabCount = 0
		c.LastPrefix = ""
		return [][]rune{comp}, pos
	default:
		// Multiple matches handling
		c.TabCount++
		if c.TabCount == 1 {
			fmt.Print("\a")
			return nil, 0
		}
		//second tab -> print matches
		c.Readline.Write([]byte("\n"))
		c.Readline.Write([]byte(strings.Join(matches, "  ")))
		c.Readline.Write([]byte("\n"))
		c.Readline.Refresh()
		c.TabCount = 0
		return nil, 0
	}

}

func (c *AutoCompleter) SetInstance(rl *readline.Instance) {
	c.Readline = rl
}
func (c *AutoCompleter) CompletePathExecutables(prefix string) []string {
	if prefix == "" {
		return nil
	}
	dirs := getPathEnvDirectories()
	set := map[string]struct{}{}
	for _, dir := range dirs {
		files, _ := os.ReadDir(dir)
		for _, file := range files {
			name := file.Name()
			if !strings.HasPrefix(name, prefix) {
				continue
			}
			if _, ok := Builtins[name]; ok {
				continue
			}
			info, err := file.Info()
			if err != nil {
				continue
			}
			if info.Mode()&0111 != 0 {
				set[name] = struct{}{}
			}
		}
	}
	if len(set) == 0 {
		return nil
	}
	res := make([]string, 0, len(set))
	for i := range set {
		res = append(res, i)
	}
	sort.Strings(res)
	return res
}
func getBinariesCompletion() []readline.PrefixCompleterInterface {
	var items []readline.PrefixCompleterInterface
	commands := []string{"cd", "echo", "pwd", "exit", "type", "history"}
	for i := range commands {
		items = append(items, readline.PcItem(commands[i]))
	}
	//// Add System Binaries
	//systemBinariesNamePath := GetExternalCommandNameInPath()
	//for i := range systemBinariesNamePath {
	//	items = append(items, readline.PcItem(systemBinariesNamePath[i]))
	//}
	return items
}
func FinalCompleter() *AutoCompleter {

	completer := readline.NewPrefixCompleter(
		getBinariesCompletion()...,
	)
	finalCompleter := &AutoCompleter{Completer: completer}
	return finalCompleter
}
