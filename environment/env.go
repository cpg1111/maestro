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
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/cpg1111/maestro/config"
	"github.com/cpg1111/maestro/util"
)

type envJob interface{}

type syncEnvJob struct {
	envJob
	cmd *exec.Cmd
}

// Run runs the process and returns the child pid and/or error
func (s *syncEnvJob) Run() (int, error) {
	log.Println("Running", s.cmd.Args)
	return s.cmd.Process.Pid, s.cmd.Run()
}

type concurrentEnvJob struct {
	envJob
	cmd *exec.Cmd
}

func (c *concurrentEnvJob) Run(pid chan int, status chan error) {
	log.Println("Running", c.cmd.Args)
	status <- c.cmd.Start()
	pid <- c.cmd.Process.Pid
	status <- c.cmd.Wait()
}

func newJob(cmdStr string, sync bool) envJob {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cmd, err := util.FmtCommand(cmdStr, pwd)
	if err != nil {
		panic(err)
	}
	if sync {
		return syncEnvJob{
			cmd: cmd,
		}
	}
	return concurrentEnvJob{
		cmd: cmd,
	}
}

// Load takes an environment config and loads the environment
func Load(conf *config.Environment) error {
	if len(conf.Env) > 0 {
		for h := range conf.Env {
			kv := strings.Split(conf.Env[h], ":")
			key := strings.ToUpper(kv[0])
			os.Setenv(key, kv[1])
		}
	}
	if len(conf.ExecSync) > 0 {
		for i := range conf.ExecSync {
			job := newJob(conf.ExecSync[i], true).(syncEnvJob)
			_, err := job.Run()
			if err != nil {
				return err
			}
		}
	}
	if len(conf.Exec) > 0 {
		pid := make(chan int)
		status := make(chan error)
		for j := range conf.Exec {
			job := newJob(conf.Exec[j], false).(concurrentEnvJob)
			go job.Run(pid, status)
		}
		count := 0
		for {
			select {
			case msg := <-status:
				count++
				if msg != nil {
					log.Fatal(msg)
				} else if count == len(conf.Exec)*2 {
					close(status)
					return nil
				}
			case _ = <-pid:
			}
		}
	}
	return nil
}
