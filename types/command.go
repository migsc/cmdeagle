package types

type CommandDefinition struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Aliases     []string `yaml:"aliases,omitempty"`
	// If omitted: Args will be nil
	// If key exists but empty (args: []): Args will be empty slice
	// If has values: Args will contain the values
	Args     ArgsConfig          `yaml:"args,omitempty"`
	Flags    []FlagDefinition    `yaml:"flags,omitempty"`
	Commands []CommandDefinition `yaml:"commands,omitempty"`
	Requires map[string]string   `yaml:"requires,omitempty"`
	Includes []string            `yaml:"includes,omitempty"`
	Build    string              `yaml:"build,omitempty"`
	Start    string              `yaml:"start,omitempty"`
}
