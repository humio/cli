package main

import (
	"fmt"
	"github.com/humio/cli/cmd/humioctl/internal/helpers"
	"github.com/spf13/cobra"
)

func newGroupsAddUserCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add-user <group-id> <user-id>",
		Short: "Add user to group.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			groupID := args[0]
			userID := args[1]
			client := NewApiClient(cmd)

			err := client.Groups().AddUserToGroup(groupID, userID)
			helpers.ExitOnError(cmd, err, "Error adding user to group")

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully added user %q to group %q\n", userID, groupID)
		},
	}
}
