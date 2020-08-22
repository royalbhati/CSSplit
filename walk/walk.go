package walk

import (
	"fmt"
	"os"
	"strings"

	"github.com/karrick/godirwalk"
)

//GetDir takes in the project dir and return the list of valid files
func GetDir(path string, exclude, ext []string) ([]string, error) {
	directories := []string{}
	dirname := path
	err := godirwalk.Walk(dirname, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if excludeFile(exclude, osPathname) {
				return godirwalk.SkipThis
			}
			if validExt(ext, osPathname) {
				directories = append(directories, osPathname)
			}
			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {

			return godirwalk.SkipNode
		},
		Unsorted: true,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return nil, err
	}
	return directories, nil
}

func validExt(ext []string, pathName string) bool {
	for _, v := range ext {
		if strings.HasSuffix(pathName, v) {
			return true
		}
	}
	return false
}

func excludeFile(excluded []string, pathName string) bool {
	if strings.Contains(pathName, "node_modules") || strings.Contains(pathName, ".git") {
		return true
	}
	for _, v := range excluded {
		if strings.Contains(pathName, v) {
			return true
		}
	}
	return false
}
