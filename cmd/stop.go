package cmd

import (
	"log"
	"os"
	"os/exec"

	c "github.com/hanks/awsudo-go/configs"
)

// Stop command is to stop agent server process running in background.
// And also will do some clean up tasks, like removing socket files
func Stop(path string) {
	cmds := []string{"pkill", "-SIGINT", "-f", c.AppName}
	_, err := exec.Command(cmds[0], cmds[1:]...).Output()

	if err != nil {
		// backoff to kill process run by main.go, mainly in dev env
		cmds = []string{"pkill", "-SIGINT", "-f", "go", "run", "main.go"}
		_, err = exec.Command(cmds[0], cmds[1:]...).Output()
		if err != nil {
			log.Println("Agent server process is not found, skip.")
		}
	}

	if err == nil {
		log.Println("Stop agent server process ok.")
	}

	if _, err = os.Stat(c.SocketFile); !os.IsNotExist(err) {
		err = os.Remove(c.SocketFile)
		if err != nil {
			log.Fatalf("Remove agent server socket file error, please check %v", err)
		} else {
			log.Println("Remove agent server socket file ok.")
		}
	} else {
		log.Println("Agent server socket file is not existed, skip.")
	}

	log.Println("Agent server stop ok.")
}
