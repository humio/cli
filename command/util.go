package command

import (
	"fmt"
	"log"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Map(vs []string, f func(string) testCase) []testCase {
	vsm := make([]testCase, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

// commandErrorText is used to easily render the same messaging across commads
// when an error is printed.
func commandErrorText(cmd NamedCommand) string {
	return fmt.Sprintf("For additional help try 'humio %s -help'", cmd.Name())
}
