package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"syscall"

	"github.com/hanks/awsudo-go/configs"
	"golang.org/x/crypto/ssh/terminal"
)

// AskUserInput is to ask user to input account info for aws role
func AskUserInput() (string, string) {
	var name, pass string

	fmt.Print("Please enter login name: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	name = scanner.Text()

	for {
		fmt.Print("Please enter login password: ")
		onePass, _ := terminal.ReadPassword(int(syscall.Stdin))

		fmt.Print("\nPlease confirm login password: ")
		twoPass, _ := terminal.ReadPassword(int(syscall.Stdin))

		if string(onePass) == string(twoPass) {
			pass = string(onePass)
			fmt.Println()
			break
		} else {
			fmt.Println("\nPassword is not the same, please input again.")
			fmt.Println()
		}
	}

	return name, pass
}

var scanner = bufio.NewScanner(os.Stdin)

// InputString is to accept user input to a string var
func InputString(s *string, name string) {
	v := *s
	if v == "" {
		v = "None"
	}
	fmt.Printf("%s [%v]: ", name, v)
	scanner.Scan()
	if scanner.Text() != "" {
		*s = scanner.Text()
	}
}

// InputInt64 is to accept user input to a int64 var
func InputInt64(v *int64, name string) {
	i := *v
	if i == 0 {
		i = configs.DefaultSessionDuration
		*v = i
	}
	fmt.Printf("%s [%v]: ", name, i)
	scanner.Scan()
	if scanner.Text() != "" {
		n, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			log.Fatalf("Please input an integer for %s.", name)
		}
		*v = n
	}
}
