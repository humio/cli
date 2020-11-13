package main

import (
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
)

func newFilesUploadCmd() *cobra.Command {
	var (
		saveAsFileName string
	)

	cmd := &cobra.Command{
		Use: "upload <view name> <input file>",
		Long: `Upload a file to a view.

Specify '-' as the input file to read from stdin.`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if args[1] == "-" && saveAsFileName == "" {
				cmd.PrintErr("When the input file is stdin, the file name must be provided with --name.\n")
				os.Exit(1)
			}

			var fileName string
			if saveAsFileName != "" {
				fileName = saveAsFileName
			} else {
				fileName = filepath.Base(args[1])
			}

			var reader io.Reader
			if args[1] == "-" {
				reader = cmd.InOrStdin()
			} else {
				var err error
				reader, err = os.Open(args[1])
				exitOnError(cmd, err, "error opening input file")
			}

			client := NewApiClient(cmd)

			err := client.Files().Upload(args[0], fileName, reader)
			exitOnError(cmd, err, "error uploading file")
		},
	}

	cmd.Flags().StringVar(&saveAsFileName, "name", "", "Specify the name of the uploaded file. Required when reading from stdin.")

	return cmd
}

