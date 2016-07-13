package util

import (
	"fmt"
	"strings"
)

// FmtDiffPath formats a path for diffing
func FmtDiffPath(clonePath, srvPath string) (newStr string) {
	if clonePath[len(clonePath)-1] != '/' {
		clonePath = fmt.Sprintf("%s/", clonePath)
	}
	newStr = strings.Replace(srvPath, clonePath, "", 1)
	if newStr[len(newStr)-1] != '/' {
		newStr = fmt.Sprintf("%s/", newStr)
	}
	return
}
