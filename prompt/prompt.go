package prompt

import (
	"fmt"

	"golang.org/x/crypto/ssh/terminal"
)

func Ask(question string) (string, error) {
	var answer string
	fmt.Print(question + ": ")
	n, err := fmt.Scanln(&answer)

	if n == 0 {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return answer, nil
}

func AskSecret(question string) (string, error) {
	fmt.Print(question + ": ")
	bytes, err := terminal.ReadPassword(0)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
