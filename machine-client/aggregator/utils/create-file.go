package utils

import (
	"os"
	"path"
)

func CreateOrReplaceFile(filepath string) (*os.File, error) {
	// Create of chmod main directory
	dirPath := path.Dir(filepath)
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0777)
		if err != nil {
			return nil, err
		}
	}

	// Create file
	return os.OpenFile(filepath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
}

func CreateOrAppendFile(filepath string) (*os.File, error) {
	// Create of chmod main directory
	dirPath := path.Dir(filepath)
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0777)
		if err != nil {
			return nil, err
		}
	}

	// Create file
	return os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
}
