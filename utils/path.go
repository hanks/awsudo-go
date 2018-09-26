package utils

import (
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

func GetAbsPath(path string) string {
	absPath, err := homedir.Expand(path)
	if err == nil {
		return absPath
	}

	absPath, _ = filepath.Abs(path)
	return absPath
}
