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
	"io"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/aws/aws-sdk-go/aws"
	awscreds "github.com/aws/aws-sdk-go/aws/credentials"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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

type remoteConfig struct {
	Storage string
	Bucket  string
	Object  string
}

func parseRemote(path string) *remoteConfig {
	storageIdx := strings.Index(path, "://")
	pathSlice := strings.Split(path[storageIdx+1:], "/")
	obj := pathSlice[1]
	if len(pathSlice) > 2 {
		for i := 2; i < len(pathSlice); i++ {
			obj = fmt.Sprintf("%s/%s", obj, pathSlice[i])
		}
	}
	return &remoteConfig{
		Storage: path[0:storageIdx],
		Bucket:  pathSlice[0],
		Object:  obj,
	}
}

func decode(r io.Reader) (*Config, error) {
	var conf Config
	if _, pErr := toml.DecodeReader(r, &conf); pErr != nil {
		return &conf, pErr
	}
	return &conf, nil
}

func loadLocal(path string) (*Config, error) {
	conf, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	return decode(conf)
}

func loadS3(path string) (*Config, error) {
	remote := parseRemote(path)
	creds := awscreds.NewEnvCredentials()
	_, err := creds.Get()
	if err != nil {
		return nil, err
	}
	config := &aws.Config{
		Region:           aws.String(os.Getenv("AWS_S3_REGION")),
		Endpoint:         aws.String("s3.amazonaws.com"),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      creds,
		LogLevel:         aws.LogLevel(aws.LogLevelType(0)),
	}
	session := awssession.New(config)
	s3Client := s3.New(session)
	query := &s3.GetObjectInput{
		Bucket: aws.String(remote.Bucket),
		Key:    aws.String(remote.Object),
	}
	resp, err := s3Client.GetObject(query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return decode(resp.Body)
}

// Load reads the config and returns a Config struct
func Load(path, clonePath string) (Config, error) {
	var (
		conf *Config
		rErr error
	)
	if strings.Contains(path, "s3://") {
		conf, rErr = loadS3(path)
	} else {
		conf, rErr = loadLocal(path)
	}
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
