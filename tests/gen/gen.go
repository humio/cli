package main

import (
	"fmt"
	"github.com/humio/cli/cmd/humioctl/humioctl"
	"github.com/humio/cli/tests/spec"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

var executable = spec.HumioctlExecutable{
	Path: "/Users/dsoerensen/src/cli/bin/humioctl",
	GlobalFlags: []string{
		"--format", "json",
		"--address", "http://127.0.0.1:8080",
		"--token", "",
	},
}

var skip = [][]string{
	{"ingest"},
	{"profiles"},
	{"welcome"},
}

func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func isSkipped(name []string) bool {
	for _, s := range skip {
		if stringSliceEqual(s, name) {
			return true
		}
	}
	return false
}

func nameToID(name []string) string {
	return strings.Join(name, "_")
}

func main() {
	commands := spec.FindCommands(humioctl.RootCmd, nil)

	var commandSpecs []spec.CommandSpec

	for _, cmd := range commands {
		if isSkipped(cmd.Name) {
			continue
		}
		if cmd.Cmd.Args != nil && cmd.Cmd.Args(cmd.Cmd, nil) == nil {
			stdout, stderr, exitCode, err := executable.Run(cmd.Name, nil)
			if err != nil {
				panic(err)
			}

			if exitCode != 0 {
				fmt.Fprintf(os.Stderr, "Command %q return exit status %d, skipping auto generation.\n", cmd.Name, exitCode)
				continue
			}

			stdoutSpec, err := spec.GenerateOutputStreamSpec(stdout)
			if err != nil {
				panic(err)
			}
			stderrSpec, err := spec.GenerateOutputStreamSpec(stderr)
			if err != nil {
				panic(err)
			}

			commandSpecs = append(commandSpecs, spec.CommandSpec{
				ID:   nameToID(cmd.Name),
				Name: cmd.Name,
				InputOutput: []spec.CommandInputOutputSpec{
					{
						ID: "default",
						Output: spec.CommandOutputSpec{
							ExitCode: exitCode,
							Stdout:   stdoutSpec,
							Stderr:   stderrSpec,
						},
					},
				},
			})
		}
	}

	yaml.NewEncoder(os.Stdout).Encode(spec.RootSpec{Commands: commandSpecs})
}
