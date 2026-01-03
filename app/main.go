package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	for {
		fmt.Print("$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}
		trimmedCommand := command[:len(command)-1]
		//fmt.Println("command:" + command + ".")
		params := strings.Split(trimmedCommand, " ")
		switch params[0] {
		case "exit":
			os.Exit(0)
		case "echo":
			fmt.Println(strings.Join(params[1:], " "))
		default:
			fmt.Printf("%s: command not found\n", trimmedCommand)
		}

	}

}
