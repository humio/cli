package humioctl

import (
	"fmt"
	"github.com/humio/cli/cmd/humioctl/internal/helpers"
	"github.com/spf13/cobra"
)

func newGroupsRemoveUserCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove-user <group-id> <user-id>",
		Short: "Remove user from group.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			groupID := args[0]
			userID := args[1]
			client := NewApiClient(cmd)

			err := client.Groups().RemoveUserFromGroup(groupID, userID)
			helpers.ExitOnError(cmd, err, "Error removing user from group")

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully removed user %q to group %q\n", userID, groupID)
		},
	}
}
