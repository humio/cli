package main

import (
	"fmt"
	"github.com/humio/cli/cmd/internal/format"
	"github.com/spf13/cobra"
	"strings"
)

func newParsersShowCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "show <repo> <parser>",
		Short: "Show details for a parser in a repository.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			repoName := args[0]
			parserName := args[1]
			client := NewApiClient(cmd)

			parser, err := client.Parsers().Get(repoName, parserName)
			exitOnError(cmd, err, "Error fetching parser")

			details := [][]format.Value{
				{format.String("ID"), format.String(parser.ID)},
				{format.String("Name"), format.String(parser.Name)},
				{format.String("Script"), format.String(parser.Script)},
				{format.String("TagFields"), format.String(strings.Join(parser.FieldsToTag, "\n"))},
				{format.String("FieldsToBeRemovedBeforeParsing"), format.String(strings.Join(parser.FieldsToBeRemovedBeforeParsing, "\n"))},
				{format.String("TestCases"), format.String(fmt.Sprintf("%+v", parser.TestCases))},
			}

			printDetailsTable(cmd, details)
		},
	}

	return &cmd
}
