package customflags

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

type StringPtrFlag struct {
	Value *string
}

func (sf *StringPtrFlag) Set(x string) error {
	sf.Value = &x
	return nil
}

func (sf *StringPtrFlag) String() string {
	if sf.Value == nil {
		return ""
	}
	return *sf.Value
}

func (sf *StringPtrFlag) Type() string {
	return "string"
}

type BoolPtrFlag struct {
	Value *bool
}

func (sf *BoolPtrFlag) Set(v string) error {
	var val bool
	if v == "true" {
		val = true
	} else if v == "false" {
		val = false
	} else {
		return errors.New("a boolean flag must be set to 'true' or 'false'")
	}
	sf.Value = &val
	return nil
}

func (sf *BoolPtrFlag) String() string {
	if sf.Value == nil {
		return ""
	}
	if *sf.Value {
		return "true"
	}
	return "false"
}

func (sf *BoolPtrFlag) Type() string {
	return "bool"
}

type UrlPtrFlag struct {
	Value *string
}

func (sf *UrlPtrFlag) Set(v string) error {
	_, err := url.Parse(v)
	if err == nil {
		sf.Value = &v
	}
	return err
}

func (sf *UrlPtrFlag) String() string {
	if sf.Value == nil {
		return ""
	}
	return *sf.Value
}

func (sf *UrlPtrFlag) Type() string {
	return "url"
}

type Float64PtrFlag struct {
	Value *float64
}

func (sf *Float64PtrFlag) Set(v string) error {
	var val float64
	val, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return err
	}
	sf.Value = &val
	return nil
}

func (sf *Float64PtrFlag) String() string {
	if sf.Value == nil {
		return ""
	}
	return fmt.Sprintf("%f", *sf.Value)
}

func (sf *Float64PtrFlag) Type() string {
	return "float64"
}
