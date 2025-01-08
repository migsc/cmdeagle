package types

type ArgsConfig struct {
	// Define the named arguments and their types
	Vars []ArgVarDef `yaml:"vars,omitempty"`
	// Define validation rules for the argument list as a whole
	Rules []ArgRuleDef `yaml:"rules,omitempty"`
}

type ArgVarDef struct {
	Name string `yaml:"name"`
	Type string `yaml:"type,omitempty"` // If omitted, defaults to "string". Other valid values: "number", "boolean"

	// Description string `yaml:"description"`
	Required bool `yaml:"required,omitempty"`
	Default  any  `yaml:"default,omitempty"`
	// Optional validation for this specific argument
	Constraints ParamConstraints `yaml:"validation,omitempty"`
	// Dependency validations
	DependsOn     []ParamDependency `yaml:"depends-on,omitempty"`
	ConflictsWith []string          `yaml:"conflicts-with,omitempty"`
}

type ArgRuleDef struct {
	// Rules from Cobra

	NoArgs         bool  `yaml:"NoArgs,no-args,omitempty"`
	OnlyValidArgs  bool  `yaml:"OnlyValidArgs,only-valid-args,omitempty"`
	ArbitraryArgs  bool  `yaml:"ArbitraryArgs,arbitrary-args,omitempty"`
	MinimumNArgs   int   `yaml:"MinimumNArgs,minimum-n-args,max-args,omitempty"`
	MaximumNArgs   int   `yaml:"MaximumNArgs,maximum-n-args,min-args,omitempty"`
	ExactArgs      int   `yaml:"ExactArgs,exact-args,omitempty"`
	RangeArgs      []int `yaml:"RangeArgs,range-args,omitempty"`
	ExactValidArgs int   `yaml:"ExactValidArgs,exact-valid-args,omitempty"`

	// Boolean logic from Cobra
	MatchAll  *([]ArgRuleDef) `yaml:"MatchAll,match-all,omitempty"`
	MatchAny  *([]ArgRuleDef) `yaml:"MatchAny,match-any,omitempty"`
	MatchNone *([]ArgRuleDef) `yaml:"MatchNone,match-none,omitempty"`

	// Boolean logic from cmdeagle
	And  *([]ArgRuleDef) `yaml:"And,omitempty"`
	Or   *([]ArgRuleDef) `yaml:"Or,omitempty"`
	Nand *([]ArgRuleDef) `yaml:"Nand,omitempty"`

	Not *ArgRuleDef `yaml:"Not,omitempty"`
}
