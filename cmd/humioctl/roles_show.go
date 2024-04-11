package main

import (
	"github.com/spf13/cobra"
)

func newRolesShowCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "show <role>",
		Short: "Show details about a role.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			roleName := args[0]
			client := NewApiClient(cmd)

			role, err := client.Roles().Get(roleName)
			exitOnError(cmd, err, "Error fetching role")

			printRoleDetailsTable(cmd, *role)
		},
	}

	return &cmd
}
