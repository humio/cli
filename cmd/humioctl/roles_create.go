package main

import (
	"fmt"

	"github.com/humio/cli/api"
	"github.com/spf13/cobra"
)

func newRolesCreateCmd() *cobra.Command {
	var colorFlag stringPtrFlag
	var orgPermissionsFlag, sysPermissionsFlag, viewPermissionsFlag []string

	cmd := cobra.Command{
		Use:   "create <role>",
		Short: "Create a role.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			role := api.NewAddRoleInput(
				args[0],
				colorFlag.value,
				viewPermissionsFlag,
				sysPermissionsFlag,
				orgPermissionsFlag,
			)

			err := client.Roles().Create(role)
			exitOnError(cmd, err, "Error creating role")

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully created role %s\n", args[0])

		},
	}

	cmd.Flags().Var(&colorFlag, "color", "The color of the role in RGB hexadecimal, e.g. #FF0000.")
	cmd.Flags().StringSliceVar(&orgPermissionsFlag, "org-permissions", []string{}, "The organization permissions of the role.")
	cmd.Flags().StringSliceVar(&sysPermissionsFlag, "system-permissions", []string{}, "The system permissions of the role.")
	cmd.Flags().StringSliceVar(&viewPermissionsFlag, "view-permissions", []string{}, "(required) The view permissions of the role.")

	cmd.MarkFlagRequired("view-permissions")

	return &cmd
}
