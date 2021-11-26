package main

import (
	"fmt"
	"github.com/humio/cli/api"
	"github.com/spf13/cobra"
	"os"
)

func newFeatureFlagsEnableCmd() *cobra.Command {
	var (
		organizationID, userID string
		global                 bool
	)

	cmd := &cobra.Command{
		Use:   "enable [--global | --user <user-id> | --organization <organization-id>] <feature-flag>",
		Short: "Enable a feature flag",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			enableDisableFeatureFlag(cmd, args, global, organizationID, userID, true)
		},
	}

	cmd.Flags().StringVar(&organizationID, "organization", "", "enable for an organization")
	cmd.Flags().StringVar(&userID, "user", "", "enable for a user")
	cmd.Flags().BoolVar(&global, "global", false, "enable globally")

	return cmd
}

func newFeatureFlagsDisableCmd() *cobra.Command {
	var (
		organizationID, userID string
		global                 bool
	)

	cmd := &cobra.Command{
		Use:   "disable [--global | --user <user-id> | --organization <organization-id>] <feature-flag>",
		Short: "disable a feature flag",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			enableDisableFeatureFlag(cmd, args, global, organizationID, userID, false)
		},
	}

	cmd.Flags().StringVar(&organizationID, "organization", "", "disable for an organization")
	cmd.Flags().StringVar(&userID, "user", "", "disable for a user")
	cmd.Flags().BoolVar(&global, "global", false, "disable globally")

	return cmd
}

func enableDisableFeatureFlag(cmd *cobra.Command, args []string, global bool, organizationID string, userID string, enable bool) {
	action := "disabling"
	postTense := "Disabled"
	if enable {
		action = "enabling"
		postTense = "Enabled"
	}

	flag := api.FeatureFlag(args[0])

	if global && len(organizationID) > 0 && len(userID) > 0 {
		cmd.PrintErrln("cannot specify --global, --user and --organization at the same time")
		os.Exit(1)
	}

	if global && len(organizationID) > 0 {
		cmd.PrintErrln("cannot specify --global with --organization")
		os.Exit(1)
	}

	if global && len(userID) > 0 {
		cmd.PrintErrln("cannot specify --global with --user")
		os.Exit(1)
	}

	if len(organizationID) > 0 && len(userID) > 0 {
		cmd.PrintErrln("cannot specify --user with --organization")
		os.Exit(1)
	}

	if len(organizationID) == 0 && len(userID) == 0 && !global {
		cmd.PrintErrln("must specify one of --global, --user or --organization")
		os.Exit(1)
	}

	client := NewApiClient(cmd)

	flags, err := client.FeatureFlags().SupportedFlags()
	exitOnError(cmd, err, "error fetching supported feature flags")

	var foundFlag bool
	for _, f := range flags {
		if f == flag {
			foundFlag = true
			break
		}
	}

	if !foundFlag {
		cmd.PrintErrf("unsupported feature flag %q\n\nSupported ones are:\n", flag)
		for _, f := range flags {
			cmd.PrintErrf("  - %s\n", f)
		}
		os.Exit(1)
	}

	var infoSelector string

	switch {
	case global:
		infoSelector = "globally"
		if enable {
			err = client.FeatureFlags().EnableGlobally(flag)
		} else {
			err = client.FeatureFlags().DisableGlobally(flag)
		}
	case len(organizationID) > 0:
		infoSelector = fmt.Sprintf("for organization %q", organizationID)
		if enable {
			err = client.FeatureFlags().EnableForOrganization(organizationID, flag)
		} else {
			err = client.FeatureFlags().DisableForOrganization(organizationID, flag)
		}
	case len(userID) > 0:
		infoSelector = fmt.Sprintf("for user %q", userID)
		if enable {
			err = client.FeatureFlags().EnableForUser(userID, flag)
		} else {
			err = client.FeatureFlags().DisableForUser(userID, flag)
		}
	}

	exitOnError(cmd, err, "error "+action+" feature flag")

	fmt.Fprintf(cmd.OutOrStdout(), "%s feature %q %s\n", postTense, flag, infoSelector)
}
