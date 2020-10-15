package main

import (
	"github.com/spf13/cobra"
	"os"
)

type humioResultType interface{}

type humioErrorExitCode struct {
	err      error
	exitCode int
}

func (h humioErrorExitCode) Error() string {
	return h.err.Error()
}

func (h humioErrorExitCode) Unwrap() error {
	return h.err
}

func WrapRun(f func(cmd *cobra.Command, args []string) (humioResultType, error)) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		res, err := f(cmd, args)
		if err != nil {
			switch err := err.(type) {
			case humioErrorExitCode:
				cmd.Printf("Error: %s", err.err)
				os.Exit(err.exitCode)
			default:
				cmd.Printf("Error: %s", err)
				os.Exit(1)
			}
		}

		if res != nil {
			cmd.Println(res)
		}
	}
}
