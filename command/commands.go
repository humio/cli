package command

import (
	"os"

	"github.com/mitchellh/cli"
)

const (
	// EnvHumioCLINoColor is an env var that toggles colored UI output.
	EnvHumioCLINoColor = `HUMIO_CLI_NO_COLOR`
	// EnvHumioFormat is the output format
	EnvHumioFormat = `HUMIO_FORMAT`
)

func Commands(metaPtr *Meta, agentUi cli.Ui) map[string]cli.CommandFactory {
	if metaPtr == nil {
		metaPtr = new(Meta)
	}

	meta := *metaPtr
	if meta.Ui == nil {
		meta.Ui = &cli.BasicUi{
			Reader:      os.Stdin,
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		}
	}

	all := map[string]cli.CommandFactory{
		"parsers": func() (cli.Command, error) {
			return &ParsersCommand{
				Meta: meta,
			}, nil
		},
		"parsers list": func() (cli.Command, error) {
			return &ParsersListCommand{
				Meta: meta,
			}, nil
		},
		"parsers rm": func() (cli.Command, error) {
			return &ParsersRemoveCommand{
				Meta: meta,
			}, nil
		},
		"parsers install": func() (cli.Command, error) {
			return &ParsersInstallCommand{
				Meta: meta,
			}, nil
		},
		"parsers export": func() (cli.Command, error) {
			return &ParsersExportCommand{
				Meta: meta,
			}, nil
		},
	}

	return all
}

// NamedCommand is a interface to denote a commmand's name.
type NamedCommand interface {
	Name() string
}
