package prompt

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

func Ask(question string) (string, error) {
	var answer string
	fmt.Print("   " + question + ": ")
	n, err := fmt.Scanln(&answer)

	if n == 0 {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return answer, nil
}

func Confirm(text string) bool {
	fmt.Print("   " + text + " [Y/n]: ")

	reader := bufio.NewReader(os.Stdin)

	for {
		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "" || response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

func AskSecret(question string) (string, error) {
	fmt.Print("   " + question + ": ")
	bytes, err := terminal.ReadPassword(0)

	if err != nil {
		return "", err
	}

	fmt.Print("***********************\n")
	return string(bytes), nil
}

func Output(text string) {
	fmt.Println("  ", text)
}

func Title(text string) {
	c := "[underline][bold]" + text + "[reset]"
	Output(Colorize(c))
}

func Description(text string) {
	c := "[gray]" + text + "[reset]"
	Output(Colorize(c))
}

func Error(text string) {
	c := "[red]" + text + "[reset]"
	Output(Colorize(c))
}

func Info(text string) {
	c := "[purple]" + text + "[reset]"
	Output(Colorize(c))
}

func Colorize(text string) string {
	replacer := strings.NewReplacer(
		"[reset]", "\x1b[0m",
		"[gray]", "\x1b[38;5;249m",
		"[purple]", "\x1b[38;5;129m",
		"[bold]", "\x1b[1m",
		"[red]", "\x1b[38;5;1m",
		"[green]", "\x1b[38;5;2m",
		"[underline]", "\x1b[4m",
	)

	return replacer.Replace(text)
}

func Owl() string {
	return `    , ,
   (O,o)
   |)__)
   -”-”-
`
}
