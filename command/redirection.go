package command

import (
	"io"
	"os"
)

func ParseAndSetupRedirection(comarr []string) ([]string, string, bool, string, bool) {
	var stdoutFilePath string
	var stderrFilePath string
	isStdoutAppend := false
	isStderrAppend := false
	args := comarr
	for i, arg := range comarr {
		if (arg == ">" || arg == "1>") && i+1 < len(comarr) {
			stdoutFilePath = comarr[i+1]
			args = comarr[:i]
			break
		} else if arg == "2>" && i+1 < len(comarr) {
			stderrFilePath = comarr[i+1]
			args = comarr[:i]
			break
		} else if (arg == ">>" || arg == "1>>") && i+1 < len(comarr) {
			stdoutFilePath = comarr[i+1]
			args = comarr[:i]
			isStdoutAppend = true
			break
		} else if arg == "2>>" && i+1 < len(comarr) {
			stderrFilePath = comarr[i+1]
			args = comarr[:i]
			isStderrAppend = true
			break
		}
	}
	return args, stdoutFilePath, isStdoutAppend, stderrFilePath, isStderrAppend
}
func OpenRedirectFile(path string, append bool, defaultWriter io.Writer) (io.Writer, func(), error) {
	if path == "" {
		return defaultWriter, nil, nil
	}
	flag := os.O_WRONLY | os.O_CREATE
	if append {
		flag |= os.O_APPEND
	} else {
		flag |= os.O_TRUNC
	}
	file, err := os.OpenFile(path, flag, 0644)
	if err != nil {
		return nil, nil, err
	}
	return file, func() { file.Close() }, nil
}
