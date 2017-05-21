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
func CheckForProcess(pid int, found chan bool, errChan chan error) {
	procDir, err := os.Open(fmt.Sprintf("/proc/%d/", pid))
	if err != nil {
		errChan <- err
		found <- false
		return
	}
	stat, err := procDir.Stat()
	if err != nil {
		errChan <- err
		found <- false
		return
	}
	found <- (stat != nil)
	syscall.Kill(pid, syscall.SIGKILL)
}

// CheckForFile is self explanatory... checks for files
func CheckForFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	_, err = file.Stat()
	if err != nil {
		return err
	}
	return nil
}

// CheckFileContents reads a file and compares it to a string
func CheckFileContents(path, expect string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	var result []byte
	_, err = file.Read(result)
	if err != nil {
		return err
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
