package humioctl

import (
	"github.com/humio/cli/cmd/humioctl/internal/helpers"
	"github.com/spf13/cobra"
	"io"
	"os"
)

func newFilesDownloadCmd() *cobra.Command {
	var (
		saveAs string
	)

	cmd := &cobra.Command{
		Use:  "download <view-name> <file-name>",
		Long: `Download a file.`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			viewName := args[0]
			fileName := args[1]
			client := NewApiClient(cmd)

			reader, err := client.Files().Download(viewName, fileName)
			helpers.ExitOnError(cmd, err, "Error downloading file")

			var writer io.Writer
			if saveAs == "-" || saveAs == "" {
				writer = cmd.OutOrStdout()
			} else {
				var err error
				writer, err = os.Create(saveAs)
				helpers.ExitOnError(cmd, err, "Error opening output file")
			}

			_, err = io.Copy(writer, reader)
			helpers.ExitOnError(cmd, err, "Error writing output")
		},
	}

	cmd.Flags().StringVar(&saveAs, "save", "", "Save to file")

	return cmd
}
