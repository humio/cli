package prompt

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

type Prompt struct {
	Out io.Writer
}

func NewPrompt(out io.Writer) *Prompt {
	return &Prompt{Out: out}
}

func (p *Prompt) Ask(question string) (string, error) {
	var answer string
	fmt.Fprint(p.Out, "  "+question+": ")
	n, err := fmt.Scanln(&answer)

	if n == 0 {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return answer, nil
}

func (p *Prompt) Confirm(text string) bool {
	p.Print(text + " [Y/n]: ")

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

func (p *Prompt) AskSecret(question string) (string, error) {
	p.Print(question + ": ")
	bytes, err := terminal.ReadPassword(int(syscall.Stdin))

	if err != nil {
		return "", err
	}

	fmt.Fprint(p.Out, "***********************\n")
	return string(bytes), nil
}

func (p *Prompt) BlankLine() {
	fmt.Fprintln(p.Out)
}

func (p *Prompt) Print(i ...interface{}) {
	fmt.Fprint(p.Out, i...)
}

func (p *Prompt) Printf(format string, i ...interface{}) {
	fmt.Fprintf(p.Out, format, i...)
}

func (p *Prompt) Title(text string) {
	c := "[underline][bold]" + text + "[reset]"
	p.Print(Colorize(c))
}

func (p *Prompt) Description(text string) {
	c := "[gray]" + text + "[reset]"
	p.Print(Colorize(c))
}

func (p *Prompt) Error(text string) {
	c := "[red]" + text + "[reset]"
	p.Print(Colorize(c))
}

func (p *Prompt) Info(text string) {
	c := "[purple]" + text + "[reset]"
	p.Print(Colorize(c))
}

func Colorize(text string) string {
	replacer := strings.NewReplacer(
		"[reset]", "\x1b[0m",
		"[gray]", "\x1b[38;5;249m",
		"[purple]", "\x1b[38;5;129m",
		"[bold]", "\x1b[1m",
		"[red]", "\x1b[38;5;1m",
		"[yellow]", "\x1b[33m",
		"[green]", "\x1b[38;5;2m",
		"[underline]", "\x1b[4m",
	)

	return replacer.Replace(text)
}

func (p *Prompt) List(items []string) string {
	var str string
	for _, value := range items {
		str += "  - " + value + "\n"
	}
	return str
}

func Owl() string {
	return `  , ,
   (O,o)
   |)__)
   -”-”-
`
}
