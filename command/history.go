package command

import (
	"fmt"
	"io"
)

var History []string

func AppendHistory(cmd string) {
	History = append(History, cmd)
}
func HistoryCommand(args []string, stdout io.Writer) {
	for i, history := range History {
		fmt.Printf("%5d  %s\n", i+1, history)
	}
}
