package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newGroupsCreate() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "create group",
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			groupName := args[0]

			err := client.Groups().Create(groupName)
			exitOnError(cmd, err, "Error creating groups")

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully created group %s", groupName)
		},
	}
}
