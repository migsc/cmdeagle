package types

// TODO: Convert all YAML field names from camelCase to kebab-case to follow Go community conventions
// Examples:
// - allowUnknownFlags -> allow-unknown-flags
// - globalFlags -> global-flags
// - logLevel -> log-level
// This is a breaking change that should be done in a major version update.
//
// TODO: Exception to kebab-case: Use CapitalCase for Cobra-compatible arg validations
// to match Cobra's function names, while keeping custom boolean logic (and/or/not) in kebab-case
// Examples:
// - MinimumNArgs (matches cobra.MinimumNArgs)
// - OnlyValidArgs (matches cobra.OnlyValidArgs)
// - and/or/not (our custom boolean logic)
// This makes it clear which validations map directly to Cobra functions.

type CmdeagleConfig struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
	Author      string `yaml:"author"`
	License     string `yaml:"license"`
	// BinDir     string `yaml:"bin-path"`
	// DataPath    string `yaml:"data-path"`

	/*
		bin-dir: "/usr/local/bin"   # macOS (system-wide): /usr/local/bin
		# macOS (user-only):   /usr/local/bin or ~/bin
		# Windows (user-only): C:\Users\<username>\AppData\Local\Programs\MyApp\bin
		# Windows (system-wide): C:\Program Files\MyApp\bin

		data-dir: "/usr/local/share/myapp"   # macOS (system-wide): /usr/local/share/myapp
		# macOS (user-only):   ~/Library/Application Support/MyApp
		# Windows (user-only): %LocalAppData%\MyApp
		# Windows (system-wide): C:\ProgramData\MyApp
	*/

	// Settings    Settings            `yaml:"settings,omitempty"`

	Args       []ArgDefinition     `yaml:"args,omitempty"`
	Flags      []FlagDefinition    `yaml:"flags,omitempty"`
	Commands   []CommandDefinition `yaml:"commands"`
	Requires   map[string]string   `yaml:"requires,omitempty"`
	Includes   []string            `yaml:"includes,omitempty"`
	Build      string              `yaml:"build,omitempty"`
	Validate   string              `yaml:"validate,omitempty"`
	Start      string              `yaml:"start,omitempty"`
	Completion bool                `yaml:"completion`
}

// type Settings struct {
// 	AllowUnknownFlags bool   `yaml:"allow_unknown_flags"`
// 	StrictArgs        bool   `yaml:"strict_args"`
// 	ColorOutput       bool   `yaml:"color_output"`
// 	LogLevel          string `yaml:"log_level"`
// }
