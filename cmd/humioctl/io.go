package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
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
				cmd.PrintErrf("Error: %s", err.err)
				os.Exit(err.exitCode)
			default:
				cmd.PrintErrf("Error: %s", err)
				os.Exit(1)
			}
		}
		if res != nil {
			if jsonOutput {
				encoder := json.NewEncoder(cmd.OutOrStdout())
				encoder.SetIndent("", "  ")
				err := encoder.Encode(res)
				if err != nil {
					cmd.PrintErrf("Error: %s", err)
					os.Exit(2)
				}
			} else {
				cmd.Println(FormatResult(res, true))
			}
		}
	}
}

func FormatResult(result interface{}, root bool) string {
	switch v := result.(type) {
	case bool:
		return strconv.FormatBool(v)
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case uintptr:
		return strconv.FormatUint(uint64(v), 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'e', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'e', -1, 64)
	case string:
		return v
	case []string:
		return strings.Join(v, ", ")
	case fmt.Stringer:
		return v.String()
	default:
		val := reflect.ValueOf(v)
		typ := val.Type()

		if typ.Kind() == reflect.Ptr {
			typ, val = typ.Elem(), val.Elem()
		}

		switch typ.Kind() {
		case reflect.Struct:
			var buf bytes.Buffer
			tw := tablewriter.NewWriter(&buf)

			tw.SetColumnAlignment([]int{tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT})
			tw.SetBorder(false)
			tw.SetColumnSeparator(":")
			tw.SetAutoWrapText(false)

			var data [][]string

			for i := 0; i < typ.NumField(); i++ {
				f := typ.Field(i)
				data = append(data, []string{
					camelCaseToWords(f.Name),
					FormatResult(val.Field(i).Interface(), false),
				})
			}

			tw.AppendBulk(data)
			tw.Render()

			return buf.String()
		case reflect.Map:
			str, _ := json.Marshal(v)
			return string(str)
		case reflect.Slice:
			if val.Len() == 0 {
				return ""
			}

			var buf bytes.Buffer
			tw := tablewriter.NewWriter(&buf)
			tw.SetBorder(false)

			var fields []string

			sliceTyp := typ.Elem()
			switch sliceTyp.Kind() {
			case reflect.Struct:
				for i := 0; i < sliceTyp.NumField(); i++ {
					fields = append(fields, camelCaseToWords(sliceTyp.Field(i).Name))
				}

				tw.SetHeader(fields)

				for i := 0; i < val.Len(); i++ {
					var row []string
					for j := range fields {
						row = append(row, FormatResult(val.Index(i).Field(j).Interface(), false))
					}
					tw.Append(row)
				}
			default:
				for i := 0; i < val.Len(); i++ {
					tw.Append([]string{FormatResult(val.Index(i).Interface(), false)})
				}
			}

			tw.Render()

			return buf.String()
		}
	}

	return "<unknown>"
}

func camelCaseToWords(s string) string {
	return regexp.MustCompile(`([a-z])([A-Z])`).ReplaceAllString(s, "$1 $2")
}
