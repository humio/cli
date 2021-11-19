package humioctl

import (
	format2 "github.com/humio/cli/cmd/humioctl/internal/format"
	"github.com/humio/cli/cmd/humioctl/internal/helpers"
	"github.com/spf13/cobra"
)

func newFilesListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list <view-name>",
		Short: "List uploaded files in a view.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			files, err := client.Files().List(args[0])
			helpers.ExitOnError(cmd, err, "Error listing files")

			var rows [][]format2.Value

			for _, file := range files {
				rows = append(rows, []format2.Value{
					format2.String(file.Name),
					format2.String(file.ContentHash),
					format2.String(file.ID),
				})
			}

			format2.PrintOverviewTable(cmd, []string{"Name", "Content Hash", "ID"}, rows)
		},
	}

	return cmd
}
