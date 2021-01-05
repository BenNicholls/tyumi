//utility functions for dealing with files

package util

import "io/ioutil"
import "strings"


//GetFileList returns a list of all files in the provided directory. If ext is provided, it only includes
//files with that extension.
func GetFileList(dirPath, ext string) (files []string, err error) {
	files = make([]string, 0)

	dirContents, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return
	} 
	
	for i, file := range dirContents {
		if !file.IsDir() {
			if ext != "" && !strings.HasSuffix(dirContents[i].Name(), ext) {
				continue
			}
			files = append(files, dirContents[i].Name())
		}
	}
	
	return
}