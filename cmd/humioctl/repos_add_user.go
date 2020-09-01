package main

import (
	"os"

	"github.com/spf13/cobra"
)

func newReposAddUserCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "add-user [flags] <repo> <username> [group - default: Member]",
		Short: "Adds a user to a repository",
		Args:  cobra.RangeArgs(2, 3),
		Run: func(cmd *cobra.Command, args []string) {
			repoName := args[0]
			userName := args[1]
			group := "Member"

			if len(args) == 3 {
				group = args[2]
			}

			if group != "Member" && group != "Admin" && group != "Eliminator" {
				cmd.Println("group was invalid must be one of: Member, Admin or Eliminator")
				os.Exit(1)
			}

			client := NewApiClient(cmd)
			apiErr := client.Repositories().AddUser(repoName, userName, group)
			exitOnError(cmd, apiErr, "error adding user")
		},
	}

	return &cmd
}
