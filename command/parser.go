package command

type testCase struct {
	Input  string
	Output map[string]string
}

type parserConfig struct {
	Name        string
	Description string     `yaml:",omitempty"`
	Tests       []testCase `yaml:",omitempty"`
	Example     string     `yaml:",omitempty"`
	Script      string     `yaml:",flow"`
}
