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

package cmd

import (
	"fmt"
	"os"

	"github.com/humio/cli/api"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// usersCmd represents the users command
func newLicenseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "license",
		Short: "Manage the Humio license [Root Only]",
	}

	cmd.AddCommand(newLicenseInstallCmd())
	cmd.AddCommand(newLicenseShowCmd())

	return cmd
}

func printLicenseInfo(license api.LicenseData) {

	data := [][]string{
		[]string{"Issued At", license.IssuedAt},
		[]string{"Expires At", license.ExpiresAt},
	}

	w := tablewriter.NewWriter(os.Stdout)
	w.AppendBulk(data)
	w.SetBorder(false)
	w.SetColumnSeparator(":")
	w.SetColumnAlignment([]int{tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT})

	fmt.Println()
	w.Render()
	fmt.Println()
}
