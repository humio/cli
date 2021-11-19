package spec

import (
	"bytes"
	"github.com/spf13/cobra"
	"os/exec"
)

type HumioctlCommand struct {
	Name []string
	Cmd  *cobra.Command
}

type HumioctlExecutable struct {
	Path        string
	GlobalFlags []string
}

func (e HumioctlExecutable) Run(name []string, args []string) ([]byte, []byte, int, error) {
	var cmdArgs []string
	cmdArgs = append(cmdArgs, e.GlobalFlags...)
	cmdArgs = append(cmdArgs, name...)
	cmdArgs = append(cmdArgs, args...)

	process := exec.Command(e.Path, cmdArgs...)

	var stdoutBuffer, stderrBuffer bytes.Buffer

	process.Stdout = &stdoutBuffer
	process.Stderr = &stderrBuffer

	err := process.Run()
	var exitCode int
	if err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			exitCode = err.ExitCode()
			err = nil
		} else {
			return nil, nil, 0, err
		}
	}

	return stdoutBuffer.Bytes(), stderrBuffer.Bytes(), exitCode, nil
}

func FindCommands(cmd *cobra.Command, path []string) []HumioctlCommand {
	var res []HumioctlCommand

	for _, c := range cmd.Commands() {
		var p []string
		p = append(p, path...)

		if c.Runnable() {
			var use []string
			use = append(use, p...)
			use = append(use, c.Name())

			res = append(res, HumioctlCommand{
				Name: use,
				Cmd:  c,
			})
		}

		p = append(p, c.Name())

		res = append(res, FindCommands(c, p)...)
	}

	return res
}
