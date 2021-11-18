// Copyright © 2018 Humio Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/humio/cli/cmd/humioctl/internal/helpers"
	"github.com/humio/cli/cmd/humioctl/internal/viperkey"
	"github.com/humio/cli/cmd/internal/format"
	"github.com/humio/cli/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Shows general status information",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)
			serverStatus, err := client.Status()
			helpers.ExitOnError(cmd, err, "Error getting server status")

			username, err := client.Viewer().Username()
			helpers.ExitOnError(cmd, err, "Error getting the current user")

			details := [][]format.Value{
				{format.String("Status"), StatusText(serverStatus.Status)},
				{format.String("Address"), format.String(viper.GetString(viperkey.Address))},
				{format.String("Version"), format.String(serverStatus.Version)},
				{format.String("Username"), format.String(username)},
			}

			format.PrintDetailsTable(cmd, details)
		},
	}

	cmd.AddCommand(newLicenseInstallCmd())
	cmd.AddCommand(newLicenseShowCmd())

	return cmd
}

type StatusText string

func (s StatusText) String() string {
	switch s {
	case "OK":
		return prompt.Colorize("[green]OK[reset]")
	case "WARN":
		return prompt.Colorize("[yellow]WARN[reset]")
	default:
		return prompt.Colorize(fmt.Sprintf("[red]%s[reset]", string(s)))
	}
}

func (s StatusText) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(s))
}
