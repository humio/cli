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
	"github.com/humio/cli/cmd/humioctl/internal/customflags"
	"github.com/humio/cli/cmd/humioctl/internal/helpers"
	"github.com/spf13/cobra"
)

func newUsersAddCmd() *cobra.Command {
	var rootFlag customflags.BoolPtrFlag
	var nameFlag, companyFlag, emailFlag, countryCodeFlag customflags.StringPtrFlag
	var pictureFlag customflags.UrlPtrFlag

	cmd := cobra.Command{
		Use:   "add [flags] <username>",
		Short: "Adds a user. [Root Only]",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]
			client := NewApiClient(cmd)

			_, err := client.Users().Add(username, api.UserChangeSet{
				IsRoot:      rootFlag.Value,
				FullName:    nameFlag.Value,
				Company:     companyFlag.Value,
				CountryCode: countryCodeFlag.Value,
				Email:       emailFlag.Value,
				Picture:     pictureFlag.Value,
			})
			helpers.ExitOnError(cmd, err, "Error creating the user")

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully created user with username %q\n", username)
		},
	}

	cmd.Flags().Var(&rootFlag, "root", "If true grants root access to the user.")
	cmd.Flags().Var(&nameFlag, "name", "The full name of the user.")
	cmd.Flags().Var(&countryCodeFlag, "country-code", "A two letter country code.")
	cmd.Flags().Var(&companyFlag, "company", "The company where the user works.")
	cmd.Flags().Var(&pictureFlag, "picture", "A URL to an avatar for user.")
	cmd.Flags().Var(&emailFlag, "email", "The user's email. Does not change the username if the email is used.")

	return &cmd
}
