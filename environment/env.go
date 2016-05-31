package environment

import (
	"log"
	"os"
	"os/exec"

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
	return s.cmd.Process.Pid, s.cmd.Run()
}

type concurrentEnvJob struct {
	envJob
	cmd *exec.Cmd
}

func (c *concurrentEnvJob) Run(pid chan int, status chan error) {
	c.cmd.Stderr = os.Stderr
	c.cmd.Stdout = os.Stdout
	err := c.cmd.Start()
	status <- err
	pid <- c.cmd.Process.Pid
	err = c.cmd.Wait()
	status <- err
}

func newJob(cmdStr string, sync bool) envJob {
	pwd, pwdErr := os.Getwd()
	if pwdErr != nil {
		panic(pwdErr)
	}
	cmd, cmdErr := util.FormatCommand(cmdStr, pwd)
	if cmdErr != nil {
		panic(cmdErr)
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
				if msg != nil {
					log.Fatal(msg)
				} else if count == len(conf.Exec)-1 {
					return nil
				} else {
					count++
				}
			case _ = <-pid:
			}
		}
	}
	return nil
}
