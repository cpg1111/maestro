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

package pipeline

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"syscall"

	"github.com/cpg1111/maestro/config"
	"github.com/cpg1111/maestro/util"
)

const PTAttach = "PTRACE_ATTACH"

func healthcheckCMD(cmd, path string) error {
	fmtCMD, err := util.FmtCommand(cmd, path)
	if err != nil {
		return err
	}
	fmtCMD.Stdout = os.Stdout
	fmtCMD.Stdout = os.Stderr
	return fmtCMD.Run()
}

func healthcheckHTTPGet(endpoint, response string) error {
	resp, err := http.Get(endpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Expected a status of 200 on endpoint '%s' but received '%d'", endpoint, resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if !strings.Contains((string)(body), response) {
		return errors.New("HTTP GET body does not match expected response")
	}
	return nil
}

func healthcheckICMPPing(ip string) error {
	_, err := net.Dial("ip:icmp", ip)
	if err != nil {
		return err
	}
	return nil
}

func healthcheckPtrace(pid int) error {
	return syscall.PtraceAttach(pid)
}

// HealthCheck performs a health on deployed services
func HealthCheck(conf *config.Service) interface{} {
	switch conf.HealthCheck.Type {
	case "CMD":
		return healthcheckCMD(conf.HealthCheck.CMD, conf.Path)
	case "HTTP_GET":
		return healthcheckHTTPGet(conf.HealthCheck.Address, conf.HealthCheck.ExpectedCondition)
	case "PING":
		return healthcheckICMPPing(conf.HealthCheck.Address)
	case PTAttach:
		return func(pid int) error { return healthcheckPtrace(pid) }
	default:
		return nil
	}
}
