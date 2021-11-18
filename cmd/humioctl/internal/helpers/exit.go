package helpers

import (
	"github.com/spf13/cobra"
	"os"
)

func ExitOnError(cmd *cobra.Command, err error, message string) {
	if err != nil {
		cmd.Printf(message+": %s\n", err)
		os.Exit(1)
	}
}
