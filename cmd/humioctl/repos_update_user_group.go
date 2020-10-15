package main

import (
	"fmt"
	"github.com/humio/cli/api"
	"github.com/spf13/cobra"
	"strings"
)

func newReposUpdateUserGroupCmd() *cobra.Command {
	var groups []string
	cmd := cobra.Command{
		Use:   "update-user-group [flags] <repo> <username>",
		Short: "Updates the users permissions to a repository based on default groups",
		Args:  cobra.ExactArgs(2),
		Run: WrapRun(func(cmd *cobra.Command, args []string) (humioResultType, error) {
			repoName := args[0]
			userName := args[1]

			var defaultGroups []api.DefaultGroupEnum
			for _, group := range groups {
				var defaultGroup api.DefaultGroupEnum
				if !defaultGroup.ParseString(group) {
					return nil, fmt.Errorf("the group %q was not valid (must be either 'Member', 'Admin' or 'Eliminator')", group)
				}
				defaultGroups = append(defaultGroups, defaultGroup)
			}

			client := NewApiClient(cmd)
			apiErr := client.Repositories().UpdateUserGroup(repoName, userName, defaultGroups...)
			if apiErr != nil {
				return nil, fmt.Errorf("error adding user: %w", apiErr)
			}

			return fmt.Sprintf("User %q's groups in repository %q changed to %s", userName, repoName, strings.Join(groups, ", ")), nil
		}),
	}
	cmd.Flags().StringSliceVarP(&groups, "groups", "g", []string{api.DefaultGroupEnumMember.String()}, "the groups that the user should be added in")

	return &cmd
}
