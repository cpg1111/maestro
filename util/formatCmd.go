package util

import (
	"os"
	"os/exec"
	"os/user"
	"strings"
)

// FormatCommand formats a string to an exec'able command
func FormatCommand(strCMD, path string) (*exec.Cmd, error) {
	var cmdArr []string
	if strings.Contains(strCMD[0:8], "bash -c") {
		cmdArr = []string{strCMD[0:4], strCMD[5:7], strCMD[8:]}
	} else {
		cmdArr = strings.Split(strCMD, " ")
	}
	cmdPath, lookupErr := exec.LookPath(cmdArr[0])
	if lookupErr != nil {
		return &exec.Cmd{}, lookupErr
	}
	cmd := exec.Command(cmdPath)
	cmdLen := len(cmdArr)
	for i := 1; i < cmdLen; i++ {
		if strings.Contains(cmdArr[i], ".") {
			cmdArr[i] = strings.Replace(cmdArr[i], ".", path, 1)
		}
		if strings.Contains(cmdArr[i], "~") {
			currUser, userErr := user.Current()
			if userErr != nil {
				return &exec.Cmd{}, userErr
			}
			cmdArr[i] = strings.Replace(cmdArr[i], "~", currUser.HomeDir, 1)
		}
		cmd.Args = append(cmd.Args, cmdArr[i])
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd, nil
}
