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

package humioctl

import (
	"github.com/humio/cli/api"
	format2 "github.com/humio/cli/cmd/humioctl/internal/format"
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
	var details [][]format2.Value

	if onprem, ok := license.(api.OnPremLicense); ok {
		details = append(details, []format2.Value{format2.String("License ID"), format2.String(onprem.ID)})
		details = append(details, []format2.Value{format2.String("Issued To"), format2.String(onprem.IssuedTo)})
		details = append(details, []format2.Value{format2.String("Number Of Seats"), format2.Int(onprem.NumberOfSeats)})
	}

	details = append(details, []format2.Value{format2.String("Issued At"), format2.String(license.IssuedAt())})
	details = append(details, []format2.Value{format2.String("Expires At"), format2.String(license.ExpiresAt())})

	format2.PrintDetailsTable(cmd, details)
}
