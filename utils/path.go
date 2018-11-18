package utils

import (
	"errors"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

// GetAbsPath is to get absolute path of file
func GetAbsPath(path string) (string, error) {
	if len(path) == 0 {
		return path, nil
	}

	absPath := ""
	err := errors.New("")
	if path[0] == '~' {
		absPath, err = homedir.Expand(path)
	} else {
		absPath, err = filepath.Abs(path)
	}

	return absPath, err
}
