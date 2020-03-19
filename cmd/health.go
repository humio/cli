package cmd

import (
	"encoding/json"
	"github.com/humio/cli/api"
	"github.com/spf13/cobra"
	"os"
)

func newHealthCmd() *cobra.Command {
	var (
		textFlag       bool
		rawFlag       bool
		versionFlag    bool
		uptimeFlag     bool
		failFlag       bool
		warnAsDownFlag bool
		selectChecks   []string
	)

	cmd := &cobra.Command{
		Use:   "health",
		Short: "Health",
		Args:  cobra.ExactArgs(0),

		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			if textFlag {
				healthStr, err := client.HealthString()

				exitOnError(cmd, err, "error getting health information")

				cmd.Println(healthStr)

				return
			}

			health, err := client.Health()
			exitOnError(cmd, err, "error getting health information")

			if rawFlag {
				_, _ = cmd.OutOrStdout().Write(health.Json())
				return
			}

			if versionFlag {
				cmd.Println(health.Version)
				return
			}

			if uptimeFlag {
				cmd.Println(health.Uptime)
				return
			}

			m := health.ChecksMap()
			if len(selectChecks) > 0 {
				newMap := map[string]api.HealthCheck{}
				for _, s := range selectChecks {
					if c, ok := m[s]; ok {
						newMap[s] = c
					}
				}
				m = newMap
			}

			_ = json.NewEncoder(cmd.OutOrStdout()).Encode(m)

			if failFlag {
				numDown := 0
				for _, c := range m {
					if c.Status == api.StatusDown || (warnAsDownFlag && c.Status == api.StatusWarn) {
						numDown++
					}
				}

				os.Exit(numDown)
			}
		},
	}

	cmd.Flags().BoolVar(&textFlag, "text", false, "Print text health information and exit.")
	cmd.Flags().BoolVar(&rawFlag, "raw", false, "Print json health information as reported by server and exit.")
	cmd.Flags().BoolVar(&versionFlag, "version", false, "Print server version and exit.")
	cmd.Flags().BoolVar(&uptimeFlag, "uptime", false, "Print uptime and exit.")
	cmd.Flags().BoolVar(&failFlag, "fail", false, "Set exit code to number of down checks.")
	cmd.Flags().BoolVar(&warnAsDownFlag, "warn-as-down", false, "When used with --fail: Treat warnings as down")
	cmd.Flags().StringSliceVarP(&selectChecks, "select", "s", nil, "Select checks to display. Specify multiple times for multiple checks.\n" +
		"If the server does not support the selected value, it will be left out.\n" +
		"Note: --select affects the checks that are considered by --fail")

	return cmd
}
