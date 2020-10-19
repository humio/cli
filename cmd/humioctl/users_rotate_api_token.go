package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

func newUsersRotateApiTokenCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "rotate-api-token",
		Short: "Rotate and retrieve a user's API token [Root Only]",
		Args:  cobra.ExactArgs(1),
		Run: WrapRun(func(cmd *cobra.Command, args []string) (humioResultType, error) {
			userID := args[0]

			client := NewApiClient(cmd)
			newToken, apiErr := client.Users().RotateUserApiTokenAndGet(userID)
			if apiErr != nil {
				return nil, fmt.Errorf("error rotating api token: %w", apiErr)
			}

			return newToken, nil
		}),
	}

	return &cmd
}
