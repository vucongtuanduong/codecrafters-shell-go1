package command

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var History []string
var lastAppendIndex int

func AppendHistory(cmd string) {
	History = append(History, cmd)
}
func HistoryCommand(args []string, stdout io.Writer) {
	length := len(History)
	if len(args) == 0 {
		for i := 0; i < length; i++ {
			fmt.Printf("%5d  %s\n", i+1, History[i])
		}
		return
	}
	if len(args) == 2 {
		if args[0] == "-r" {
			filePath := args[1]
			readHistoryFromFile(filePath, stdout)
		} else if args[0] == "-w" {
			writeHistoryToFile(args[1], stdout)
			return
		} else if args[0] == "-a" {
			path := args[1]
			appendHistoryToFile(path, stdout)
		}
		return
	}
	n, err := strconv.Atoi(args[0])
	if err != nil {
		// handle parse error (set default, return, or log)
		n = 0
	}

	minIndex := length
	if n < minIndex {
		minIndex = n
	}
	minIndex = length - minIndex
	for i := minIndex; i < length; i++ {
		fmt.Printf("%5d  %s\n", i+1, History[i])
	}
}
func readHistoryFromFile(filePath string, stdout io.Writer) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(stdout, "failed to read history from: %s: %v\n", filePath, err)
		return
	}
	lines := strings.Split(string(fileContent), "\n")
	for _, line := range lines {
		//ignore empty lines
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		AppendHistory(line)
	}
}
func writeHistoryToFile(filePath string, stdout io.Writer) {
	flag := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	file, err := os.OpenFile(filePath, flag, 0644)
	if err != nil {
		fmt.Fprintf(stdout, "failed to open file %s: %v\n", filePath, err)
		return
	}
	defer file.Close()
	for _, line := range History {
		fmt.Fprintln(file, line)
	}
}
func appendHistoryToFile(filePath string, stdout io.Writer) {
	flag := os.O_CREATE | os.O_WRONLY | os.O_APPEND
	file, err := os.OpenFile(filePath, flag, 0644)
	if err != nil {
		fmt.Fprintf(stdout, "failed to open file %s: %v\n", filePath, err)
		return
	}
	defer file.Close()
	for i := lastAppendIndex; i < len(History); i++ {
		fmt.Fprintln(file, History[i])
	}
	lastAppendIndex = len(History)
	return
}
func markHistoryPersisted() {
	lastAppendIndex = len(History)
}
func InitHistoryFromFile() {
	path := os.Getenv("HISTFILE")
	if path != "" {
		readHistoryFromFile(path, os.Stdout)
		markHistoryPersisted()
	}
}
func WriteHistoryToFileWhenExit(stdout io.Writer) {
	path := os.Getenv("HISTFILE")
	if path != "" {
		writeHistoryToFile(path, os.Stdout)
	}
}
