package types

// type ArgsConfig struct {
// 	// Define the named arguments and their types
// 	Vars []ArgVarDef `yaml:"vars,omitempty"`
// 	// Define validation rules for the argument list as a whole
// 	// Rules []ArgRuleDef `yaml:"rules,omitempty"`
// }

// type ArgsConfig = []ArgDef

type ArgDefinition struct {
	Name string `yaml:"name"`
	Type string `yaml:"type,omitempty"` // If omitted, defaults to "string". Other valid values: "number", "boolean"

	// Description string `yaml:"description"`
	Required bool `yaml:"required,omitempty"`
	Default  any  `yaml:"default,omitempty"`
	// Optional validation for this specific argument
	// TODO: rename this to rules? right?
	Constraints ParamConstraints `yaml:"validation,omitempty"`
	// Dependency validations
	DependsOn     []ParamDependency `yaml:"depends-on,omitempty"`
	ConflictsWith []string          `yaml:"conflicts-with,omitempty"`
}

// type ArgRuleDef struct {
// 	// Rules from Cobra

// 	NoArgs         bool  `yaml:"no-args,omitempty"`
// 	OnlyValidArgs  bool  `yaml:"only-valid-args,omitempty"`
// 	ArbitraryArgs  bool  `yaml:"arbitrary-args,omitempty"`
// 	MinimumNArgs   int   `yaml:"minimum-n-args,omitempty"`
// 	MaximumNArgs   int   `yaml:"maximum-n-args,omitempty"`
// 	ExactArgs      int   `yaml:"exact-args,omitempty"`
// 	RangeArgs      []int `yaml:"range-args,omitempty"`
// 	ExactValidArgs int   `yaml:"exact-valid-args,omitempty"`

// 	// Boolean logic from Cobra
// 	MatchAll  *([]ArgRuleDef) `yaml:"match-all,omitempty"`
// 	MatchAny  *([]ArgRuleDef) `yaml:"match-any,omitempty"`
// 	MatchNone *([]ArgRuleDef) `yaml:"match-none,omitempty"`

// 	// Boolean logic from cmdeagle
// 	And  *([]ArgRuleDef) `yaml:"and,omitempty"`
// 	Or   *([]ArgRuleDef) `yaml:"or,omitempty"`
// 	Nand *([]ArgRuleDef) `yaml:"nand,omitempty"`

// 	Not *ArgRuleDef `yaml:"not,omitempty"`
// }
