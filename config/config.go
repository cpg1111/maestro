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
	"io"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

// HealthCheck is a struct to check for a service's 'upness'
type HealthCheck struct {
	Type              string // Either cmd, http_get, icmp_ping or ptrace_attach
	CMD               string
	Address           string
	ExpectedCondition string
	Retrys            int
}

// Service is a struct of the service to build
type Service struct {
	Name             string
	Tag              string
	TagType          string
	Path             string
	BuildLogFilePath string
	BuildCMD         []string
	TestCMD          []string
	CheckCMD         []string
	CreateCMD        []string
	UpdateCMD        []string
	HealthCheck      HealthCheck `toml:"[Services.HealthCheck]"`
	DependsOn        []string
}

// Project is the struct of the project to build
type Project struct {
	RepoURL         string
	CloneCMD        string
	AuthType        string
	SSHPrivKeyPath  string
	SSHPubKeyPath   string
	Username        string
	Password        string
	PromptForPWD    bool
	MaestrodHostEnv string
	MaestrodPortEnv string
	MaestrodHost    string
	MaestrodPort    int
}

// Environment is the config for the environment
type Environment struct {
	Env      []string
	ExecSync []string
	Exec     []string
}

// Artifact is struct for what artifacts to save post-pipeline
type Artifact struct {
	RuntimeFilePath string
	SaveFilePath    string
}

// CleanUp is a struct for the post-pipeline actions to clean up the build and save artifacts
type CleanUp struct {
	AdditionalCMDs []string
	InDaemon       bool
	Artifacts      []Artifact
}

// Config is a struct to parse the TOML config into
type Config struct {
	Environment Environment
	Project     Project
	Services    []Service
	CleanUp     CleanUp
}

func decode(r io.Reader) (*Config, error) {
	var conf Config
	if _, pErr := toml.DecodeReader(r, &conf); pErr != nil {
		return &conf, pErr
	}
	return &conf, nil
}

func loadLocal(path string) (*Config, error) {
	conf, readErr := os.OpenFile(path, os.O_RDONLY, 0644)
	if readErr != nil {
		return nil, readErr
	}
	return decode(conf)
}

// Load reads the config and returns a Config struct
func Load(path, clonePath string) (Config, error) {
	conf, rErr := loadLocal(path)
	if rErr != nil {
		return *conf, rErr
	}
	for i := range conf.Services {
		if strings.Contains(conf.Services[i].Path, ".") {
			conf.Services[i].Path = strings.Replace(conf.Services[i].Path, ".", clonePath, 1)
		}
	}
	return *conf, nil
}
