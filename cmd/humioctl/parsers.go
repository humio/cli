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
	"github.com/spf13/cobra"
)

func newParsersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parsers",
		Short: "Manage parsers",
	}

	cmd.AddCommand(newParsersInstallCmd())
	cmd.AddCommand(newParsersListCmd())
	cmd.AddCommand(newParsersRemoveCmd())
	cmd.AddCommand(newParsersExportCmd())

	return cmd
}
