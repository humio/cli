package main

import (
	"github.com/spf13/cobra"
	"io"
	"os"
)

func newFilesDownloadCmd() *cobra.Command {
	var (
		saveAs string
	)

	cmd := &cobra.Command{
		Use:  "download <view name> <file name>",
		Long: `Download a file.`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			reader, err := client.Files().Download(args[0], args[1])
			exitOnError(cmd, err, "error downloading file")

			var writer io.Writer
			if saveAs == "-" || saveAs == "" {
				writer = cmd.OutOrStdout()
			} else {
				var err error
				writer, err = os.Create(saveAs)
				exitOnError(cmd, err, "error opening output file")
			}

			_, err = io.Copy(writer, reader)
			exitOnError(cmd, err, "error writing output")
		},
	}

	cmd.Flags().StringVar(&saveAs, "save", "", "Save to file")

	return cmd
}

