package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newGroupsFind() *cobra.Command {
	return &cobra.Command{
		Use:   "find",
		Short: "find group by name",
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			group, err := client.Groups().Find(args[0])
			exitOnError(cmd, err, "Error finding group")

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully found group %s with id %s", group.DisplayName, group.ID)
		},
	}
}
