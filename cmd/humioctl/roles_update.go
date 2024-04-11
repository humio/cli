package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/humio/cli/api"
)

func newRolesUpdateCmd() *cobra.Command {
	var colorFlag stringPtrFlag
	var orgPermissionsFlag, sysPermissionsFlag, viewPermissionsFlag []string

	cmd := cobra.Command{
		Use:   "update [flags] <role>",
		Short: "Updates a role.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			roleName := args[0]
			client := NewApiClient(cmd)

			if viewPermissionsFlag == nil && sysPermissionsFlag == nil && orgPermissionsFlag == nil && colorFlag.value == nil {
				exitOnError(cmd, fmt.Errorf("you must specify at least one flag to update"), "Nothing specified to update")
			}

			roleId, err := client.Roles().GetRoleID(roleName)
			exitOnError(cmd, err, "Error getting role ID")
			if roleId == "" {
				exitOnError(cmd, fmt.Errorf("role %q not found", roleName), "Role not found")
			}

			roleUpdate := api.NewUpdateRoleInput(
				roleId,
				roleName,
				viewPermissionsFlag,
				sysPermissionsFlag,
				orgPermissionsFlag,
				colorFlag.value,
			)

			err = client.Roles().Update(roleUpdate)
			exitOnError(cmd, err, "Error updating role")

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully updated role %q\n", roleName)
		},
	}

	cmd.Flags().Var(&colorFlag, "color", "The color of the role in RGB hexadecimal, e.g. #FF0000.")
	cmd.Flags().StringSliceVar(&orgPermissionsFlag, "org-permissions", []string{}, "The organization permissions of the role.")
	cmd.Flags().StringSliceVar(&sysPermissionsFlag, "system-permissions", []string{}, "The system permissions of the role.")
	cmd.Flags().StringSliceVar(&viewPermissionsFlag, "view-permissions", []string{}, "(required) The view permissions of the role.")

	return &cmd
}
