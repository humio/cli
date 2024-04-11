package main

import (
	"sort"

	"github.com/humio/cli/api"
	"github.com/humio/cli/cmd/internal/format"
	"github.com/spf13/cobra"
)

func newRolesListCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "list [flags]",
		Short: "List roles",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			roles, err := client.Roles().List()
			exitOnError(cmd, err, "Error fetching roles")

			sort.Slice(roles, func(i, j int) bool {
				var a, b api.Role
				a = roles[i]
				b = roles[j]

				return a.DisplayName < b.DisplayName
			})

			rows := make([][]format.Value, len(roles))
			for i, view := range roles {
				rows[i] = []format.Value{
					format.String(view.DisplayName),
					format.String(view.Description),
					format.Int(view.GroupsCount),
					format.Int(view.UsersCount),
					format.String(view.ID),
				}
			}

			printOverviewTable(cmd, []string{"Name", "Description", "Groups", "Users", "ID"}, rows)
		},
	}

	return &cmd
}
