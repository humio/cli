package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func exitOnError(cmd *cobra.Command, err error, message string) {
	if err != nil {
		cmd.Println(fmt.Errorf(message+": %s", err))
		os.Exit(1)
	}
}

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

func (sf *stringPtrFlag) Type() string {
	return "string"
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

type urlPtrFlag struct {
	value *string
}

func (sf *urlPtrFlag) Set(v string) error {
	_, err := url.Parse(v)
	if err == nil {
		sf.value = &v
	}
	return err
}

func (sf *urlPtrFlag) String() string {
	if sf.value == nil {
		return ""
	}
	return *sf.value
}

func (sf *urlPtrFlag) Type() string {
	return "url"
}
