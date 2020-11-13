package main

import "github.com/spf13/cobra"

func newFilesDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <view name> <file name>",
		Short: "Delete an uploaded file.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			err := client.Files().Delete(args[0], args[1])
			exitOnError(cmd, err, "error deleting file")
		},
	}

	return cmd
}

