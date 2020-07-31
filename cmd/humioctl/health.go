package main

import (
	"encoding/json"
	"fmt"
	"github.com/humio/cli/api"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"io"
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
		Args:  cobra.ExactArgs(0),

		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			health, err := client.Health()
			exitOnError(cmd, err, "error getting health information")

			switch {
			case versionFlag:
				cmd.Println(health.Version)
				return
			case uptimeFlag:
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

			result := healthCheckResult{
				Checks:        m,
				Version:       health.Version,
				Uptime:        health.Uptime,
				Status:        health.Status,
				StatusMessage: health.StatusMessage,
			}

			if jsonFlag {
				_ = json.NewEncoder(cmd.OutOrStdout()).Encode(result)
			} else {
				encodeAsText(cmd.OutOrStdout(), result)
			}

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

	cmd.Flags().BoolVarP(&jsonFlag, "json", "j", false, "Output as json.")
	cmd.Flags().BoolVar(&versionFlag, "version", false, "Print server version and exit.")
	cmd.Flags().BoolVar(&uptimeFlag, "uptime", false, "Print uptime and exit.")
	cmd.Flags().BoolVar(&failFlag, "fail", false, "Set exit code to number of down checks.")
	cmd.Flags().BoolVar(&warnAsDownFlag, "warn-as-down", false, "When used with --fail: Treat warnings as down")
	cmd.Flags().StringSliceVarP(&selectChecks, "select", "s", nil, "Select checks to display. Specify multiple times for multiple checks.\n"+
		"If the server does not support the selected value, it will be left out.\n"+
		"Note: --select affects the checks that are considered by --fail")

	return cmd
}

func encodeAsText(writer io.Writer, result healthCheckResult) {
	tw := tablewriter.NewWriter(writer)
	tw.SetAutoWrapText(false)
	tw.Append([]string{"STATUS", string(result.Status)})
	tw.Append([]string{"MESSAGE", result.StatusMessage})
	tw.Append([]string{"VERSION", result.Version})
	tw.Append([]string{"UPTIME", result.Uptime})
	tw.Render()

	tw = tablewriter.NewWriter(writer)
	tw.SetAutoWrapText(false)
	tw.SetHeader([]string{"name", "status", "message", "fields"})

	for _, c := range result.Checks {
		var keys []string
		for f := range c.Fields {
			keys = append(keys, f)
		}
		sort.Strings(keys)

		var fields []string

		for _, f := range keys {
			fields = append(fields, fmt.Sprintf("%s=%q", f, c.Fields[f]))
		}

		row := []string{c.Name, string(c.Status), c.StatusMessage, strings.Join(fields, " ")}

		tw.Append(row)
	}

	tw.Render()
}
