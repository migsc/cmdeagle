# Argument Handling Implementation

## Current State
- Using Cobra for command structure
- Have a visitor pattern for commands
- Config defined in YAML:
```yaml
commands:
  - name: example
    args:
      - name: input
        type: string
        description: "Input file"
        required: true
```

## Tasks to Implement

### ARG-001: Basic arg parsing & validation
- Parse args from command line
- Match against schema definition
- Validate correct number of args

```yaml
commands:
name: example
    args:
        - name: input
          type: string
          description: "Input file"
    required: true
```

All args are required for now until optional args are implemented in #ARG-002.


### ARG-002: Required/optional arg support
- Check required args are provided
- Handle optional args correctly
- Error messaging for missing required args

### ARG-003: Default arg values
- Support default values in schema
- Apply defaults when arg not provided
- Type conversion for defaults

## Relevant Code

```go
type Arg struct {
    Name        string `yaml:"name"`
    Type        string `yaml:"type"`
    Description string `yaml:"description"`
    Required    bool   `yaml:"required,omitempty"`
    Default     any    `yaml:"default,omitempty"`
}
```

## Goals
- Clean error messages
- Type safety
- Pass args to run scripts properly
- Maintain Cobra compatibility

## Implementation Notes
- Arguments need to be validated at runtime
- Need to pass validated args to run scripts
- Must maintain compatibility with Cobra's arg handling
- Should support basic types (string, int, bool)

