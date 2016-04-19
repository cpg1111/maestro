package pipeline

import (
	"os"
	"os/exec"

	"github.com/cpg1111/maestro/config"
	"github.com/cpg1111/maestro/credentials"

	git "gopkg.in/libgit2/git2go.v22"
)

type Service struct {
	config.Service
	Diff     bool
	State    string
	creds    credentials.RawCredentials
	gitCreds git.Cred
}

func NewService(srv config.Service, creds credentials.RawCredentials) *Service {
	return &Service{
		srv,
		Diff:     false,
		State:    "Pending",
		creds:    creds,
		gitCreds: creds.ToGitCredentials(),
	}
}

func (s *Service) ShouldRunPipeLine() (bool, error) {

}

type Project struct {
	config.Project
	State     string
	ABSPath   string
	Services  map[string]*Service
	creds     credentials.RawCredentials
	gitCreds  git.Cred
	clonePath string
}

func New(conf config.Config, creds credentials.RawCredentials, clonePath string) *Project {
	newServices := make(map[string]*Service)
	for i := range conf.Services {
		newServices[conf.Services[i].Name] = NewService(conf.Services[i], creds)
	}
	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		panic(cwdErr)
	}
	return &Project{
		conf.Project,
		State:     "Pending",
		ABSPath:   fmt.Sprintf("%s%s", cwd, clonePath),
		Services:  newServices,
		creds:     creds,
		gitCreds:  creds.ToGitCredentials,
		clonePath: clonePath,
	}
}

func (p *Project) Clone(opts *git.CloneOptions) (git.Repository, error) {
	return git.Clone(p.RepoURL, p.clonePath, opts)
}
