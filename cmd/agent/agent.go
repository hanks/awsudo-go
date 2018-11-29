package agent

import (
	"log"
	"os"
	"os/exec"

	c "github.com/hanks/awsudo-go/configs"
	"github.com/hanks/awsudo-go/pkg/parser"
	"github.com/hanks/awsudo-go/utils"
)

func getConfig(path string) *parser.Config {
	if path == "" {
		path = c.DefaultConfPath
	}
	absPath, err := utils.GetAbsPath(path)
	if err != nil {
		log.Fatalf("Config (%s) load error: %v", absPath, err)
	}
	conf, err := parser.LoadConfig(absPath)
	if err != nil {
		log.Fatalf("Config (%s) load error: %v", absPath, err)
	}

	return conf
}

// RunAgentServer is a wrapper to run awsudo agent server process in background
func RunAgentServer(path string) {
	conf := getConfig(path)
	server := newServer(c.SocketFile, conf.Agent.Expiration)
	server.run()
}

// RunAgentClient is a wrapper to run awsudo agent client
func RunAgentClient(configPath string, roleName string) {
	// create subprocess to run agent server background
	ex, err := os.Executable()
	if err != nil {
		log.Fatalf("Can not find current executable binary.")
	}
	serverCmd := exec.Command(ex, "start-agent", "-c", configPath)
	err = serverCmd.Start()
	if err != nil {
		log.Fatalf("Start server process background err. Please check %v", err)
	}
	err = serverCmd.Process.Release()
	if err != nil {
		log.Fatalf("Can not detach server subprocess. Please check %v", err)
	}

	// run agent client foreground
	conf := getConfig(configPath)
	existed, _ := conf.GetARN(roleName)
	if !existed {
		log.Fatalf("Role (%s) is not support now, please try another one.", roleName)
	}

	client := newClient(c.SocketFile, conf, roleName)
	client.run()
}
