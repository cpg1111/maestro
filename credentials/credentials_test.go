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
