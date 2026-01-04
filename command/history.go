package command

import (
	"fmt"
	"io"
	"strconv"
)

var History []string

func AppendHistory(cmd string) {
	History = append(History, cmd)
}
func HistoryCommand(args []string, stdout io.Writer) {
	if len(args) > 1 {
		fmt.Println("Too many arguments")
		return
	}
	n, err := strconv.Atoi(args[0])
	if err != nil {
		// handle parse error (set default, return, or log)
		n = 0
	}
	minIndex := len(History)
	if n < minIndex {
		minIndex = n
	}
	for i := 0; i < minIndex; i++ {
		fmt.Printf("%5d  %s\n", i+1, History[i])
	}
}
