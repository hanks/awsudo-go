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
	absPath := utils.GetAbsPath(path)
	conf, err := parser.LoadConfig(absPath)
	if err != nil {
		log.Fatalf("Config (%s) load error: %v", absPath, err)
	}

	return conf
}

func RunAgentServer(path string) {
	conf := getConfig(path)
	log.Printf("Start agent server to handle new request, and will be expired after %d seconds", conf.Agent.Expiration)
	server := newServer(c.SocketFile, conf.Agent.Expiration)
	server.run()
}

func RunAgentClient(path string, roleName string) {
	// create subprocess to run agent server background
	ex, err := os.Executable()
	if err != nil {
		log.Fatalf("Can not find current executable binary.")
	}
	serverCmd := exec.Command(ex, "agent", "-c", path)
	err = serverCmd.Start()
	if err != nil {
		log.Fatalf("Start server process background err. Please check %v", err)
	}
	err = serverCmd.Process.Release()
	if err != nil {
		log.Fatalf("Can not detach server subprocess. Please check %v", err)
	}

	// run agent client foreground
	conf := getConfig(path)
	existed, _ := conf.GetARN(roleName)
	if !existed {
		log.Fatalf("Role (%s) is not support now, please try another one.", roleName)
	}

	client := newClient(c.SocketFile, conf, roleName)
	client.run()
}
