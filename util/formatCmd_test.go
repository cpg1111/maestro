package util

import (
	"fmt"
	"os"
	"testing"
)

func TestFormatCommand(t *testing.T) {
	path, pathErr := os.Getwd()
	if pathErr != nil {
		t.Error(pathErr)
	}
	cmd, err := FormatCommand("bash -c 'echo \"hello world ./\"'", path)
	if err != nil {
		t.Error(err)
	}
	expected := []string{"/bin/bash", "-c", fmt.Sprintf("'echo \"hello world %s\"'", path)}
	for i := range cmd.Args {
		if cmd.Args[i] != expected[i] {
			t.Errorf("expected %s, found %s", expected[i], cmd.Args[i])
		}
	}
}
