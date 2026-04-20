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
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/humio/cli/internal/api"
	"github.com/spf13/cobra"
)

func newUsersUpdateUsernameCmd() *cobra.Command {
	var csvFile string
	var dryRun bool

	cmd := cobra.Command{
		Use:   "update-username [flags] <current-username> <new-username>",
		Short: "Updates a user's username. Requires the system permission 'ChangeUsername'. [Root Only]",
		Args: func(cmd *cobra.Command, args []string) error {
			if csvFile != "" {
				return nil
			}
			return cobra.ExactArgs(2)(cmd, args)
		},
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			if csvFile != "" {
				rows, err := ReadCSV(csvFile)
				if err != nil {
					exitOnError(cmd, err, "Error reading CSV file")
				}

				for _, row := range rows {
					updateUsername(client, row["username"], row["new_username"], dryRun)
				}
			} else {
				updateUsername(client, args[0], args[1], dryRun)
			}
		},
	}

	cmd.Flags().StringVar(&csvFile, "csv", "", "Bulk update usernames via CSV file (use headers: username, new_username)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without updating usernames")

	return &cmd
}

func updateUsername(client *api.Client, username, newUsername string, dryRun bool) {
	if dryRun {
		fmt.Printf("[DRYRUN] User %s would be renamed to %s\n", username, newUsername)
		return
	}

	_, err := client.Users().UpdateUsername(
		username,
		newUsername,
	)
	if err != nil {
		fmt.Printf("Error updating user: %s\n", err)
		return
	}

	fmt.Printf("Successfully updated user: %s renamed to %s\n", username, newUsername)
}

// Simple function to read a CSV file and return an array of maps, assumes the first row is always headers.
func ReadCSV(filePath string) ([]map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) == 0 {
		return []map[string]string{}, nil
	}

	headers := records[0]
	results := make([]map[string]string, 0, len(records)-1)

	for _, row := range records[1:] {
		if len(row) != len(headers) {
			continue // Skip rows with missing columns
		}

		rowMap := make(map[string]string)
		for i, value := range row {
			rowMap[strings.TrimSpace(headers[i])] = strings.TrimSpace(value)
		}
		results = append(results, rowMap)
	}

	return results, nil
}
