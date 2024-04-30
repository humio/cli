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
		switch filename {
		case "api/interfaces.go":
			// Handling the ordering of interface declarations and their implementations is unpleasant within the template.
			// To make it easier we preprocess the schema by pulling out just the data we need and in the order we need.
			err := t.Execute(&buf, interfaces(*schema))
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to execute template: %v\n", err)
				os.Exit(1)
			}
		case "api/types.go":
			err := t.Execute(&buf, types(*schema))
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to execute template: %v\n", err)
				os.Exit(1)
			}
		default:
			err := t.Execute(&buf, schema)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to execute template: %v\n", err)
				os.Exit(1)
			}
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

type response struct {
	Data struct {
		Schema schema `json:"__schema"`
	} `json:"data"`
}

type schema struct {
	QueryType struct {
		Name string `json:"name"`
	} `json:"queryType,omitempty"`

	MutationType struct {
		Name string `json:"name"`
	} `json:"mutationType,omitempty"`

	SubscriptionType struct {
		Name string `json:"name"`
	} `json:"subscriptionType,omitempty"`

	Types []humioType `json:"types"`

	Directives []struct {
		Name        string       `json:"name"`
		Description string       `json:"description,omitempty"`
		Args        []inputValue `json:"args,omitempty"`
		Locations   []string     `json:"locations,omitempty"`
	} `json:"directives,omitempty"`
}

type humioType struct {
	Kind        string `json:"kind"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Fields      []struct {
		Name              string       `json:"name"`
		Description       string       `json:"description"`
		Args              []inputValue `json:"args,omitempty"`
		Type              typeRef      `json:"type,omitempty"`
		IsDeprecated      bool         `json:"isDeprecated"`
		DeprecationReason string       `json:"deprecatedReason"`
	} `json:"fields,omitempty"`
	InputFields   []inputValue `json:"inputFields,omitempty"`
	Interfaces    []typeRef    `json:"interfaces,omitempty"`
	EnumValues    []enumValue  `json:"enumValues"`
	PossibleTypes []typeRef    `json:"possibleTypes,omitempty"`
}

type enumValue struct {
	Name              string `json:"name"`
	Description       string `json:"description"`
	IsDeprecated      bool   `json:"isDeprecated"`
	DeprecationReason string `json:"deprecatedReason"`
}

type inputValue struct {
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Type         typeRef `json:"type"`
	DefaultValue string  `json:"defaultValue,omitempty"` // What type is this?
}

type typeRef struct {
	Kind   string   `json:"kind"`
	Name   string   `json:"name"`
	OfType *typeRef `json:"ofType,omitempty"`
}

func loadSchema(token, humioEndpoint string) (*schema, error) {
	var schema response
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
	return &schema.Data.Schema, err
}

// interfaces returns a map of GraphQL interface names to their implementation details.
func interfaces(schema schema) map[string][]humioType {
	ifaces := make(map[string][]humioType)

	for _, t := range schema.Types {
		if t.Kind == "INTERFACE" {
			ifaces[t.Name] = append([]humioType{t}, ifaces[t.Name]...)
		} else if t.Kind == "OBJECT" {
			for _, interfaceType := range t.Interfaces {
				ifaces[interfaceType.Name] = append(ifaces[interfaceType.Name], t)
			}
		}
	}
	return ifaces
}

func types(schema schema) []humioType {
	types := make([]humioType, 0)
	for _, t := range schema.Types {
		if t.Kind == "UNION" ||
			(t.Kind == "OBJECT" && t.EnumValues == nil && t.PossibleTypes == nil && len(t.Interfaces) == 0 && t.InputFields == nil) {
			// Ignore "internal" fields which all begin with "__".
			if !strings.HasPrefix(t.Name, "__") {
				if t.Name != "Mutation" {
					types = append(types, t)
				}
			}
		}
	}
	return types
}

var templates = map[string]*template.Template{
	"api/interfaces.go": t(`// Code generated by generate.go; DO NOT EDIT.
package api

import graphql "github.com/cli/shurcooL-graphql"

{{range $iface, $elems := .}}
{{range $elems}}{{if eq .Kind "INTERFACE"}}
{{- if .Description}}
// {{.Name}} {{.Description | clean | endSentence}}
{{- end}}
type {{.Name | identifier}} struct {
	{{- range .Fields}}
		{{- if ne .Name "fileFieldSearch"}}
			{{.Name | identifier}} {{.Type | type}}
		{{- end}}
	{{- end}}
}
{{else}}
{{- if .Description}}
// {{.Name |identifier }} {{.Description | clean | endSentence}}
{{- end}}
type {{.Name | identifier}} struct {
	{{$iface}}
	{{- range .Fields}}
		{{- if .Description}}
			// {{.Name |identifier }} {{.Description | clean | endSentence}}
		{{- end}}
		{{.Name | identifier}} {{.Type | type}}
	{{- end}}
}
{{end}}{{end}}{{end}}
`),

	"api/types.go": t(`// Code generated by generate.go; DO NOT EDIT.
package api

import graphql "github.com/cli/shurcooL-graphql"

// The following type aliases are curated in cmd/api_gen/generate.go by manual inspection of GraphQL type information.
type DateTime = string // Expected in ISO-8601 instant format, e.g. "2024-04-30T18:00:00.00Z"
type Email = string
type JSON = string
type Long = int64
type Markdown = string
type PackageName = string
type PackageScope = string
type PackageTag = string
type PackageVersion = string
type RepoOrViewName = string
type SemanticVersion = string
type URL = string
type UnversionedPackageSpecifier = string
type UrlOrData = string
type VersionedPackageSpecifier = string
type YAML = string
// End of manually curated type aliases.

{{range .}}
{{if eq .Kind "UNION"}}{{template "union" .}}{{end}}
{{if eq .Kind "OBJECT"}}{{template "object" .}}{{end}}
{{end}}

{{- define "union" -}}
{{- if .Description}}
// {{.Name | identifier}} {{.Description | clean | endSentence}}
{{- end}}
type {{.Name | identifier}} struct {
	{{- range .PossibleTypes}}
		{{.Name | identifier}}
	{{- end }}
}
{{- end -}}

{{- define "object" -}}
{{- if .Description}}
// {{.Name | identifier}} {{.Description | clean | endSentence}}
{{- end}}
type {{.Name | identifier}} struct {
	{{- range .Fields}}
	{{- if .Description}}
		// {{.Name | identifier}} {{.Description | clean | endSentence}}
	{{- end}}
	{{.Name | identifier}} {{.Type | type}}
	{{- end}}
}
{{- end -}}
`),

	"api/enum.go": t(`// Code generated by generate.go; DO NOT EDIT.
package api

{{range .Types | sortByName}}{{if and (eq .Kind "ENUM") (not (internal .Name))}}
{{template "enum" .}}
{{end}}{{end}}

{{- define "switchValues" -}}
{{- range .EnumValues }}
	case {{.Name | quote}}:
		return {{enumIdentifier $.Name .Name}}, true
{{- end -}}
{{- end -}}

{{- define "enum" -}}
{{- if .Description}}
// {{.Name}} {{.Description | clean | endSentence}}
{{- end}}
type {{.Name}} string

{{if .Description}}// {{.Description | clean | fullSentence}}{{end}}
const ({{range .EnumValues}}
	{{enumIdentifier $.Name .Name}} {{$.Name}} = {{.Name | quote}} {{if .Description}}// {{.Description | clean | fullSentence}}{{end}}{{end}}
)

func Valid{{.Name}}(s string) ({{.Name}}, bool) {
	switch s {
	{{- template "switchValues" .}}
	}
	return {{$.Name}}(""), false
}
{{- end -}}
	`),

	"api/input.go": t(`// Code generated by generate.go; DO NOT EDIT.

package api

import "github.com/cli/shurcooL-graphql"

// Input represents one of the Input structs:
//
// {{join (inputObjects .Types) ", "}}.
type Input interface{}
{{range .Types | sortByName}}{{if eq .Kind "INPUT_OBJECT"}}
{{template "inputObject" .}}
{{end}}{{end}}

{{- define "inputObject" -}}
{{- if .Description}}
// {{.Name}} {{.Description | clean | endSentence}}
{{- end}}
     type {{.Name}} struct {
     	{{- range .InputFields}}
     		{{- if eq .Type.Kind "NON_NULL"}}
     			{{- if .Description}}
     				{{printf "// %s" .Description | clean | fullSentence}}
     			{{- end}}
     			{{.Name | identifier}} {{.Type | type}} ` + "`" + `json:"{{.Name}}"` + "`" + ` // Required
     		{{- end}}
     	{{- end }}
     	{{- range .InputFields}}
     		{{- if ne .Type.Kind "NON_NULL"}}
     			{{- if .Description}}
     				{{printf "// %s" .Description | clean | fullSentence}}
     			{{- end}}
     			{{.Name | identifier}} {{.Type | type}} ` + "`" + `json:"{{.Name}},omitempty"` + "`" + ` // Optional
     		{{- end}}
     	{{- end}}
     }

{{- end -}}
`),
}

func t(text string) *template.Template {
	// typeString returns a string representation of GraphQL type t.
	var typeString func(t *typeRef) string
	typeString = func(t *typeRef) string {
		if t == nil {
			return ""
		}
		baseTypes := []string{"Boolean", "Int", "Float", "String", "ID"}
		switch t.Kind {
		case "NON_NULL":
			s := typeString(t.OfType)
			if !strings.HasPrefix(s, "*") {
				panic(fmt.Errorf("nullable type %q doesn't begin with '*'", s))
			}
			strippedName := s[1:] // Strip star from nullable type to make it non-null.
			if slices.Contains(baseTypes, strippedName) {
				return "graphql." + strippedName
			}
			return strippedName
		case "LIST":
			return "*[]" + typeString(t.OfType)
		default:
			name := t.Name
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
		"sortByName": func(types []humioType) []humioType {
			sort.Slice(types, func(i, j int) bool {
				ni := types[i].Name
				nj := types[j].Name
				return ni < nj
			})
			return types
		},
		"inputObjects": func(types []humioType) []string {
			var names []string
			for _, t := range types {
				if t.Kind != "INPUT_OBJECT" {
					continue
				}
				names = append(names, t.Name)
			}
			sort.Strings(names)
			return names
		},
		"union": func(types []typeRef) []string {
			var names []string
			for _, pt := range types {
				names = append(names, pt.Name)
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
	    }
	  }
	}
}`
