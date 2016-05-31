package credentials

import (
	"testing"

	"github.com/cpg1111/maestro/config"
)

func TestToGitCredentials(t *testing.T) {
	conf, loadErr := config.Load("../test_conf.toml", ".")
	if loadErr != nil {
		t.Error(loadErr)
	}
	project := &conf.Project
	testCreds, createErr := NewCreds(project)
	if createErr != nil {
		t.Error(createErr)
	}
	if testCreds.SSHPrivKey == "" {
		t.Error("Expected to load private key, did not load private key")
	}
	if testCreds.SSHPubKey == "" {
		t.Error("Expected to load public key, did not load private key")
	}
	gitCreds := testCreds.ToGitCredentials()
	if !gitCreds.HasUsername() {
		t.Error("Could not create git credentials")
	}
}
