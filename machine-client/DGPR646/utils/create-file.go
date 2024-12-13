package utils

import (
	"os"
	"path"
)

func CreateOrAppendFile(filepath string) (*os.File, error) {
	// Create of chmod main directory
	dirPath := path.Dir(filepath)
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0777)
		if err != nil {
			return nil, err
		}

		err = os.Chmod(dirPath, 0777)
		if err != nil {
			return nil, err
		}
	} else {
		err = os.Chmod(dirPath, 0777)
		if err != nil {
			return nil, err
		}
	}

	// Create file
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	// Chmod file
	err = f.Chmod(0666)
	return f, err
}
