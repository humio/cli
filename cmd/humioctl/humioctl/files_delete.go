package humioctl

import (
	"fmt"
	"github.com/humio/cli/cmd/humioctl/internal/helpers"
	"github.com/spf13/cobra"
)

func newFilesDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <view-name> <file-name>",
		Short: "Delete an uploaded file.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			viewName := args[0]
			fileName := args[1]
			client := NewApiClient(cmd)

			err := client.Files().Delete(viewName, fileName)
			helpers.ExitOnError(cmd, err, "Error deleting file")

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully deleted file %q in repo %q\n", fileName, viewName)
		},
	}

	return cmd
}
