package utils

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

func ExecCommand(cmds []string) error {
	cmdName := cmds[0]
	cmdArgs := cmds[1:]
	cmd := exec.Command(cmdName, cmdArgs...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer stdout.Close()
	if err := cmd.Start(); err != nil {
		return err
	}

	// output command log in real-time style
	buf := make([]byte, 1024)
	for {
		_, err := stdout.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		r := bufio.NewReader(stdout)
		line, _, _ := r.ReadLine()
		fmt.Println(string(line))
	}

	return nil
}
