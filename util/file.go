package util

import (
	"io/ioutil"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

//GetAllFilesInDir returns all files in a directory
func GetAllFile(pathname string) ([]string, error) {
	var s []string
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		logrus.Error("Read File Error")
		return s, err
	}

	for _, fi := range rd {
		if !fi.IsDir() {
			fullName := filepath.Join(pathname, fi.Name())
			s = append(s, fullName)
		}
	}
	return s, nil
}
