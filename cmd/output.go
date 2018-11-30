package cmd

import (
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
)

func printTable(cmd *cobra.Command, rows []string) {

	table := columnize.SimpleFormat(rows)

	cmd.Println()
	cmd.Println(table)
	cmd.Println()
}

func yesNo(isTrue bool) string {
	if isTrue {
		return "yes"
	}
	return "no"
}

func checkmark(v bool) string {
	if v {
		return "✓"
	}
	return ""
}

func valueOrEmpty(v string) string {
	if v == "" {
		return "-"
	}
	return v
}
