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

package config

import (
	"fmt"
	"testing"

	"github.com/fatih/structs"
)

func checkStructs(test, expected interface{}) error {
	expectedMap := structs.Map(expected)
	testMap := structs.Map(test)
	for key := range expectedMap {
		val := testMap[key]
		expectedVal := expectedMap[key]
		if key == "DependsOn" {
			expectedArr := expectedVal.([]string)
			testArr := val.([]string)
			if len(expectedArr) != len(testArr) {
				return fmt.Errorf("Expected a length of %d in DependsOn, found a length of %d", len(expectedArr), len(testArr))
			}
			if len(expectedArr) == 0 && len(testArr) == 0 {
				return nil
			}
			for dep := range expectedArr {
				if testArr[dep] == "" || testArr[dep] != expectedArr[dep] {
					return fmt.Errorf("Exepcted %s for DependsOn of %d, found %s", expectedArr[dep], dep, testArr[dep])
				}
			}
		} else if key == "HealthCheck" || key == "CleanUp" || key == "Artifacts" {
			testSubMap := val.(map[string]interface{})
			expectedSubMap := expectedVal.(map[string]interface{})
			for i := range expectedSubMap {
				if testSubMap[i] != expectedSubMap[i] {
					return fmt.Errorf("Exepcted %v for %v, found %v %+v", expectedSubMap[i], i, testSubMap[i], testSubMap)
				}
			}
		} else if val == nil || val != expectedVal {
			return fmt.Errorf("Exepcted %s for %s, found %s", expectedVal, key, val)
		}
	}
	return nil
}

// TestLoad tests the loading of a config file
func TestLoad(t *testing.T) {
	expected := &Config{
		Environment: Environment{
			Exec: []string{"docker pull cpg1111/maestro"},
		},
		Project: Project{
			RepoURL:        "git@github.com:cpg1111/maestro.git",
			CloneCMD:       "git clone",
			AuthType:       "SSH",
			SSHPrivKeyPath: "~/.ssh/id_rsa",
			SSHPubKeyPath:  "~/.ssh/id_rsa.pub",
			Username:       "git",
			Password:       "",
			PromptForPWD:   false,
		},
		Services: []Service{
			Service{
				Name:             "test",
				Tag:              "0.0.1",
				TagType:          "",
				Path:             ".",
				BuildLogFilePath: "./test.log",
				BuildCMD:         []string{"bash -c 'docker build -f Dockerfile_build -t maestro_build . && docker run -v $(pwd)/dist:/opt/bin/ && docker build -t cpg1111/maestro .'"},
				TestCMD:          []string{"go test -v ./..."},
				CheckCMD:         []string{"docker ps -a"},
				CreateCMD:        []string{"docker push cpg1111/maestro"},
				UpdateCMD:        []string{"docker rm -f test && docker run -n test -d test"},
				DependsOn:        []string{},
				HealthCheck: HealthCheck{
					Type:              "",
					ExpectedCondition: "",
					Retrys:            0,
				},
			},
		},
		CleanUp: CleanUp{
			AdditionalCMDs: []string{"docker inspect maestro"},
			InDaemon:       false,
			Artifacts: []Artifact{
				Artifact{
					RuntimeFilePath: "./dist/maestro",
					SaveFilePath:    "/tmp/maestro",
				},
			},
		},
	}
	conf, loadErr := Load("../test_conf.toml", ".")
	if loadErr != nil {
		t.Error(loadErr)
	}
	projectErr := checkStructs(conf.Project, expected.Project)
	if projectErr != nil {
		t.Error(projectErr)
	}
	for i := range expected.Services {
		serviceErr := checkStructs(conf.Services[i], expected.Services[i])
		if serviceErr != nil {
			t.Error(serviceErr)
		}
	}
}
