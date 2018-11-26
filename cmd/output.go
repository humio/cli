package cmd

import (
	"fmt"

	"github.com/ryanuber/columnize"
)

func printTable(rows []string) {

	table := columnize.SimpleFormat(rows)

	fmt.Println()
	fmt.Println(table)
	fmt.Println()
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
