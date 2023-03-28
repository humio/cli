package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newGroupsDelete() *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "delete group",
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			groupID := args[0]

			err := client.Groups().Delete(groupID)
			exitOnError(cmd, err, "Error deleting groups")

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully created group with GroupID: %s", groupID)
		},
	}
}
