package humioctl

import (
	format2 "github.com/humio/cli/cmd/humioctl/internal/format"
	"github.com/humio/cli/cmd/humioctl/internal/helpers"
	"github.com/spf13/cobra"
)

func newGroupsList() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List groups",
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			groups, err := client.Groups().List()
			helpers.ExitOnError(cmd, err, "Error listing groups")

			rows := make([][]format2.Value, len(groups))
			for i, group := range groups {
				rows[i] = []format2.Value{
					format2.String(group.DisplayName),
					format2.String(group.ID),
				}
			}

			format2.PrintOverviewTable(cmd, []string{"Display Name", "ID"}, rows)
		},
	}
}
