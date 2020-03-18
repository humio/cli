package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/humio/cli/api"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func newHealthCmd() *cobra.Command {
	var (
		jsonFlag       bool
		versionFlag    bool
		uptimeFlag     bool
		failFlag       bool
		warnAsDownFlag bool
		selectChecks   []string
	)

	cmd := &cobra.Command{
		Use:   "health",
		Short: "Health",
		Long: "Query health checks\n" +
			"\n" +
			"Only one of the following flags can be used at a time: [--fail, --json, --select|-s, --uptime, --version]\n" +
			"--select|-s may of course be specified multiple times.",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			if !mutuallyExclusive(jsonFlag, versionFlag, uptimeFlag, failFlag, len(selectChecks) > 0) {
				cmd.Println("some flags are mutually exclusive")
				os.Exit(1)
			}

			var health api.Health
			if jsonFlag || versionFlag || uptimeFlag || failFlag || len(selectChecks) > 0 {
				var err error
				health, err = client.Health()

				exitOnError(cmd, err, "error getting health information")

				switch {
				case jsonFlag:
					_ = json.NewEncoder(cmd.OutOrStdout()).Encode(health)
				case uptimeFlag:
					cmd.Println(health.Uptime)
				case versionFlag:
					cmd.Println(health.Version)
				case failFlag:
					if warnAsDownFlag {
						os.Exit(len(health.Down) + len(health.Warn))
					} else {
						os.Exit(len(health.Down))
					}
				case len(selectChecks) > 0:
					m := health.ChecksMap()
					for _, s := range selectChecks {
						c, has := m[s]

						if !has {
							cmd.Printf("%s=UNKNOWN\n", s)
						} else {
							cmd.Println(healthCheckToString(c))
						}
					}
				}

				return
			}

			healthStr, err := client.HealthString()

			exitOnError(cmd, err, "error getting health information")

			cmd.Println(healthStr)
		},
	}

	cmd.Flags().BoolVarP(&jsonFlag, "json", "j", false, "Output json.")
	cmd.Flags().BoolVar(&versionFlag, "version", false, "Print version and exit.")
	cmd.Flags().BoolVar(&uptimeFlag, "uptime", false, "Print uptime and exit.")
	cmd.Flags().BoolVar(&failFlag, "fail", false, "Set exit code to number of down checks.")
	cmd.Flags().BoolVar(&warnAsDownFlag, "warn-as-down", false, "When used with --fail: Treat warnings as down")
	cmd.Flags().StringSliceVarP(&selectChecks, "select", "s", nil, "Select checks to display. Specify multiple times for multiple checks.")

	return cmd
}

func mutuallyExclusive(bools ...bool) bool {
	foundTrue := false

	for _, b := range bools {
		if b && foundTrue {
			return false
		}
		foundTrue = foundTrue || b
	}

	return true
}

func healthCheckToString(healthCheck api.HealthCheck) string {
	s := fmt.Sprintf("%s=%s", healthCheck.Name, healthCheck.Status)
	if len(healthCheck.StatusMessage) > 0 {
		s = s + fmt.Sprintf(" (%s)", healthCheck.StatusMessage)
	}
	var fields []string
	for k, v := range healthCheck.Fields {
		fields = append(fields, fmt.Sprintf("%s=%q", k, v))
	}

	if len(fields) > 0 {
		s = s + fmt.Sprintf(" [%s]", strings.Join(fields, " "))
	}

	return s
}

