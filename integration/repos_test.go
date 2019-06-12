package integration

import (
	"fmt"
	"os/exec"
	"regexp"
	"testing"
)

func Test_Repos(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			"create new repo",
			[]string{"repos", "create", "test-repo"},
			fmt.Sprintf(repoCreate, "test-repo"),
		},
		{
			"show repo",
			[]string{"repos", "show", "test-repo"},
			fmt.Sprintf(repoShow, "test-repo"),
		},
		{
			"list repos",
			[]string{"repos", "list"},
			reposList,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			humioCtlCommand, err := humioCtl()
			if err != nil {
				t.Fatal(err)
			}
			args := []string{"run", humioCtlCommand}
			args = append(args, tt.args...)

			cmd := exec.Command("go", args...)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("error: %s, output: %s", err, output)
			}

			actual := stripWhitespace(string(output))
			expected := stripWhitespace(tt.expected)

			re := regexp.MustCompile(expected)

			if !re.Match([]byte(actual)) {
				t.Fatalf("actual = '%+v', expected = '%+v'", actual, expected)
			}
		})
	}
}
