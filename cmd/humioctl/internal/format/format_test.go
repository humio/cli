package format

import (
	"bytes"
	"testing"
)

func TestJSONFormatterDetails(t *testing.T) {
	var buf bytes.Buffer
	f := JSONFormatter{&buf}

	details := ToValues([][]string{
		{"Name", "Developer"},
		{"Is Root", "yes"},
	})

	f.Details(details)

	expected := `{
  "Is Root": "yes",
  "Name": "Developer"
}
`

	actual := buf.String()
	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestJSONFormatterTable(t *testing.T) {
	var buf bytes.Buffer
	f := JSONFormatter{&buf}

	header := []string{"Name", "Is Root"}

	rows := [][]string{
		{"Developer", "yes"},
		{"Operator", "no"},
	}

	f.Table(header, ToValues(rows))

	expected := `[
  {
    "Is Root": "yes",
    "Name": "Developer"
  },
  {
    "Is Root": "no",
    "Name": "Operator"
  }
]
`

	actual := buf.String()
	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}
