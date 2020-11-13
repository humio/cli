package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

func newFilesListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list <view name>",
		Short: "List uploaded files in a view.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			files, err := client.Files().List(args[0])
			exitOnError(cmd, err, "error listing files")

			table := []string{"ID | Name | Content Hash"}

			for _, file := range files {
				table = append(table, fmt.Sprintf("%s | %s | %s", file.ID, file.Name, file.ContentHash))
			}

			printTable(cmd, table)
		},
	}

	return cmd
}

