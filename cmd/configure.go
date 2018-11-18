package cmd

import (
	"log"
	"os"
	"path/filepath"

	c "github.com/hanks/awsudo-go/configs"
	"github.com/hanks/awsudo-go/pkg/parser"
	"github.com/hanks/awsudo-go/utils"
)

// RunConfigure command is to config awsudo with 'aws configure style',
// And will save configs into ~/.awsudo.conf
func RunConfigure(path string) {
	if path == "" {
		path = c.DefaultConfPath
	}
	absPath, err := utils.GetAbsPath(path)
	if err != nil {
		log.Fatalf("Config (%s) load error: %v", absPath, err)
	}

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

	err = conf.InputConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = conf.WriteConfig(absPath)
	if err != nil {
		log.Fatal(err)
	}
}
