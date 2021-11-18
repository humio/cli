package format

import "github.com/spf13/cobra"

func PrintDetailsTable(cmd *cobra.Command, data [][]Value) {
	FormatterFromCommand(cmd).Details(data)
}

func PrintOverviewTable(cmd *cobra.Command, header []string, data [][]Value) {
	FormatterFromCommand(cmd).Table(header, data)
}
