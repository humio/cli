package main

import (
	"github.com/humio/cli/internal/format"
	"github.com/spf13/cobra"
)

func newFeatureFlagsSupportedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "supported",
		Short: "List supported feature flags.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			flags, err := client.FeatureFlags().SupportedFlags()
			exitOnError(cmd, err, "Error listing feature flags")

			var rows [][]format.Value
			for _, flag := range flags {
				rows = append(rows, []format.Value{
					format.String(flag.Flag),
					format.Bool(flag.Experimental),
					format.String(flag.Description),
				})
			}

			printOverviewTable(cmd, []string{"Feature Flag", "Experimental", "Description"}, rows)
		},
	}

	return cmd
}
