package cmd

import (
	"log"
	"os"
	"path/filepath"

	c "github.com/hanks/awsudo-go/configs"
	"github.com/hanks/awsudo-go/pkg/parser"
	"github.com/hanks/awsudo-go/utils"
)

func RunConfigure(path string) {
	if path == "" {
		path = c.DefaultConfPath
	}
	absPath := utils.GetAbsPath(path)

	// create config directory if not existed
	absDir := filepath.Dir(absPath)
	if _, err := os.Stat(absDir); os.IsNotExist(err) {
		os.MkdirAll(absDir, os.ModePerm)
	}

	conf := new(parser.Config)
	if _, err := os.Stat(absPath); !os.IsNotExist(err) {
		// load existed config, and update
		conf, err = parser.LoadConfig(absPath)
		if err != nil {
			log.Fatalf("Config (%s) load error: %v", absPath, err)
		}
	}

	conf.InputConfig()
	parser.WriteConfig(absPath, conf)
}
