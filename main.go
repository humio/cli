package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	// "github.com/skratchdot/open-golang/open"
	"github.com/humio/cli/command"
	mit "github.com/mitchellh/cli"
	"golang.org/x/crypto/ssh/terminal"
)

////////////////////////////////////////////////////////////////////////////////
///// Globals //////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// This is set by the release script.
var version = "master"
var client = &http.Client{}

type server struct {
	URL   string
	Token string
	Repo  string
}

var (
	// Hidden hides the commands from both help and autocomplete. Commands that
	// users should not be running should be placed here, versus hiding
	// subcommands from the main help, which should be filtered out of the
	// commands above.
	hidden = []string{}

	// aliases is the list of aliases we want users to be aware of. We hide
	// these form the help output but autocomplete them.
	aliases = []string{}

	// Common commands are grouped separately to call them out to operators.
	commonCommands = []string{
		"parsers",
	}
)

////////////////////////////////////////////////////////////////////////////////
///// main function ////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func main() {
	os.Exit(Run(os.Args[1:]))
}

func Run(args []string) int {
	// Parse flags into env vars for global use
	args = setupEnv(args)

	// Create the meta object
	metaPtr := new(command.Meta)

	// Don't use color if disabled
	color := true
	if os.Getenv(command.EnvHumioCLINoColor) != "" {
		color = false
	}

	isTerminal := terminal.IsTerminal(int(os.Stdout.Fd()))
	metaPtr.Ui = &mit.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	// The Nomad agent never outputs color
	agentUi := &mit.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	// Only use colored UI if stdout is a tty, and not disabled
	if isTerminal && color {
		metaPtr.Ui = &mit.ColoredUi{
			ErrorColor: mit.UiColorRed,
			WarnColor:  mit.UiColorYellow,
			InfoColor:  mit.UiColorGreen,
			Ui:         metaPtr.Ui,
		}
	}

	commands := command.Commands(metaPtr, agentUi)
	cli := &mit.CLI{
		Name:     "humio",
		Version:  version,
		Args:     args,
		Commands: commands,
		// HiddenCommands:             hidden,
		Autocomplete:               true,
		AutocompleteNoDefaultFlags: true,
		HelpFunc: groupedHelpFunc(
			mit.BasicHelpFunc("humio"),
		),
		HelpWriter: os.Stdout,
	}

	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	return exitCode
}

func groupedHelpFunc(f mit.HelpFunc) mit.HelpFunc {
	return func(commands map[string]mit.CommandFactory) string {
		var b bytes.Buffer
		tw := tabwriter.NewWriter(&b, 0, 2, 6, ' ', 0)

		fmt.Fprintf(tw, "Usage: humio [-version] [-help] [-autocomplete-(un)install] <command> [args]\n\n")
		fmt.Fprintf(tw, "Common commands:\n")
		for _, v := range commonCommands {
			printCommand(tw, v, commands[v])
		}

		// Filter out common commands and aliased commands from the other
		// commands output
		otherCommands := make([]string, 0, len(commands))
		for k := range commands {
			found := false
			for _, v := range commonCommands {
				if k == v {
					found = true
					break
				}
			}

			for _, v := range aliases {
				if k == v {
					found = true
					break
				}
			}

			if !found {
				otherCommands = append(otherCommands, k)
			}
		}
		sort.Strings(otherCommands)

		fmt.Fprintf(tw, "\n")
		fmt.Fprintf(tw, "Other commands:\n")
		for _, v := range otherCommands {
			printCommand(tw, v, commands[v])
		}

		tw.Flush()

		return strings.TrimSpace(b.String())
	}
}

func printCommand(w io.Writer, name string, cmdFn mit.CommandFactory) {
	cmd, err := cmdFn()
	if err != nil {
		panic(fmt.Sprintf("failed to load %q command: %s", name, err))
	}
	fmt.Fprintf(w, "    %s\t%s\n", name, cmd.Synopsis())
}

// setupEnv parses args and may replace them and sets some env vars to known
// values based on format options
func setupEnv(args []string) []string {
	noColor := false
	for _, arg := range args {
		// Check if color is set
		if arg == "-no-color" || arg == "--no-color" {
			noColor = true
		}
	}

	// Put back into the env for later
	if noColor {
		os.Setenv(command.EnvHumioCLINoColor, "true")
	}

	return args
}

////////////////////////////////////////////////////////////////////////////////
// Utils
////////////////////////////////////////////////////////////////////////////////

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
