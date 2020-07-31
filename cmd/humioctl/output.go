package main

import (
	"fmt"

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

// ByteCountDecimal returns a human-readable size of a byte count
func ByteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}
