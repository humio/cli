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

	"github.com/humio/cli/api"
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
	var details [][]string

	if onprem, ok := license.(api.OnPremLicense); ok {
		details = append(details, []string{"License ID", onprem.ID})
		details = append(details, []string{"Issued To", onprem.IssuedTo})
		details = append(details, []string{"Number Of Seats", fmt.Sprintf("%d", onprem.NumberOfSeats)})
	}

	details = append(details, []string{"Issued At", license.IssuedAt()})
	details = append(details, []string{"Expires At", license.ExpiresAt()})

	printDetailsTable(cmd, details)
}
