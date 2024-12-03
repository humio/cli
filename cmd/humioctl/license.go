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
	"github.com/humio/cli/internal/api"
	"github.com/humio/cli/internal/format"
	"github.com/spf13/cobra"
)

func newLicenseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "license",
		Short: "Manage the Humio license [Root Only]",
	}

	cmd.AddCommand(newLicenseInstallCmd())
	cmd.AddCommand(newLicenseShowCmd())

	return cmd
}

func printLicenseDetailsTable(cmd *cobra.Command, license api.License) {
	var details [][]format.Value

	if onprem, ok := license.(api.OnPremLicense); ok {
		details = append(details, []format.Value{format.String("License ID"), format.String(onprem.ID)})
		details = append(details, []format.Value{format.String("Issued To"), format.String(onprem.IssuedTo)})
		if onprem.NumberOfSeats != nil {
			details = append(details, []format.Value{format.String("Number Of Seats"), format.Int(*onprem.NumberOfSeats)})
		}
	}

	details = append(details, []format.Value{format.String("Issued At"), format.String(license.IssuedAt())})
	details = append(details, []format.Value{format.String("Expires At"), format.String(license.ExpiresAt())})

	printDetailsTable(cmd, details)
}
