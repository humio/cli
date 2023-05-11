package main

import (
	"github.com/spf13/cobra"
)

func newFilesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "files",
		Short: "Manage files",
	}

	cmd.AddCommand(newFilesListCmd())
	cmd.AddCommand(newFilesDeleteCmd())
	cmd.AddCommand(newFilesUploadCmd())
	cmd.AddCommand(newFilesDownloadCmd())

	return cmd
}
