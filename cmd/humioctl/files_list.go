package main

import (
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

			var rows [][]string

			for _, file := range files {
				rows = append(rows, []string{file.Name, file.ContentHash, file.ID})
			}

			printOverviewTable(cmd, []string{"Name", "Content Hash", "ID"}, rows)
		},
	}

	return cmd
}
