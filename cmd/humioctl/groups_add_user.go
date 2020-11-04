package main

import "github.com/spf13/cobra"

func newGroupsAddUserCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add-user <group-id> <user-id>",
		Short: "Add user to group.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			err := client.Groups().AddUserToGroup(args[0], args[1])
			exitOnError(cmd, err, "error adding user to group")
		},
	}
}