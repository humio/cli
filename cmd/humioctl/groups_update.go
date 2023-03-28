package main

import (
	"fmt"

	"github.com/humio/cli/api"
	"github.com/spf13/cobra"
)

func newGroupsUpdate() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "update group",
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			groupID := args[0]
			displayName := args[1]
			changeSet := api.GroupChangeSet{
				GroupID: &groupID,
			}

			err := client.Groups().Update(displayName, changeSet)
			exitOnError(cmd, err, "Error updating groups")

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully updated group %s, with GroupID: %s", displayName, groupID)
		},
	}
}
