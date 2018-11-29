package utils

import (
	"fmt"
	"strconv"
	"strings"
	"syscall"

	"github.com/hanks/awsudo-go/configs"
	"golang.org/x/crypto/ssh/terminal"
)

type iScanner interface {
	Scan() bool
	Text() string
}

var readPassword = terminal.ReadPassword

// AskUserInput is to ask user to input account info for aws role
func AskUserInput(scanner iScanner) (string, string) {
	var name, pass string

	fmt.Print("Please enter login name: ")
	scanner.Scan()
	name = scanner.Text()

	for {
		fmt.Print("Please enter login password: ")
		onePass, _ := readPassword(int(syscall.Stdin))

		fmt.Print("\nPlease confirm login password: ")
		twoPass, _ := readPassword(int(syscall.Stdin))

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

// InputString is to accept user input to a string var
func InputString(scanner iScanner, original string, name string) string {
	// for UX, to print out None when the value is empty
	var input = original
	if original == "" {
		fmt.Printf("%s [%v]: ", name, "None")
	} else {
		fmt.Printf("%s [%v]: ", name, original)
	}

	scanner.Scan()
	text := scanner.Text()
	if text != "" {
		input = text
	}

	return input
}

// InputInt64 is to accept user input to a int64 var
func InputInt64(scanner iScanner, original int64, name string) (int64, error) {
	// for UX, to print out default value when the value is empty
	var input = original
	if original == 0 {
		if strings.Contains(name, "Duration") {
			input = configs.DefaultSessionDuration
		} else if strings.Contains(name, "Expiration") {
			input = configs.DefaultAgentExpiration
		}
	}
	fmt.Printf("%s [%v]: ", name, input)

	scanner.Scan()
	text := scanner.Text()
	if text != "" {
		n, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return 0, err
		}
		input = n
	}

	return input, nil
}
