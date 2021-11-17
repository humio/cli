package main

import (
	"github.com/spf13/cobra"
)

func newGroupsList() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List groups",
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			groups, err := client.Groups().List()
			exitOnError(cmd, err, "Error listing groups")

			rows := make([][]string, len(groups))
			for i, group := range groups {
				rows[i] = []string{group.DisplayName, group.ID}
			}

			printOverviewTable(cmd, []string{"Display Name", "ID"}, rows)
		},
	}
}
