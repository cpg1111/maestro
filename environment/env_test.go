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

package environment

import (
	"testing"
	"time"

	"github.com/cpg1111/maestro/config"
	"github.com/cpg1111/maestro/util"
)

var conf = &config.Environment{
	ExecSync: []string{"echo '1'"},
	Exec:     []string{"docker run --rm -d nginx"},
}

// TestSyncRun tests running environment processes running synchronously
func TestSyncRun(t *testing.T) {
	job := newJob(conf.ExecSync[0], true).(syncEnvJob)
	_, runErr := job.Run()
	if runErr != nil {
		t.Error(runErr)
	}
}

// TestConcurrentRun tests running environment processes concurrently
func TestConcurrentRun(t *testing.T) {
	job := newJob(conf.Exec[0], false).(concurrentEnvJob)
	pidChan := make(chan int)
	foundChan := make(chan bool)
	errChan := make(chan error)
	go job.Run(pidChan, errChan)
	for {
		select {
		case <-time.After(3 * time.Second):
			t.Error("concurrentEnvJob timed out")
			break
		case pid := <-pidChan:
			go util.CheckForProcess(pid, foundChan, errChan)
			break
		case err := <-errChan:
			if err != nil {
				t.Error(err)
				return
			}
			break
		case found := <-foundChan:
			if found {
				return
			}
			t.Error("could not find child process")
			break
		}
	}
}
