package humioctl

import (
	"encoding/json"
	"fmt"
	"github.com/humio/cli/api"
	format2 "github.com/humio/cli/cmd/humioctl/internal/format"
	"github.com/humio/cli/cmd/humioctl/internal/helpers"
	"github.com/spf13/cobra"
	"os"
	"sort"
	"strings"
)

type healthCheckResult struct {
	Checks        map[string]api.HealthCheck `json:"checks"`
	Version       string                     `json:"version"`
	Uptime        string                     `json:"uptime"`
	Status        api.StatusValue            `json:"status"`
	StatusMessage string                     `json:"statusMessage"`
}

func newHealthCmd() *cobra.Command {
	var (
		versionFlag    bool
		uptimeFlag     bool
		failFlag       bool
		warnAsDownFlag bool
		selectChecks   []string
	)

	cmd := &cobra.Command{
		Use:   "health",
		Short: "Health",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			health, err := client.Health()
			helpers.ExitOnError(cmd, err, "Error getting health information")

			switch {
			case versionFlag:
				fmt.Fprintf(cmd.OutOrStdout(), "%s\n", health.Version)
				return
			case uptimeFlag:
				fmt.Fprintf(cmd.OutOrStdout(), "%s\n", health.Uptime)
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

			result := healthCheckResult{
				Checks:        m,
				Version:       health.Version,
				Uptime:        health.Uptime,
				Status:        health.Status,
				StatusMessage: health.StatusMessage,
			}

			printHealthDetailsTable(cmd, result)
			printHealthOverviewTable(cmd, result)

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

	cmd.Flags().BoolVar(&versionFlag, "version", false, "Print server version and exit.")
	cmd.Flags().BoolVar(&uptimeFlag, "uptime", false, "Print uptime and exit.")
	cmd.Flags().BoolVar(&failFlag, "fail", false, "Set exit code to number of down checks.")
	cmd.Flags().BoolVar(&warnAsDownFlag, "warn-as-down", false, "When used with --fail: Treat warnings as down")
	cmd.Flags().StringSliceVarP(&selectChecks, "select", "s", nil, "Select checks to display. Specify multiple times for multiple checks.\n"+
		"If the server does not support the selected value, it will be left out.\n"+
		"Note: --select affects the checks that are considered by --fail")

	return cmd
}

func printHealthDetailsTable(cmd *cobra.Command, result healthCheckResult) {
	details := [][]format2.Value{
		{format2.String("Status"), format2.String(result.Status)},
		{format2.String("Message"), format2.String(result.StatusMessage)},
		{format2.String("Version"), format2.String(result.Version)},
		{format2.String("Uptime"), format2.String(result.Uptime)},
	}

	format2.PrintDetailsTable(cmd, details)
}

func printHealthOverviewTable(cmd *cobra.Command, result healthCheckResult) {
	var healthChecksNames []string
	for name := range result.Checks {
		healthChecksNames = append(healthChecksNames, name)
	}
	sort.Strings(healthChecksNames)

	var rows [][]format2.Value
	for _, name := range healthChecksNames {
		var keys []string
		for f := range result.Checks[name].Fields {
			keys = append(keys, f)
		}
		sort.Strings(keys)

		fields := HealthFields{}
		for _, f := range keys {
			fields[f] = format2.String(fmt.Sprint(result.Checks[name].Fields[f]))
		}
		rows = append(rows, []format2.Value{
			format2.String(result.Checks[name].Name),
			format2.String(result.Checks[name].Status),
			format2.String(result.Checks[name].StatusMessage),
			fields,
		})
	}

	format2.PrintOverviewTable(cmd, []string{"name", "status", "message", "fields"}, rows)
}

type HealthFields map[string]format2.Value

func (h HealthFields) String() string {
	var s []string
	for k, v := range h {
		s = append(s, fmt.Sprintf("%s=%q", k, v.String()))
	}

	sort.Strings(s)

	return strings.Join(s, " ")
}

func (h HealthFields) MarshalJSON() ([]byte, error) {
	return json.Marshal((map[string]format2.Value)(h))
}
