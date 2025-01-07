# Flag Handling Implementation

## Current State
- Using Cobra for command structure
- Have visitor pattern for commands
- Config defined in YAML:
```yaml
commands:
  - name: example
    flags:
      - name: verbose
        shorthand: v
        type: boolean
        description: "Enable verbose output"
        default: false
```

## Implementation Tickets

### FLAG-001: Basic flag parsing
- Parse flags from command line
- Match against schema definition
- Pass to run scripts as `${flags.name}`
- Support basic boolean flags only

### FLAG-002: Flag type support
- Support string flags
- Support integer flags
- Support float flags
- Type validation
- Type conversion

### FLAG-003: Default flag values
- Support default values in schema
- Apply defaults when flag not provided
- Type-specific default handling

### FLAG-004: Flag shorthand support
- Support -v style shorthand flags
- Validate no shorthand conflicts
- Support both styles in scripts

## Relevant Code
```go
type Flag struct {
    Name        string `yaml:"name"`
    Shorthand   string `yaml:"shorthand,omitempty"`
    Type        string `yaml:"type"`
    Description string `yaml:"description"`
    Default     any    `yaml:"default,omitempty"`
}
```

## Goals
- Clean error messages
- Pass flags to run scripts properly
- Maintain Cobra compatibility
- Support basic boolean flags first
- Use same interpolation system as args (`${flags.name}`)

## Implementation Notes
- Flags need to be registered with Cobra during command creation
- Need to pass validated flags to run scripts
- Will use similar structure to args package
- Should start with just boolean flags for FLAG-001