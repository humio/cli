package main

import (
	"strings"

	"github.com/humio/cli/api"
	"github.com/humio/cli/cmd/internal/format"
	"github.com/spf13/cobra"
)

func newRolesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "roles",
		Short: "Manage roles",
	}

	cmd.AddCommand(newRolesShowCmd())
	cmd.AddCommand(newRolesListCmd())
	cmd.AddCommand(newRolesCreateCmd())
	cmd.AddCommand(newRolesUpdateCmd())
	cmd.AddCommand(newRolesDeleteCmd())

	return cmd
}

func printRoleDetailsTable(cmd *cobra.Command, role api.Role) {
	details := [][]format.Value{
		{format.String("ID"), format.String(role.ID)},
		{format.String("Name"), format.String(role.DisplayName)},
		{format.String("Description"), format.String(role.Description)},
		{format.String("View Permissions"), format.String(strings.Join(role.ViewPermissions, "\n"))},
		{format.String("System Permissions"), format.String(strings.Join(role.SystemPermissions, "\n"))},
		{format.String("Organization Permissions"), format.String(strings.Join(role.OrgPermissions, "\n"))},
	}

	printDetailsTable(cmd, details)
}
