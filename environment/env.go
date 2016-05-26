package environment

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/cpg1111/maestro/config"
)

type envJob interface{}

type syncEnvJob struct {
	envJob
	cmd []string
}

// Run runs the process and returns the child pid and/or error
func (s *syncEnvJob) Run() (int, error) {
	cmdPath, lookErr := exec.LookPath(s.cmd[0])
	if lookErr != nil {
		return -1, lookErr
	}
	cmd := exec.Command(cmdPath, s.cmd[1:]...)
	return cmd.Process.Pid, cmd.Run()
}

type concurrentEnvJob struct {
	envJob
	cmd []string
}

func (c *concurrentEnvJob) Run(pid chan int, status chan error) {
	cmdPath, lookErr := exec.LookPath(c.cmd[0])
	if lookErr != nil {
		status <- lookErr
	}
	cmd := exec.Command(cmdPath, c.cmd[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	status <- err
	pid <- cmd.Process.Pid
	err = cmd.Wait()
	status <- err
}

func newJob(cmdStr string, sync bool) envJob {
	if sync {
		return syncEnvJob{
			cmd: strings.Split(cmdStr, " "),
		}
	}
	return concurrentEnvJob{
		cmd: strings.Split(cmdStr, " "),
	}
}

// Load takes an environment config and loads the environment
func Load(conf config.Environment) error {
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
