package spec

type CommandInputSpec struct {
	Args  []string `yaml:"args,flow"`
	Stdin string   `yaml:"stdin,omitempty"`
}

type JSONValueSpec struct {
	Nullable     bool   `yaml:"nullable,omitempty"`
	Int          bool   `yaml:"int,omitempty"`
	String       bool   `yaml:"string,omitempty"`
	Double       bool   `yaml:"double,omitempty"`
	Bool         bool   `yaml:"bool,omitempty"`
	Regex        string `yaml:"regex,omitempty"`
	ExampleValue string `yaml:"exampleValue"`
}

type JSONKeySpec struct {
	Name      string        `yaml:"name"`
	Required  bool          `yaml:"required,omitempty"`
	Forbidden bool          `yaml:"forbidden,omitempty"`
	Value     JSONValueSpec `yaml:"value,omitempty"`
}

type JSONSpec struct {
	Strict bool          `yaml:"strict"`
	Keys   []JSONKeySpec `yaml:"keys"`
}

type JSONArraySpec struct {
	Min   int      `yaml:"min,omitempty"`
	Max   int      `yaml:"max,omitempty"`
	Empty bool     `yaml:"empty,omitempty"`
	JSON  JSONSpec `yaml:"json"`
}

type OutputStreamSpec struct {
	Ignore    bool          `yaml:"ignore,omitempty"`
	Empty     bool          `yaml:"empty,omitempty"`
	JSON      JSONSpec      `yaml:"json,omitempty"`
	JSONArray JSONArraySpec `yaml:"jsonArray,omitempty"`
	Regex     string        `yaml:"regex,omitempty"`
	Outputs   []string      `yaml:"outputs,omitempty"`
}

type CommandOutputSpec struct {
	ExitCode int              `yaml:"exitCode"`
	Stdout   OutputStreamSpec `yaml:"stdout,omitempty"`
	Stderr   OutputStreamSpec `yaml:"stderr,omitempty"`
}

type CommandInputOutputSpec struct {
	ID     string
	Uses   []string          `yaml:"uses,omitempty"`
	After  []string          `yaml:"after,omitempty"`
	Before []string          `yaml:"before,omitempty"`
	Input  CommandInputSpec  `yaml:"input"`
	Output CommandOutputSpec `yaml:"output"`
}

type CommandSpec struct {
	ID          string
	Name        []string                 `yaml:"name,flow"`
	InputOutput []CommandInputOutputSpec `yaml:"inputOutput"`
}

type RootSpec struct {
	Commands []CommandSpec
}
