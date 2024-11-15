package main

import (
	"github.com/spf13/cobra"
)

func newTokensRotateApiTokenCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "rotate-api-token <token-id>",
		Short: "Rotate and retrieve an API token [Root Only]",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tokenID := args[0]

			client := NewApiClient(cmd)
			newToken, apiErr := client.Tokens().Rotate(tokenID)
			exitOnError(cmd, apiErr, "Error updating token")

			cmd.Printf("New API Token: %s\n", newToken)
		},
	}

	return &cmd
}
