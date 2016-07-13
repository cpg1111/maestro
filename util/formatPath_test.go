package util

import (
	"testing"
)

func TestFmtDiffPath(t *testing.T) {
	testVals := []string{
		"/usr/lib",
		"usr/lib/",
		"/usr/lib/",
		"usr/lib",
	}
	for i := range testVals {
		res := FmtDiffPath("/usr/", testVals[i])
		if res != "lib/" {
			t.Errorf("expected lib/ found %s", res)
		}
	}
}
