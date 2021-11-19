package spec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func GenerateOutputStreamSpec(stream []byte) (OutputStreamSpec, error) {
	if len(stream) == 0 {
		return OutputStreamSpec{Empty: true}, nil
	}

	if stream[0] == '{' {
		return generateJsonOutputStreamSpec(stream)
	}

	if stream[0] == '[' {
		return generateJsonArrayOutputStreamSpec(stream)
	}

	return generateRegexOutputStreamSpec(stream)
}

var wsBeforeRegex = regexp.MustCompile(`^\s+`)
var wsAfterRegex = regexp.MustCompile(`\s+$`)
var nlAfterRegex = regexp.MustCompile(`\n+$`)
var wordCouldBeIdentifier = regexp.MustCompile(`[a-zA-Z0-9]{16,}`)

func generateRegexOutputStreamSpec(stream []byte) (OutputStreamSpec, error) {
	str := string(stream)

	withoutSpace := strings.TrimSpace(str)

	hasPrefixWS := wsBeforeRegex.MatchString(str)
	hasSuffixWS := wsAfterRegex.MatchString(str)
	hasNewline := nlAfterRegex.MatchString(str)

	lines := strings.Split(withoutSpace, "\n")

	lineCount := len(lines)

	if lineCount == 1 {
		line := lines[0]

		words := strings.Split(line, " ")

		var regexpWords []string

		i := 0
		var outputs []string
		for _, w := range words {
			if wordCouldBeIdentifier.MatchString(w) {
				regexpWords = append(regexpWords, fmt.Sprintf(`(?P<identifier%d>[a-zA-Z0-9]{%d})`, i, len(w)))
				outputs = append(outputs, fmt.Sprintf("identifier%d", i))
				i++
			} else {
				regexpWords = append(regexpWords, regexp.QuoteMeta(w))
			}
		}

		var prefix, suffix string
		if hasPrefixWS {
			prefix = `\s+`
		}
		if hasSuffixWS {
			if hasNewline {
				suffix = `\s*\n+`
			} else {
				suffix = `\s+`
			}
		}

		re := "^" + prefix + strings.Join(regexpWords, " ") + suffix + "$"

		return OutputStreamSpec{
			Regex:   re,
			Outputs: outputs,
		}, nil
	}

	return OutputStreamSpec{
		Regex: "^\\s*" + regexp.QuoteMeta(strings.TrimSpace(string(stream))) + "\\s*$",
	}, nil
}

func generateJsonArrayOutputStreamSpec(stream []byte) (OutputStreamSpec, error) {
	var res OutputStreamSpec

	decoder := json.NewDecoder(bytes.NewReader(stream))
	decoder.UseNumber()

	var m []json.RawMessage

	err := decoder.Decode(&m)
	if err != nil {
		return OutputStreamSpec{}, err
	}

	if len(m) == 0 {
		return OutputStreamSpec{
			JSONArray: JSONArraySpec{
				Empty: true,
			},
		}, nil
	}

	jsonSpec, err := generateJsonOutputStreamSpec(m[0])
	if err != nil {
		return OutputStreamSpec{}, err
	}

	res.JSONArray.Min = 1
	res.JSONArray.JSON = jsonSpec.JSON

	return res, nil
}

func generateJsonOutputStreamSpec(stream []byte) (OutputStreamSpec, error) {
	var res OutputStreamSpec
	res.JSON.Strict = true

	decoder := json.NewDecoder(bytes.NewReader(stream))
	decoder.UseNumber()

	var m map[string]interface{}

	err := decoder.Decode(&m)
	if err != nil {
		return OutputStreamSpec{}, err
	}

	for k, v := range m {
		var keySpec JSONKeySpec

		keySpec.Name = k
		keySpec.Required = true

		switch x := v.(type) {
		case json.Number:
			if strings.Contains(x.String(), ".") {
				keySpec.Value.Double = true
			} else {
				keySpec.Value.Int = true
			}
		case string:
			keySpec.Value.String = true
		case bool:
			keySpec.Value.Bool = true
		case nil:
			keySpec.Value.Nullable = true
		}

		keySpec.Value.ExampleValue = fmt.Sprint(v)

		res.JSON.Keys = append(res.JSON.Keys, keySpec)

		res.Outputs = append(res.Outputs, k)
	}

	return res, nil
}
