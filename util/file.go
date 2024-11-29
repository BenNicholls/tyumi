// utility functions for dealing with files
package util

import (
	"os"
	"strings"
)

// GetFileList returns a list of all files in the provided directory. If ext is provided, it only includes files with
// that extension.
func GetFileList(dirPath, ext string) (files []string, err error) {
	files = make([]string, 0)

	dirContents, err := os.ReadDir(dirPath)
	if err != nil {
		return
	}

	for _, entry := range dirContents {
		if !entry.IsDir() {
			if ext != "" && !strings.HasSuffix(entry.Name(), ext) {
				continue
			}
			files = append(files, entry.Name())
		}
	}

	return
}
