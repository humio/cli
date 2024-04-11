package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newRolesDeleteCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "delete [flags] <role>",
		Short: "Delete a role.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			repo := args[0]
			client := NewApiClient(cmd)

			roleId, err := client.Roles().GetRoleID(repo)
			exitOnError(cmd, err, "Error getting role ID")
			if roleId == "" {
				exitOnError(cmd, fmt.Errorf("role %q not found", repo), "Role not found")
			}

			err = client.Roles().Delete(repo)
			exitOnError(cmd, err, "Error removing repository")

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully deleted repository: %q\n", repo)
		},
	}

	return &cmd
}
