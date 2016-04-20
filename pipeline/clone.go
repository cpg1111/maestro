package pipeline

import (
	"fmt"
	"log"
	"os"

	"github.com/cpg1111/maestro/config"
	"github.com/cpg1111/maestro/credentials"

	git "gopkg.in/libgit2/git2go.v22"
)

// Project is a struct for the Project in the pipeline
type Project struct {
	conf      config.Project
	State     string
	ABSPath   string
	Services  map[string]*Service
	creds     *credentials.RawCredentials
	gitCreds  *git.Cred
	clonePath string
	CloneOpts *git.CloneOptions
}

// New returns a new instance of a pipeline project
func New(conf *config.Config, creds *credentials.RawCredentials, clonePath, branch string) *Project {
	newServices := make(map[string]*Service)
	for i := range conf.Services {
		newServices[conf.Services[i].Name] = NewService(conf.Services[i], creds)
	}
	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		panic(cwdErr)
	}
	gitCreds := creds.ToGitCredentials()
	cloneOpts := &git.CloneOptions{
		RemoteCallbacks: &git.RemoteCallbacks{
			CredentialsCallback: func(url string, username string, allowedTypes git.CredType) (git.ErrorCode, *git.Cred) {
				log.Println(url, username)
				return 0, &gitCreds
			},
			CertificateCheckCallback: func(cert *git.Certificate, valid bool, hostname string) git.ErrorCode {
				return 0
			},
		},
		CheckoutOpts:   &git.CheckoutOpts{},
		Bare:           true,
		CheckoutBranch: branch,
	}
	//conf.CloneOpts.RemoteCreateCallback = createRemote
	return &Project{
		conf:      conf.Project,
		State:     "Pending",
		ABSPath:   fmt.Sprintf("%s/%s", cwd, clonePath),
		Services:  newServices,
		creds:     creds,
		gitCreds:  &gitCreds,
		clonePath: clonePath,
		CloneOpts: cloneOpts,
	}
}

// Clone clones a git repo
func (p *Project) Clone(opts *git.CloneOptions) (*git.Repository, error) {
	return git.Clone(p.conf.RepoURL, p.clonePath, opts)
}

// Unpack Git repo
func (p *Project) Unpack() {

}
