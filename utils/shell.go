package utils

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

// ExecCommand is to execute shell command in golang
func ExecCommand(cmds []string) error {
	cmdName := cmds[0]
	binPath, err := exec.LookPath(cmdName)
	if err != nil {
		log.Fatalf("%s is not found", cmdName)
	}
	cmds[0] = binPath
	// use syscall.Exec to replace process with new one
	err = syscall.Exec(binPath, cmds, os.Environ())
	return err
}
