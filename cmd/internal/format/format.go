package format

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func FormatterFromCommand(cmd *cobra.Command) Formatter {
	factories := map[string]func(cmd *cobra.Command) Formatter{
		"json": func(cmd *cobra.Command) Formatter { return JSONFormatter{cmd.OutOrStdout()} },
	}

	var formatter string
	if value := cmd.Flags().Lookup("format").Value; value != nil {
		formatter = value.String()
	}
	factory := factories[formatter]

	if factory == nil {
		return TextFormatter{w: cmd.OutOrStdout()}
	}

	return factory(cmd)
}

type Formatter interface {
	Details(rows [][]Value)
	Table(header []string, rows [][]Value)
}

type TextFormatter struct {
	w io.Writer
}

type Value interface {
	fmt.Stringer
	json.Marshaler
}

func ToValues(strings [][]string) [][]Value {
	var values [][]Value
	for _, row := range strings {
		var valueSlice []Value
		for _, r := range row {
			valueSlice = append(valueSlice, String(r))
		}
		values = append(values, valueSlice)
	}

	return values
}

type String string
type Int int
type Float float64
type Bool bool
type MultiValue []Value

func (m MultiValue) String() string {
	var s []string
	for _, v := range m {
		s = append(s, v.String())
	}

	return strings.Join(s, "\n")
}

func (m MultiValue) MarshalJSON() ([]byte, error) {
	var js []json.RawMessage
	for _, v := range m {
		r, err := v.MarshalJSON()
		if err != nil {
			return nil, err
		}
		js = append(js, r)
	}

	return json.Marshal(js)
}

func (s String) String() string {
	return string(s)
}

func (s String) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(s))
}

func (b Bool) String() string {
	return strconv.FormatBool(bool(b))
}

func (b Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(bool(b))
}

func (i Int) String() string {
	return strconv.Itoa(int(i))
}

func (i Int) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(i))
}

func (f Float) String() string {
	return strconv.FormatFloat(float64(f), 'G', 3, 64)
}

func (f Float) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(f))
}

func ValuesToStrings(sliceOfValues [][]Value) [][]string {
	var result [][]string
	for _, values := range sliceOfValues {
		var strings []string
		for _, value := range values {
			if value != nil {
				strings = append(strings, value.String())
			} else {
				strings = append(strings, "")
			}
		}
		result = append(result, strings)
	}
	return result
}

func (t TextFormatter) Details(rows [][]Value) {
	w := tablewriter.NewWriter(t.w)
	w.AppendBulk(ValuesToStrings(rows))
	w.SetBorder(false)
	w.SetColumnAlignment([]int{tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT})
	w.Render()
}

func (t TextFormatter) Table(header []string, rows [][]Value) {
	w := tablewriter.NewWriter(t.w)
	w.AppendBulk(ValuesToStrings(rows))
	w.SetBorder(false)
	w.SetHeader(header)
	w.Render()
}

type JSONFormatter struct {
	w io.Writer
}

func (j JSONFormatter) Details(rows [][]Value) {
	var err error
	m := make(map[string]json.RawMessage, len(rows))

	for _, row := range rows {
		if len(row) < 2 {
			continue
		}
		m[row[0].String()], err = json.Marshal(row[1])
		if err != nil {
			panic(err)
		}
	}

	encoder := json.NewEncoder(j.w)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(m)
}

func (j JSONFormatter) Table(header []string, rows [][]Value) {
	var err error
	m := make([]map[string]json.RawMessage, 0, len(rows))

	for _, row := range rows {
		r := make(map[string]json.RawMessage, len(header))
		for i := range header {
			r[header[i]], err = json.Marshal(row[i])
			if err != nil {
				panic(err)
			}
		}
		m = append(m, r)
	}

	encoder := json.NewEncoder(j.w)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(m)
}
