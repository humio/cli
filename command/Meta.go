package command

import (
	"bufio"
	"flag"
	"io"
	"os"
	"strings"

	"github.com/humio/cli/api"
	"github.com/mitchellh/cli"
	"github.com/mitchellh/colorstring"
	"github.com/posener/complete"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	// Constants for CLI identifier length
	shortId = 8
	fullId  = 36
)

// FlagSetFlags is an enum to define what flags are present in the
// default FlagSet returned by Meta.FlagSet.
type FlagSetFlags uint

const (
	FlagSetNone    FlagSetFlags = 0
	FlagSetClient  FlagSetFlags = 1 << iota
	FlagSetDefault              = FlagSetClient
)

// Meta contains the meta-options and functionality that nearly every
// Nomad command inherits.
type Meta struct {
	Ui cli.Ui

	// These are set by the command line flags.
	flagAddress string

	// Whether to not-colorize output
	noColor bool

	// the api token used to access Humio's API
	token string
}

// FlagSet returns a FlagSet with the common flags that every
// command implements. The exact behavior of FlagSet can be configured
// using the flags as the second parameter, for example to disable
// server settings on the commands that don't talk to a server.
func (m *Meta) FlagSet(n string, fs FlagSetFlags) *flag.FlagSet {
	f := flag.NewFlagSet(n, flag.ContinueOnError)

	// FlagSetClient is used to enable the settings for specifying
	// client connectivity options.
	if fs&FlagSetClient != 0 {
		f.StringVar(&m.flagAddress, "address", "", "")
		f.BoolVar(&m.noColor, "no-color", false, "")
		f.StringVar(&m.token, "token", "", "")
	}

	// Create an io.Writer that writes to our UI properly for errors.
	// This is kind of a hack, but it does the job. Basically: create
	// a pipe, use a scanner to break it into lines, and output each line
	// to the UI. Do this forever.
	errR, errW := io.Pipe()
	errScanner := bufio.NewScanner(errR)
	go func() {
		for errScanner.Scan() {
			m.Ui.Error(errScanner.Text())
		}
	}()
	f.SetOutput(errW)

	return f
}

// AutocompleteFlags returns a set of flag completions for the given flag set.
func (m *Meta) AutocompleteFlags(fs FlagSetFlags) complete.Flags {
	if fs&FlagSetClient == 0 {
		return nil
	}

	return complete.Flags{
		"-address":  complete.PredictAnything,
		"-no-color": complete.PredictNothing,
		"-token":    complete.PredictAnything,
	}
}

// ApiClientFactory is the signature of a API client factory
type ApiClientFactory func() (*api.Client, error)

// Client is used to initialize and return a new API client using
// the default command line arguments and env vars.
func (m *Meta) Client() (*api.Client, error) {
	config := api.DefaultConfig()
	if m.flagAddress != "" {
		config.Address = m.flagAddress
	}

	return api.NewClient(config)
}

func (m *Meta) Colorize() *colorstring.Colorize {
	return &colorstring.Colorize{
		Colors:  colorstring.DefaultColors,
		Disable: m.noColor || !terminal.IsTerminal(int(os.Stdout.Fd())),
		Reset:   true,
	}
}

// generalOptionsUsage returns the help string for the global options.
func generalOptionsUsage() string {
	helpText := `
  -address=<addr>
    The address of the Humio server.
    Overrides the HUMIO_ADDR environment variable if set.
    Default = http://localhost:8080/
  -no-color
    Disables colored command output. Alternatively, HUMIO_CLI_NO_COLOR may be
    set.
  -token
    The API token to use to authenticate API requests with.
    Overrides the HUMIO_API_TOKEN environment variable if set.
`
	return strings.TrimSpace(helpText)
}

// funcVar is a type of flag that accepts a function that is the string given
// by the user.
type funcVar func(s string) error

func (f funcVar) Set(s string) error { return f(s) }
func (f funcVar) String() string     { return "" }
func (f funcVar) IsBoolFlag() bool   { return false }
