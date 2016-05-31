package util

import (
	"os"
	"syscall"
)

// CheckForProcess checks for a running process using the ptrace syscall
func CheckForProcess(pid int, found chan bool, err chan error) {
	ptraceErr := syscall.PtraceAttach(pid)
	if ptraceErr != nil {
		err <- ptraceErr
		found <- false
		return
	}
	found <- true
	syscall.Kill(pid, syscall.SIGKILL)
}

// CheckForFile is self explanatory... checks for files
func CheckForFile(path string) error {
	file, openErr := os.Open(path)
	if openErr != nil {
		return openErr
	}
	_, statErr := file.Stat()
	if statErr != nil {
		return statErr
	}
	return nil
}
