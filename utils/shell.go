package utils

import (
	"os"
	"os/exec"
	"syscall"
)

var execv = syscall.Exec

// ExecCommand is to execute shell command in golang
func ExecCommand(cmds []string) error {
	cmdName := cmds[0]
	binPath, err := exec.LookPath(cmdName)
	if err != nil {
		return err
	}
	cmds[0] = binPath
	// use syscall.Exec to replace process with new one
	err = execv(binPath, cmds, os.Environ())
	return err
}
