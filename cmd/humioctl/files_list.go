package main

import (
	"github.com/humio/cli/cmd/internal/format"
	"github.com/spf13/cobra"
)

func newFilesListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list <view-name>",
		Short: "List uploaded files in a view.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			files, err := client.Files().List(args[0])
			exitOnError(cmd, err, "Error listing files")

			var rows [][]format.Value

			for _, file := range files {
				rows = append(rows, []format.Value{
					format.String(file.Name),
					format.String(file.ContentHash),
					format.String(file.ID),
				})
			}

			format.PrintOverviewTable(cmd, []string{"Name", "Content Hash", "ID"}, rows)
		},
	}

	return cmd
}
