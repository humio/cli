package main

import (
	"github.com/humio/cli/cmd/internal/format"
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

			rows := make([][]format.Value, len(groups))
			for i, group := range groups {
				rows[i] = []format.Value{
					format.String(group.DisplayName),
					format.String(group.ID),
				}
			}

			printOverviewTable(cmd, []string{"Display Name", "ID"}, rows)
		},
	}
}
