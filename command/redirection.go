package command

import (
	"fmt"
)

func ParseAndSetupRedirection(comarr []string) ([]string, string) {
	var redirectPath string
	args := comarr
	for i, arg := range comarr {
		if (arg == ">" || arg == "1>") && i+1 < len(comarr) {
			redirectPath = comarr[i+1]
			args = comarr[:i]
			fmt.Printf("args = comarr[:i]: %s\n", args)
			break
		}
	}
	return args, redirectPath
}
