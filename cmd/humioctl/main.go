package main

import (
	"github.com/humio/cli/cmd/humioctl/humioctl"
	"github.com/humio/cli/cmd/humioctl/internal/helpers"
)

func main() {
	err := humioctl.RootCmd.Execute()
	helpers.ExitOnError(humioctl.RootCmd, err, "Unable to execute command")
}
