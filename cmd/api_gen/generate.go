package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io"
	"net/http"
	"net/url"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/cli/shurcooL-graphql/ident"
)

func main() {
	token, ok := os.LookupEnv("HUMIO_TOKEN")
	if !ok {
		fmt.Fprintf(os.Stderr, "HUMIO_TOKEN environment variable not set\n")
		os.Exit(1)
	}
	endpoint, ok := os.LookupEnv("HUMIO_ENDPOINT")
	if !ok {
		fmt.Fprintf(os.Stderr, "HUMIO_ENDPOINT environment variable not set\n")
		os.Exit(1)
	}
	schema, err := loadSchema(token, endpoint)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load schema: %v\n", err)
		os.Exit(1)
	}

	for filename, t := range templates {
		var buf bytes.Buffer
		err := t.Execute(&buf, schema)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to execute template: %v\n", err)
			os.Exit(1)
		}
		out, err := format.Source(buf.Bytes())
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to format source: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("writing", filename)
		err = os.WriteFile(filename, out, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to write file: %v\n", err)
			os.Exit(1)
		}
	}
}

func loadSchema(token, humioEndpoint string) (interface{}, error) {
	var schema interface{}
	path, err := url.JoinPath(humioEndpoint, "/graphql")
	if err != nil {
		return nil, err
	}

	in := struct {
		Query string `json:"query"`
	}{
		Query: introspectionQuery,
	}
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(in)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, path, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("authorization", "bearer "+token)
	req.Header.Set("content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("non-200 status code: %v body: %q", resp.StatusCode, body)
	}
	err = json.NewDecoder(resp.Body).Decode(&schema)
	return schema, err
}

var templates = map[string]*template.Template{
	"api/enum.go": t(`// Code generated by generate.go; DO NOT EDIT.
package api

import graphql "github.com/cli/shurcooL-graphql"

{{range .data.__schema.types | sortByName}}{{if and (eq .kind "ENUM") (not (internal .name))}}
{{template "enum" .}}
{{end}}{{end}}

{{- define "switchValues" -}}
{{- $valCount := len .enumValues -}}
{{- range .enumValues }}
	case {{.name | quote}}:
		return {{printf "%s%s" $.name .name}}, true
{{- end -}}
{{- end -}}

{{- define "enum" -}}
// {{.name}} {{if .description}}{{.description | clean | endSentence}}{{end}}
type {{.name}} string

{{if .description}}// {{.description | clean | fullSentence}}{{end}}
const ({{range .enumValues}}
	{{enumIdentifier $.name .name}} {{$.name}} = {{.name | quote}} {{if .description}}// {{.description | clean | fullSentence}}{{end}}{{end}}
)

func Valid{{.name}}(s string) ({{.name}}, bool) {
	switch s {
	{{template "switchValues" .}}
	}
	return {{$.name}}(""), false
}
{{- end -}}
	`),

	"api/input.go": t(`// Code generated by generate.go; DO NOT EDIT.

package api

import graphql "github.com/cli/shurcooL-graphql"

type Long int64
type JSON string
type YAML string
type DateTime string
type RepoOrViewName string

// Input represents one of the Input structs:
//
// {{join (inputObjects .data.__schema.types) ", "}}.
type Input interface{}
{{range .data.__schema.types | sortByName}}{{if eq .kind "INPUT_OBJECT"}}
{{template "inputObject" .}}
{{end}}{{end}}

{{- define "inputObject" -}}
// {{.name}} {{if .description}}{{.description | clean | endSentence}}{{end}}
type {{.name}} struct {{"{"}}{{range .inputFields}}{{if eq .type.kind "NON_NULL"}}

	// {{if .description}}{{.description | clean | fullSentence}} {{end}}(Required)
	{{.name | identifier}} {{.type | type}} ` + "`" + `json:"{{.name}}"` + "`" + `{{end}}{{end}}

{{range .inputFields}}{{if ne .type.kind "NON_NULL"}}

		// {{if .description}}{{.description | clean | fullSentence}} {{end}}(Optional)
		{{.name | identifier}} {{.type | type}} ` + "`" + `json:"{{.name}},omitempty"` + "`" + `{{end}}{{end}}
	}

{{- end -}}
`),
}

func t(text string) *template.Template {
	// typeString returns a string representation of GraphQL type t.
	var typeString func(t map[string]interface{}) string
	typeString = func(t map[string]interface{}) string {
		baseTypes := []string{"Boolean", "Int", "Float", "String", "ID"}
		switch t["kind"] {
		case "NON_NULL":
			s := typeString(t["ofType"].(map[string]interface{}))
			if !strings.HasPrefix(s, "*") {
				panic(fmt.Errorf("nullable type %q doesn't begin with '*'", s))
			}
			strippedName := s[1:] // Strip star from nullable type to make it non-null.
			if slices.Contains(baseTypes, strippedName) {
				return "graphql." + strippedName
			}
			return strippedName
		case "LIST":
			return "*[]" + typeString(t["ofType"].(map[string]interface{}))
		default:
			name := t["name"].(string)
			if slices.Contains(baseTypes, name) {
				return "*graphql." + name
			}
			return "*" + name
		}
	}

	return template.Must(template.New("").Funcs(template.FuncMap{
		"sub":      func(a, b int) int { return a - b },
		"internal": func(s string) bool { return strings.HasPrefix(s, "__") },
		"quote":    strconv.Quote,
		"join":     strings.Join,
		"sortByName": func(types []interface{}) []interface{} {
			sort.Slice(types, func(i, j int) bool {
				ni := types[i].(map[string]interface{})["name"].(string)
				nj := types[j].(map[string]interface{})["name"].(string)
				return ni < nj
			})
			return types
		},
		"inputObjects": func(types []interface{}) []string {
			var names []string
			for _, t := range types {
				t := t.(map[string]interface{})
				if t["kind"].(string) != "INPUT_OBJECT" {
					continue
				}
				names = append(names, t["name"].(string))
			}
			sort.Strings(names)
			return names
		},
		"identifier": func(name string) string { return ident.ParseLowerCamelCase(name).ToMixedCaps() },
		"enumIdentifier": func(enum, value string) string {
			return enum + ident.ParseMixedCaps(value).ToMixedCaps()
		},
		"type":  typeString,
		"clean": func(s string) string { return strings.Join(strings.Fields(s), " ") },
		"endSentence": func(s string) string {
			s = strings.ToLower(s[0:1]) + s[1:]
			switch {
			default:
				s = "represents " + s
			case strings.HasPrefix(s, "autogenerated "):
				s = "is an " + s
			case strings.HasPrefix(s, "specifies "):
				// Do nothing.
			}
			if !strings.HasSuffix(s, ".") {
				s += "."
			}
			return s
		},
		"fullSentence": func(s string) string {
			if !strings.HasSuffix(s, ".") {
				s += "."
			}
			return s
		},
	}).Parse(text))
}

const introspectionQuery = `query IntrospectionQuery {
	__schema {
	  queryType { name }
	  mutationType { name }
	  subscriptionType { name }
	  types {
		...FullType
	  }
	  directives {
		name
		description
		args {
		  ...InputValue
		}
		locations
	  }
	}
	}

	fragment FullType on __Type {
	kind
	name
	description
	fields(includeDeprecated: true) {
	  name
	  description
	  args {
		...InputValue
	  }
	  type {
		...TypeRef
	  }
	  isDeprecated
	  deprecationReason
	}
	inputFields {
	  ...InputValue
	}
	interfaces {
	  ...TypeRef
	}
	enumValues(includeDeprecated: true) {
	  name
	  description
	  isDeprecated
	  deprecationReason
	}
	possibleTypes {
	  ...TypeRef
	}
	}

	fragment InputValue on __InputValue {
	name
	description
	type { ...TypeRef }
	defaultValue
	}

	fragment TypeRef on __Type {
	kind
	name
	ofType {
	  kind
	  name
	  ofType {
		kind
		name
		ofType {
		  kind
		  name
		}
	  }
	}
}`
