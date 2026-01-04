package command

import (
	"fmt"
	"github.com/chzyer/readline"
	"strings"
	"unicode"
)

type BellCompleter struct {
	Completer readline.AutoCompleter
}

func (b *BellCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	newLine, length = b.Completer.Do(line, pos)
	if len(newLine) == 0 {
		fmt.Println("\a")
	}
	return newLine, length
}
func getBinariesCompletion() []readline.PrefixCompleterInterface {
	var items []readline.PrefixCompleterInterface
	commands := []string{"cd", "echo", "pwd", "exit", "type", "history"}
	for i := range commands {
		items = append(items, readline.PcItem(commands[i]))
	}
	// Add System Binaries
	systemBinariesNamePath := GetExternalCommandNameInPath()
	for i := range systemBinariesNamePath {
		items = append(items, readline.PcItem(systemBinariesNamePath[i]))
	}
	return items
}
func FinalCompleter() *BellCompleter {

	completer := readline.NewPrefixCompleter(
		getBinariesCompletion()...,
	)
	finalCompleter := &BellCompleter{Completer: completer}
	return finalCompleter
}
func ParseInput(line string) []string {
	var args []string
	var b strings.Builder
	inSingle := false
	inDouble := false
	escaped := false
	if len(line) > 0 && line[len(line)-1] == '\n' {
		line = line[:len(line)-1]
	}
	for i, r := range line {
		if escaped {
			// Outside quotes: backslash escapes any character (backslash removed)
			// Inside double quotes: handled by setting escaped only when allowed below,
			// so here escaped means "append rune literally"
			b.WriteRune(r)
			escaped = false
			continue
		}
		if inSingle {
			// Single quotes: everything literal until next single quote
			if r == '\'' {
				inSingle = false
				continue
			}
			// everything literal inside single quotes
			b.WriteRune(r)
			continue
		}
		if inDouble {
			// In double quotes: backslash only escapes ", \, $, `, and newline.
			if r == '"' {
				inDouble = false
				continue
			}
			if r == '\\' {
				// Lookahead: if next rune is one of the special ones, consume it literally.
				// Otherwise backslash is literal.
				if i+1 < len(line) {
					nextRune := rune(line[i+1])
					switch nextRune {
					case '"', '\\', '$', '`', '\n':
						//set escaped = true so next iteration appends next literally and skips special handling
						escaped = true
						continue
					default:
						//backlash is literal
						b.WriteRune(r)
						continue
					}
				} else {
					// trailing backslash in double quotes — treat as literal backslash
					b.WriteRune(r)
					continue
				}
			}
			b.WriteRune(r)
			continue
		}
		switch r {
		case '\\':
			escaped = true
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
	// If a backslash was the last character outside quotes, it escaped nothing — treat it as literal
	if escaped {
		b.WriteByte('\\')
	}
	// flush if ended while not at rune-end (covers non-rune-indexed case)
	if b.Len() > 0 {
		args = append(args, b.String())
	}

	return args
}
