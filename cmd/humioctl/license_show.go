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
	"errors"

	"github.com/humio/cli/internal/api"
	"github.com/spf13/cobra"
)

func newLicenseShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show the current Humio license installed",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)
			license, err := client.Licenses().Get()
			noLicense := api.OnPremLicense{}
			if license == noLicense {
				err = errors.New("no license currently installed")
			}
			exitOnError(cmd, err, "Error fetching the license")

			printLicenseDetailsTable(cmd, license)
		},
	}

	return cmd
}
