package command

import (
	"fmt"
	"os"
)

func PwdCommandHandling() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error printing working directory")
		return
	}
	fmt.Println(dir)
}
