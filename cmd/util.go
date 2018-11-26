package cmd

import (
	"errors"
	"log"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//
// func Map(vs []string, f func(string) testCase) []testCase {
// 	vsm := make([]testCase, len(vs))
// 	for i, v := range vs {
// 		vsm[i] = f(v)
// 	}
// 	return vsm
// }

// commandErrorText is used to easily render the same messaging across commads
// when an error is printed.
// func commandErrorText(cmd NamedCommand) string {
// 	return fmt.Sprintf("For additional help try 'humio %s -help'", cmd.Name())
// }

type stringPtrFlag struct {
	value *string
}

func (sf *stringPtrFlag) Set(x string) error {
	sf.value = &x
	return nil
}

func (sf *stringPtrFlag) String() string {
	if sf.value == nil {
		return ""
	}
	return *sf.value
}

type boolPtrFlag struct {
	value *bool
}

func (sf *boolPtrFlag) Set(v string) error {
	var val bool
	if v == "true" {
		val = true
	} else if v == "false" {
		val = false
	} else {
		return errors.New("a boolean flag must be set to 'true' or 'false'")
	}
	sf.value = &val
	return nil
}

func (sf *boolPtrFlag) String() string {
	if sf.value == nil {
		return ""
	}
	if *sf.value {
		return "true"
	}
	return "false"
}

func (sf *boolPtrFlag) Type() string {
	return "bool"
}
