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
	if clonePath[len(clonePath)-1] != '/' && srvPath[len(srvPath)-1] == '/' {
		clonePath = fmt.Sprintf("%s/", clonePath)
	}
	newStr = strings.Replace(srvPath, clonePath, "", 1)
	fmt.Println(newStr, clonePath, len(newStr))
	if len(newStr) <= 1 {
		newStr = "*"
		return
	}
	if newStr[len(newStr)-1] != '/' {
		newStr = fmt.Sprintf("%s/", newStr)
	}
	return
}

// FmtClonePath formats clonePath to a uniform format regardless of input
func FmtClonePath(clonePath *string) *string {
	clPath := *clonePath
	if clPath[len(clPath)-1] == '/' {
		clPath = clPath[0:(len(clPath) - 1)]
		return &clPath
	}
	return clonePath
}
