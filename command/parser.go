package command

type testCase struct {
	Input  string
	Output map[string]string
}

type parserConfig struct {
	Name        string
	Description string
	Tests       []testCase
	Example     string
	Script      string
}
