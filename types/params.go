package types

type ParamDependency struct {
	Name string            `yaml:"arg"` // can be `flags.your_flag` or `args.your_arg` or `args[0]`
	When *ParamConstraints `yaml:"when,omitempty"`
}

type ParamConstraints struct {
	// Numeric validations
	MinValue   *float64 `yaml:"min,omitempty"` // Will be nil if not specified in YAML
	MaxValue   *float64 `yaml:"max,omitempty"`
	MultipleOf *float64 `yaml:"multipleOf,omitempty"`
	Eq         any      `yaml:"eq,omitempty"`    // equals
	Neq        any      `yaml:"neq,omitempty"`   // not equals
	Gt         any      `yaml:"gt,omitempty"`    // greater than
	Gte        any      `yaml:"gte,omitempty"`   // greater than or equal
	Lt         any      `yaml:"lt,omitempty"`    // less than
	Lte        any      `yaml:"lte,omitempty"`   // less than or equal
	In         []any    `yaml:"in,omitempty"`    // value is in list
	NotIn      []any    `yaml:"notIn,omitempty"` // value is not in list
	// Could add more like:
	// Pattern string `yaml:"pattern,omitempty"` // regex match
	// Exists  bool   `yaml:"exists,omitempty"`  // just check if arg exists

	// String validations
	MinLength *int   `yaml:"min-length,omitempty"`
	MaxLength *int   `yaml:"max-length,omitempty"`
	Pattern   string `yaml:"pattern,omitempty"`

	// File/Path validations
	FileExists     string `yaml:"file-exists,omitempty"`
	DirExists      string `yaml:"dir-exists,omitempty"`
	IsFileType     string `yaml:"is-file-type,omitempty"`
	HasPermissions string `yaml:"has-permissions,omitempty"`

	// Conditionals
	And  ([]*ParamConstraints) `yaml:"and,omitempty"`
	Nand ([]*ParamConstraints) `yaml:"nand,omitempty"`
	Or   ([]*ParamConstraints) `yaml:"or,omitempty"`
	Not  (*ParamConstraints)   `yaml:"not,omitempty"`
}

var ConstraintFileKeys = []string{"FileExists", "DirExists", "HasPermissions", "IsFileType"}

var ConstraintConditionals = []string{"And", "Or", "Not"}
