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
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newLicenseInstallCmd() *cobra.Command {
	var license string

	cmd := &cobra.Command{
		Use:   "install [flags] (<license-file> | --license=<string>)",
		Short: "Install a Humio license",
		Args:  cobra.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			var licenseString string
			if len(args) == 1 {
				filepath := args[0]

				// #nosec G304
				licenseBytes, err := os.ReadFile(filepath)
				exitOnError(cmd, err, "Error reading license file")

				licenseString = string(licenseBytes)
			} else if license != "" {
				// License set from flag
				licenseString = license
			} else {
				cmd.Println("Expected either an argument <filename> or flag --license=<license>.")
				_ = cmd.Help()
				os.Exit(1)
			}

			err := client.Licenses().Install(licenseString)
			exitOnError(cmd, err, "Error installing license")

			fmt.Fprintln(cmd.OutOrStdout(), "Successfully installed license")
		},
	}

	cmd.Flags().StringVarP(&license, "license", "l", "", "A string with the content license license file.")

	return cmd
}
