package humioctl

import "github.com/spf13/cobra"

func newGroupsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "groups",
		Short: "Manage groups",
	}

	cmd.AddCommand(newGroupsAddUserCmd())
	cmd.AddCommand(newGroupsRemoveUserCmd())
	cmd.AddCommand(newGroupsList())

	return cmd
}
