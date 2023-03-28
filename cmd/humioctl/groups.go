package main

import "github.com/spf13/cobra"

func newGroupsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "groups",
		Short: "Manage groups",
	}

	cmd.AddCommand(newGroupsAddUserCmd())
	cmd.AddCommand(newGroupsRemoveUserCmd())
	cmd.AddCommand(newGroupsList())
	cmd.AddCommand(newGroupsFind())
	cmd.AddCommand(newGroupsCreate())
	cmd.AddCommand(newGroupsDelete())
	cmd.AddCommand(newGroupsUpdate())

	return cmd
}
