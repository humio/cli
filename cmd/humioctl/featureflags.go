package main

import "github.com/spf13/cobra"

func newFeatureFlagsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feature-flags",
		Short: "Manage feature flags",
	}

	cmd.AddCommand(newFeatureFlagsSupportedCmd())
	cmd.AddCommand(newFeatureFlagsEnableCmd())
	cmd.AddCommand(newFeatureFlagsDisableCmd())

	return cmd
}
