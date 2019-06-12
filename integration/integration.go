package integration

import (
	"os"
	"path"
	"strings"
)

func humioCtl() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return path.Join(dir, "..", "main.go"), nil
}

func stripWhitespace(data string) string {
	var str strings.Builder
	strSplit := strings.Split(string(data), "\n")

	for i, s := range strSplit {
		trimmed := strings.TrimSpace(s)
		if trimmed != "" {
			str.WriteString(trimmed)

			if i < len(strSplit)-1 {
				str.WriteString("\n")
			}
		}
	}

	return str.String()
}
