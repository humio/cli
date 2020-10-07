package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

func newGroupsList() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List groups",
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			groups, err := client.Groups().List()
			exitOnError(cmd, err, "error listing groups")

			rows := make([]string, len(groups))
			for i, group := range groups {
				rows[i] = fmt.Sprintf("%s | %s", group.ID, group.DisplayName)
			}

			printTable(cmd, append([]string{
				"ID | Display name"},
				rows...,
			))
		},
	}
}