package main

import (
	"os"

	"github.com/humio/cli/api"
	"github.com/spf13/cobra"
)

func newReposUpdateUserGroupCmd() *cobra.Command {
	var groups []string
	cmd := cobra.Command{
		Use:   "update-user-group [flags] <repo> <username>",
		Short: "Updates the users permissions to a repository based on default groups",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			repoName := args[0]
			userName := args[1]
			client := NewApiClient(cmd)

			var defaultGroups []api.DefaultGroupEnum
			for _, group := range groups {
				var defaultGroup api.DefaultGroupEnum
				if !defaultGroup.ParseString(group) {
					cmd.Println("The group '%s' was not valid (must be either 'Member', 'Admin' or 'Eliminator')")
					os.Exit(1)
				}
				defaultGroups = append(defaultGroups, defaultGroup)
			}

			err := client.Repositories().UpdateUserGroup(repoName, userName, defaultGroups...)
			exitOnError(cmd, err, "Error adding user")
		},
	}
	cmd.Flags().StringSliceVarP(&groups, "groups", "g", []string{api.DefaultGroupEnumMember.String()}, "the groups that the user should be added in")

	return &cmd
}
