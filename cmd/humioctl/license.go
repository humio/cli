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
	"github.com/olekukonko/tablewriter"
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

func printLicenseInfo(cmd *cobra.Command, license api.License) {

	var data [][]string

	data = append(data, []string{"License Type", license.LicenseType()})

	if onprem, ok := license.(api.OnPremLicense); ok {
		data = append(data, []string{"License ID", onprem.ID})
		data = append(data, []string{"Issued To", onprem.IssuedTo})
		data = append(data, []string{"Number Of Seats", fmt.Sprintf("%d", onprem.NumberOfSeats)})
	}

	data = append(data, []string{"Issued At", license.IssuedAt()})
	data = append(data, []string{"Expires At", license.ExpiresAt()})

	w := tablewriter.NewWriter(cmd.OutOrStdout())
	w.AppendBulk(data)
	w.SetBorder(false)
	w.SetColumnSeparator(":")
	w.SetColumnAlignment([]int{tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT})

	w.Render()
}
