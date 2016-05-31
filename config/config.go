package config

import (
	"io/ioutil"
	"strings"

	"github.com/BurntSushi/toml"
)

// Service is a struct of the service to build
type Service struct {
	Name      string
	Tag       string
	TagType   string
	Path      string
	BuildCMD  string
	TestCMD   string
	CheckCMD  string
	CreateCMD string
	UpdateCMD string
	DependsOn []string
}

// Project is the struct of the project to build
type Project struct {
	RepoURL        string
	CloneCMD       string
	AuthType       string
	SSHPrivKeyPath string
	SSHPubKeyPath  string
	Username       string
	Password       string
	PromptForPWD   bool
}

// Environment is the config for the environment
type Environment struct {
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

func readConfig(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

// Load reads the config and returns a Config struct
func Load(path, clonePath string) (Config, error) {
	var conf Config
	confData, readErr := readConfig(path)
	if readErr != nil {
		return conf, readErr
	}
	if _, pErr := toml.Decode((string)(confData), &conf); pErr != nil {
		return conf, pErr
	}
	for i := range conf.Services {
		if strings.Contains(conf.Services[i].Path, ".") {
			conf.Services[i].Path = strings.Replace(conf.Services[i].Path, ".", clonePath, 1)
		}
	}
	return conf, nil
}
