package environment

import (
	"syscall"
	"testing"

	"github.com/cpg1111/maestro/config"
)

var conf = &config.Environment{
	ExecSync: []string{"echo '1'"},
	Exec:     []string{"ping github.com"},
}

func checkForProcess(pid int, found chan bool, err chan error) {
	ptraceErr := syscall.PtraceAttach(pid)
	if ptraceErr != nil {
		err <- ptraceErr
		found <- false
		return
	}
	found <- true
	syscall.Kill(pid, syscall.SIGKILL)
}

func TestSyncRun(t *testing.T) {
	job := newJob(conf.ExecSync[0], true).(syncEnvJob)
	_, runErr := job.Run()
	if runErr != nil {
		t.Error(runErr)
	}
}

func TestConcurrentRun(t *testing.T) {
	job := newJob(conf.Exec[0], false).(concurrentEnvJob)
	pidChan := make(chan int)
	foundChan := make(chan bool)
	errChan := make(chan error)
	go job.Run(pidChan, errChan)
	for {
		select {
		case pid := <-pidChan:
			go checkForProcess(pid, foundChan, errChan)
		case err := <-errChan:
			if err != nil {
				t.Error(err)
			}
		case found := <-foundChan:
			if found {
				return
			}
			t.Error("Could not find child process")
		}
	}
}
