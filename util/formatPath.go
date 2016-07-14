package util

import (
	"fmt"
	"strings"
)

// FmtDiffPath formats a path for diffing
func FmtDiffPath(clonePath, srvPath string) (newStr string) {
	if strings.Index(srvPath, clonePath) == -1 && strings.Index(srvPath, clonePath[1:]) > -1 {
		clonePath = clonePath[1:]
	}
	if clonePath[len(clonePath)-1] != '/' {
		clonePath = fmt.Sprintf("%s/", clonePath)
	}
	newStr = strings.Replace(srvPath, clonePath, "", 1)
	if newStr[len(newStr)-1] != '/' {
		newStr = fmt.Sprintf("%s/", newStr)
	}
	return
}

func FmtClonePath(clonePath *string) *string {
	clPath := *clonePath
	if clPath[len(clPath)-1] == '/' {
		clPath = clPath[0:(len(clPath) - 1)]
		return &clPath
	}
	return clonePath
}
