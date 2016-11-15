/*
Copyright 2016 Christian Grabowski All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"fmt"
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

// CheckFileContents reads a file and compares it to a string
func CheckFileContents(path, expect string) error {
	file, openErr := os.Open(path)
	if openErr != nil {
		return openErr
	}
	var result []byte
	_, readErr := file.Read(result)
	if readErr != nil {
		return readErr
	}
	if (string)(result) != expect {
		return fmt.Errorf(
			"expected: %s, in file: %s, found: %s",
			expect,
			path,
			(string)(result),
		)
	}
	os.Remove(path)
	return nil
}
