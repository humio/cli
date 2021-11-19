package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/humio/cli/tests/spec"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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

func mergeInto(receiver map[string]string, values map[string]string, valuePrefix string) {
	for k, v := range values {
		if valuePrefix != "" {
			receiver[valuePrefix+"."+k] = v
		} else {
			receiver[k] = v
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <specs dir>\n", os.Args[0])
		os.Exit(1)
	}

	specsDir := os.Args[1]

	specs, err := filepath.Glob(filepath.Join(specsDir, "*.yaml"))
	if err != nil {
		panic(err)
	}

	for _, spec := range specs {
		specBytes, err := ioutil.ReadFile(spec)
		if err != nil {
			panic(err)
		}

		var updatedSpec bytes.Buffer

		specInput := bytes.NewReader(specBytes)
		err = processSpec(specInput, &updatedSpec)
		if err != nil {
			panic(err)
		}

		updatedSpecFile := spec + ".updated"
		err = ioutil.WriteFile(updatedSpecFile, updatedSpec.Bytes(), 0600)
		if err != nil {
			panic(err)
		}
	}
}

func processSpec(specInput io.Reader, updateTo io.Writer) error {
	var rootSpec spec.RootSpec
	decoder := yaml.NewDecoder(specInput)
	decoder.SetStrict(true)
	err := decoder.Decode(&rootSpec)
	if err != nil {
		return err
	}

	type CmdWithIO struct {
		id        string
		cmd       *spec.CommandSpec
		io        *spec.CommandInputOutputSpec
		dependsOn []string
	}

	var cmdsWithIO []CmdWithIO
	for _, cmdSpec := range rootSpec.Commands {
		cmdSpec := cmdSpec
		for i := range cmdSpec.InputOutput {
			io := CmdWithIO{
				id:  fmt.Sprintf("%s.%s", cmdSpec.ID, cmdSpec.InputOutput[i].ID),
				cmd: &cmdSpec,
				io:  &cmdSpec.InputOutput[i],
			}

			for _, use := range cmdSpec.InputOutput[i].Uses {
				pieces := strings.SplitN(use, ".", 3)
				if len(pieces) != 3 {
					return fmt.Errorf("invalid Use, must be <cmd>.<io>.<output>, got %q", use)
				}

				io.dependsOn = append(io.dependsOn, fmt.Sprintf("%s.%s", pieces[0], pieces[1]))
			}

			for _, after := range cmdSpec.InputOutput[i].After {
				io.dependsOn = append(io.dependsOn, after)
			}

			cmdsWithIO = append(cmdsWithIO, io)
		}
	}

	for _, io := range cmdsWithIO {
		for _, before := range io.io.Before {
			var found bool
			for i := range cmdsWithIO {
				io2 := cmdsWithIO[i]
				if io2.id == before {
					cmdsWithIO[i].dependsOn = append(io2.dependsOn, io.id)
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("invalid before target %q, not found", before)
			}
		}
	}

	var sorted []CmdWithIO
	handled := map[string]bool{}

	// all without depends
	for _, io := range cmdsWithIO {
		if len(io.dependsOn) == 0 {
			sorted = append(sorted, io)
			handled[io.id] = true
		}
	}

	removeHandled := func() {
		old := cmdsWithIO
		cmdsWithIO = cmdsWithIO[:0]
		for _, i := range old {
			if !handled[i.id] {
				cmdsWithIO = append(cmdsWithIO, i)
			}
		}
	}

	n := 1
	for n > 0 {
		removeHandled()
		n = 0
		for _, io := range cmdsWithIO {
			dependencyMet := true
			for _, dep := range io.dependsOn {
				if !handled[dep] {
					dependencyMet = false
				}
			}

			if dependencyMet {
				sorted = append(sorted, io)
				handled[io.id] = true
				n++
			}
		}
	}

	if len(cmdsWithIO) > 0 {
		var ids []string
		for _, x := range cmdsWithIO {
			ids = append(ids, x.id)
		}
		fmt.Printf("circular dependency exists: %s\n", strings.Join(ids, ", "))
	}

	values := map[string]string{}

	for i := range sorted {
		io := sorted[i]
		outputs, err, updated := RunSpec(io.cmd, io.io, func(s string) (string, bool) {
			v, ok := values[s]
			return v, ok
		})
		if err != nil {
			fmt.Printf("%s (%s) failed: %s\n", io.cmd.Name, io.io.Input.Args, err)
		} else {
			mergeInto(values, outputs, "")
		}

		if updated != nil {
			*sorted[i].io = *updated
		}
	}

	if updateTo != nil {
		yaml.NewEncoder(updateTo).Encode(rootSpec)
	}

	return nil
}

func randomName() string {
	words := make([]string, 2)
	for i := range words {
		words[i] = wordlist[rand.Intn(len(wordlist))]
	}
	return strings.Join(words, "-")
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RunSpec(commandSpec *spec.CommandSpec, io *spec.CommandInputOutputSpec, input func(string) (string, bool)) (map[string]string, error, *spec.CommandInputOutputSpec) {
	res := map[string]string{}

	var replacements []string
	for _, use := range io.Uses {
		v, ok := input(use)
		if !ok {
			return nil, fmt.Errorf("missing input %q", use), nil
		}

		replacements = append(replacements, "{"+use+"}", v)
	}

	replacements = append(replacements, "{RANDOM_IDENTIFIER}", randomName())
	replacements = append(replacements, "{RANDOM_ID}", RandStringRunes(32))

	var args []string
	for _, a := range io.Input.Args {
		a = strings.NewReplacer(replacements...).Replace(a)

		args = append(args, a)
	}

	stdout, stderr, exitCode, err := executable.Run(commandSpec.Name, args)
	if err != nil {
		return nil, err, nil
	}

	if exitCode != io.Output.ExitCode {
		return nil, fmt.Errorf("expected exit code %d, got %d", io.Output.ExitCode, exitCode), nil
	}

	values, err := checkOutputStream(io.Output.Stdout, stdout)
	if err != nil {
		return nil, err, nil
	}
	_, err = checkOutputStream(io.Output.Stderr, stderr)
	if err != nil {
		return nil, err, nil
	}

	mergeInto(res, values, fmt.Sprintf("%s.%s", commandSpec.ID, io.ID))

	stdoutSpec, err := spec.GenerateOutputStreamSpec(stdout)
	if err != nil {
		return nil, err, nil
	}
	stderrSpec, err := spec.GenerateOutputStreamSpec(stderr)
	if err != nil {
		return nil, err, nil
	}

	updatedIOSpec := *io

	updatedIOSpec.Output.Stdout = stdoutSpec
	updatedIOSpec.Output.Stderr = stderrSpec

	return res, nil, &updatedIOSpec
}

func checkOutputStream(outputStreamSpec spec.OutputStreamSpec, outputStream []byte) (map[string]string, error) {
	res := map[string]string{}
	if outputStreamSpec.Ignore {
		return nil, nil
	}

	if outputStreamSpec.Empty && len(outputStream) > 0 {
		return nil, fmt.Errorf("expected stream to be empty, got %d bytes", len(outputStream))
	}

	if outputStreamSpec.Regex != "" {
		if re, err := regexp.Compile(outputStreamSpec.Regex); err != nil {
			return nil, err
		} else {
			submatch := re.FindStringSubmatch(string(outputStream))
			if submatch == nil {
				return nil, fmt.Errorf("output didn't match regex %q", outputStreamSpec.Regex)
			}

			for i, name := range re.SubexpNames() {
				if name != "" {
					res[name] = submatch[i]
				}
			}
		}
	}

	if len(outputStreamSpec.JSON.Keys) > 0 {
		js, err := jsonParse(outputStream)
		if err != nil {
			return nil, fmt.Errorf("expected output of JSON object type")
		}

		jsonSpec := outputStreamSpec.JSON
		values, err := checkJsonMap(jsonSpec.Keys, js, jsonSpec.Strict)
		if err != nil {
			return nil, err
		}

		mergeInto(res, values, "")
	}

	if len(outputStreamSpec.JSONArray.JSON.Keys) > 0 {
		js, err := jsonParseArray(outputStream)
		if err != nil {
			return nil, fmt.Errorf("expected output of JSON array type")
		}

		jsonArraySpec := outputStreamSpec.JSONArray
		jsonSpec := jsonArraySpec.JSON

		l := len(js)

		if l < jsonArraySpec.Min {
			return nil, fmt.Errorf("expected %d elements", jsonArraySpec.Min)
		}

		if 0 < jsonArraySpec.Max && l > jsonArraySpec.Max {
			return nil, fmt.Errorf("expected less than %d elements", jsonArraySpec.Max)
		}

		if jsonArraySpec.Empty && l > 0 {
			return nil, fmt.Errorf("expected no elements")
		}

		for i := range js {
			_, err = checkJsonMap(jsonSpec.Keys, js[i], jsonSpec.Strict)
			if err != nil {
				return nil, fmt.Errorf("element %d: %w", i, err)
			}
		}
	}

	return res, nil
}

func jsonParse(b []byte) (map[string]interface{}, error) {
	var m map[string]interface{}

	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.UseNumber()

	err := decoder.Decode(&m)
	return m, err
}

func jsonParseArray(b []byte) ([]map[string]interface{}, error) {
	var m []map[string]interface{}

	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.UseNumber()

	err := decoder.Decode(&m)
	return m, err
}

func checkJsonMap(keySpecs []spec.JSONKeySpec, js map[string]interface{}, strict bool) (map[string]string, error) {
	res := map[string]string{}

	for _, keySpec := range keySpecs {
		v, present := js[keySpec.Name]

		if keySpec.Required && !present {
			return nil, fmt.Errorf("missing required JSON key %q", keySpec.Name)
		}

		if keySpec.Forbidden && present {
			return nil, fmt.Errorf("found forbidden JSON key %q", keySpec.Name)
		}

		var empty spec.JSONValueSpec
		if keySpec.Value != empty && present {
			valueSpec := keySpec.Value

			if !valueSpec.Nullable && v == nil {
				return nil, fmt.Errorf("key %q contained null, but wasn't nullable", keySpec.Name)
			}

			if _, is := v.(bool); v != nil && valueSpec.Bool && !is {
				return nil, fmt.Errorf("key %q expected to be boolean, but was %T", keySpec.Name, v)
			}

			if number, is := v.(json.Number); v != nil && valueSpec.Int && !is {
				return nil, fmt.Errorf("key %q expected to be integer, but was %T", keySpec.Name, v)
			} else if v != nil && valueSpec.Int && is && strings.Contains(number.String(), ".") {
				return nil, fmt.Errorf("key %q expected to be integer, but value %q contained '.'", keySpec.Name, v)
			}

			if _, is := v.(json.Number); v != nil && valueSpec.Double && !is {
				return nil, fmt.Errorf("key %q expected to be double, but was %T", keySpec.Name, v)
			}

			if _, is := v.(string); v != nil && valueSpec.String && !is {
				return nil, fmt.Errorf("key %q expected to be string, but was %T", keySpec.Name, v)
			}

			var stringValue string
			switch t := v.(type) {
			case bool:
				if t {
					stringValue = "true"
				} else {
					stringValue = "false"
				}
			case json.Number:
				stringValue = t.String()
			case string:
				stringValue = t
			case nil:
				stringValue = "null"
			}

			if valueSpec.Regex != "" {
				re, err := regexp.Compile(valueSpec.Regex)
				if err != nil {
					return nil, err
				}

				if !re.MatchString(stringValue) {
					return nil, fmt.Errorf("value didn't match regex %q", valueSpec.Regex)
				}
			}

			res[keySpec.Name] = stringValue
		}

		delete(js, keySpec.Name)
		_ = v
	}

	if strict {
		if len(js) > 0 {
			var leftOverKeys []string
			for k := range js {
				leftOverKeys = append(leftOverKeys, strconv.Quote(k))
			}

			return nil, fmt.Errorf("found unmatched JSON keys: %s", strings.Join(leftOverKeys, ", "))
		}
	}
	return res, nil
}
