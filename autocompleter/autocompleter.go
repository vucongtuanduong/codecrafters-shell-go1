package autocompleter

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/command"
)

type AutoCompleter struct {
	Completer  *readline.PrefixCompleter
	Readline   *readline.Instance
	TabCount   int
	LastPrefix string
	LastCall   time.Time
}

func (c *AutoCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	newLine, length = c.Completer.Do(line, pos)
	if len(newLine) == 0 || length == 0 {
		upToCursor := string(line[:pos])
		lastSpace := strings.LastIndex(upToCursor, " ")
		var current string
		if lastSpace == -1 {
			current = upToCursor
		} else {
			current = upToCursor[lastSpace+1:]
		}

		if current != "" {
			suggestions := c.CompletePathExecutables(current)
			if suggestions != nil && len(suggestions) > 0 {
				res := make([][]rune, len(suggestions))
				for i, s := range suggestions {
					res[i] = []rune(s)
				}
				// length is the number of characters to replace (the token we completed)
				return res, len(current)
			}
		}
		fmt.Printf("\a")
	}
	return newLine, length
}

func (c *AutoCompleter) SetInstance(rl *readline.Instance) {
	c.Readline = rl
}

func (c *AutoCompleter) CompletePathExecutables(prefix string) []string {
	now := time.Now()

	// Reset tab count when prefix changes
	if prefix != c.LastPrefix {
		c.TabCount = 0
		c.LastPrefix = prefix
		// reset LastCall to zero time so first press always increments
		c.LastCall = time.Time{}
	}

	// Debounce: if this call is essentially the same call (same prefix, very quick),
	// don't increment TabCount (avoid double-increment from multiple caller paths).
	isDuplicate := !c.LastCall.IsZero() && prefix == c.LastPrefix && now.Sub(c.LastCall) < 100*time.Millisecond
	if !isDuplicate {
		c.TabCount++
	}
	c.LastCall = now

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
		return nil
	}

	// Only one match -> return the suffix + trailing space
	if len(matches) == 1 {
		c.TabCount = 0
		c.LastCall = time.Time{}
		suffix := matches[0][len(prefix):]
		return []string{suffix + " "}
	}

	// multiple matches -> try longest common prefix (LCP)
	sort.Strings(matches)
	lcp := longestCommonPrefix(matches)

	// If lcp extends the typed prefix, return the *suffix* (no space)
	if len(lcp) > len(prefix) {
		c.TabCount = 0
		c.LastCall = time.Time{}
		return []string{lcp[len(prefix):]}
	}

	// lcp == prefix -> first tab rings bell, second tab prints matches
	if c.TabCount == 1 {
		fmt.Printf("\a")
		return nil
	}

	// second tab: list matches
	fmt.Println()
	fmt.Println(strings.Join(matches, "  "))
	c.Readline.Refresh()
	c.TabCount = 0
	c.LastCall = time.Time{}
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
