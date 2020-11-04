package main

import "github.com/spf13/cobra"

func newGroupsRemoveUserCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove-user <group-id> <user-id>",
		Short: "Remove user from group.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			err := client.Groups().RemoveUserFromGroup(args[0], args[1])
			exitOnError(cmd, err, "error removing user from group")
		},
	}
}