package humioctl

import (
	"github.com/humio/cli/cmd/humioctl/internal/helpers"
	"github.com/spf13/cobra"
)

func newUsersRotateApiTokenCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "rotate-api-token <user-id>",
		Short: "Rotate and retrieve a user's API token [Root Only]",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			userID := args[0]

			client := NewApiClient(cmd)
			newToken, apiErr := client.Users().RotateUserApiTokenAndGet(userID)
			helpers.ExitOnError(cmd, apiErr, "Error updating user")

			cmd.Printf("New API Token: %s\n", newToken)
		},
	}

	return &cmd
}
