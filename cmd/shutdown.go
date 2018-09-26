package cmd

import (
	"log"
	"os"
	"os/exec"

	c "github.com/hanks/awsudo-go/configs"
)

func Shutdown(path string) {
	_, err := exec.Command("pkill", "-SIGINT", c.AppName).Output()

	if err != nil {
		log.Fatalf("Shutdown agent server in background error, please check %v", err)
	}

	if _, err := os.Stat(c.SocketFile); !os.IsNotExist(err) {
		err = os.Remove(c.SocketFile)
		if err != nil {
			log.Fatalf("Remove agent server socket file error, please check %v", err)
		}
	}

	log.Println("Shutdown is ok.")
}
