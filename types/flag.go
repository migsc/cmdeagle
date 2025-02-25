package types

type FlagDefinition struct {
	Name          string              `yaml:"name"`
	Type          string              `yaml:"type"`
	Required      bool                `yaml:"required,omitempty"`
	Default       any                 `yaml:"default,omitempty"`
	Description   string              `yaml:"description,omitempty"`
	Shorthand     string              `yaml:"short,omitempty"`
	Hidden        bool                `yaml:"hidden,omitempty"`
	DependsOn     []*ParamDependency  `yaml:"depends-on,omitempty"`
	ConflictsWith []string            `yaml:"conflicts-with,omitempty"`
	Constraints   *ParamConstraints   `yaml:"constraints,omitempty"`
	Rules         []*ParamConstraints `yaml:"rules,omitempty"`
	Pattern       string              `yaml:"pattern,omitempty"`
}
