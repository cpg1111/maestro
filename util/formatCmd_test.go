package util

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestFmtCommand(t *testing.T) {
	path, pathErr := os.Getwd()
	if pathErr != nil {
		t.Error(pathErr)
	}
	cmd, err := FmtCommand("bash -c 'echo \"hello world ./\"'", path)
	if err != nil {
		t.Error(err)
	}
	expected := []string{"/bin/bash", "-c", fmt.Sprintf("'echo \"hello world %s/\"'", path)}
	for i := range cmd.Args {
		if strings.Compare(cmd.Args[i], expected[i]) != 0 {
			t.Errorf("expected %s, found %s", expected[i], cmd.Args[i])
		}
	}
}
