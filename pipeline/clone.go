package pipeline

import (
	"fmt"
	"log"
	"os"

	"github.com/cpg1111/maestro/config"
	"github.com/cpg1111/maestro/credentials"

	pb "gopkg.in/cheggaaa/pb.v1"
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

func credCB(gitCreds *git.Cred) git.CredentialsCallback {
	return func(url string, username string, allowedTypes git.CredType) (git.ErrorCode, *git.Cred) {
		return 0, gitCreds
	}
}

func certCheckCB(cert *git.Certificate, valid bool, hostname string) git.ErrorCode {
	return 0
}

var (
	progbar  *pb.ProgressBar
	received uint
)

func handleProgress(stats git.TransferProgress) git.ErrorCode {
	if progbar == nil {
		progbar = pb.StartNew((int)(stats.TotalObjects))
	}
	newObjs := stats.ReceivedObjects - received
	for i := 0; i < (int)(newObjs); i++ {
		progbar.Increment()
	}
	received = stats.ReceivedObjects
	return git.ErrOk
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
			CredentialsCallback:      credCB(&gitCreds),
			CertificateCheckCallback: certCheckCB,
			TransferProgressCallback: handleProgress,
		},
		CheckoutOpts: &git.CheckoutOpts{
			Strategy: git.CheckoutSafeCreate,
		},
		Bare:           false,
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
	log.Println("Cloning Repo...")
	return git.Clone(p.conf.RepoURL, fmt.Sprintf("%s/", p.clonePath), opts)
}
