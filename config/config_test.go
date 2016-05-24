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
		}
		if val == nil || val != expectedVal {
			return fmt.Errorf("Exepcted %s for %s, found %s", expectedVal, key, val)
		}
	}
	return nil
}

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
				Name:      "test",
				Tag:       "0.0.1",
				TagType:   "",
				Path:      ".",
				BuildCMD:  "docker build -t test .",
				TestCMD:   "go test -v ./...",
				CheckCMD:  "docker ps -a",
				CreateCMD: "docker run -n test -d test",
				UpdateCMD: "docker rm -f test && docker run -n test -d test",
				DependsOn: []string{},
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
