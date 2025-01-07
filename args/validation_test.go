package args

import (
	"cmdeagle/types"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a test cobra command
func createTestCommand() *cobra.Command {
	return &cobra.Command{
		Use: "test",
	}
}

// Helper function to create a basic args store for testing
func createTestArgsStore(cmd *cobra.Command, config *types.ArgsConfig, arguments []string) *ArgsStateStore {
	return CreateArgsStore(cmd, config, arguments)
}

func TestValidateArgs_BasicValidation(t *testing.T) {
	// Test cases for basic argument validation
	tests := []struct {
		name        string
		config      *types.ArgsConfig
		args        []string
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid required string arg",
			config: &types.ArgsConfig{
				Vars: []types.ArgVarDef{
					{
						Name:     "input",
						Type:     "string",
						Required: true,
					},
				},
			},
			args:        []string{"test.txt"},
			expectError: false,
		},
		{
			name: "missing required arg",
			config: &types.ArgsConfig{
				Vars: []types.ArgVarDef{
					{
						Name:     "input",
						Type:     "string",
						Required: true,
					},
				},
			},
			args:        []string{},
			expectError: true,
			errorMsg:    "missing required argument: input",
		},
		{
			name: "valid number type conversion",
			config: &types.ArgsConfig{
				Vars: []types.ArgVarDef{
					{
						Name: "count",
						Type: "number",
					},
				},
			},
			args:        []string{"42"},
			expectError: false,
		},
		{
			name: "invalid number type conversion",
			config: &types.ArgsConfig{
				Vars: []types.ArgVarDef{
					{
						Name: "count",
						Type: "number",
					},
				},
			},
			args:        []string{"not-a-number"},
			expectError: true,
		},
		{
			name: "valid boolean type conversion",
			config: &types.ArgsConfig{
				Vars: []types.ArgVarDef{
					{
						Name: "flag",
						Type: "boolean",
					},
				},
			},
			args:        []string{"true"},
			expectError: false,
		},
		{
			name: "default value when not provided",
			config: &types.ArgsConfig{
				Vars: []types.ArgVarDef{
					{
						Name:    "opt",
						Type:    "string",
						Default: "default-value",
					},
				},
			},
			args:        []string{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := createTestCommand()
			store := createTestArgsStore(cmd, tt.config, tt.args)
			err := ValidateArgs(cmd, tt.config, store)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateArgs_DependencyValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      *types.ArgsConfig
		args        []string
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid dependency - both args present",
			config: &types.ArgsConfig{
				Vars: []types.ArgVarDef{
					{
						Name: "source",
						Type: "string",
					},
					{
						Name: "destination",
						Type: "string",
						DependsOn: []types.ParamDependency{
							{
								Name: "source",
								When: &types.ParamConstraints{
									MinLength: toPtr(1),
								},
							},
						},
					},
				},
			},
			args:        []string{"src.txt", "dest.txt"},
			expectError: false,
		},
		{
			name: "invalid dependency - missing required dependent",
			config: &types.ArgsConfig{
				Vars: []types.ArgVarDef{
					{
						Name: "source",
						Type: "string",
					},
					{
						Name: "destination",
						Type: "string",
						DependsOn: []types.ParamDependency{
							{
								Name: "source",
								When: &types.ParamConstraints{
									MinLength: toPtr(1),
								},
							},
						},
					},
				},
			},
			args:        []string{"", "dest.txt"},
			expectError: true,
			errorMsg:    "Value is less than the minimum character length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := createTestCommand()
			store := createTestArgsStore(cmd, tt.config, tt.args)
			err := ValidateArgs(cmd, tt.config, store)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateArgs_ConflictValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      *types.ArgsConfig
		args        []string
		expectError bool
		errorMsg    string
	}{
		{
			name: "no conflict - only one arg provided",
			config: &types.ArgsConfig{
				Vars: []types.ArgVarDef{
					{
						Name:          "file",
						Type:          "string",
						ConflictsWith: []string{"url"},
					},
					{
						Name: "url",
						Type: "string",
					},
				},
			},
			args:        []string{"test.txt", ""},
			expectError: false,
		},
		{
			name: "conflict - both conflicting args provided",
			config: &types.ArgsConfig{
				Vars: []types.ArgVarDef{
					{
						Name:          "file",
						Type:          "string",
						ConflictsWith: []string{"url"},
					},
					{
						Name: "url",
						Type: "string",
					},
				},
			},
			args:        []string{"test.txt", "http://example.com"},
			expectError: true,
			errorMsg:    "conflicts with",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := createTestCommand()
			store := createTestArgsStore(cmd, tt.config, tt.args)
			err := ValidateArgs(cmd, tt.config, store)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateArgs_RuleValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      *types.ArgsConfig
		args        []string
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid - exact args",
			config: &types.ArgsConfig{
				Rules: []types.ArgRuleDef{
					{
						ExactArgs: 2,
					},
				},
			},
			args:        []string{"arg1", "arg2"},
			expectError: false,
		},
		{
			name: "invalid - too many args",
			config: &types.ArgsConfig{
				Rules: []types.ArgRuleDef{
					{
						ExactArgs: 2,
					},
				},
			},
			args:        []string{"arg1", "arg2", "arg3"},
			expectError: true,
			errorMsg:    "accepts 2 arg(s)",
		},
		{
			name: "valid - minimum args",
			config: &types.ArgsConfig{
				Rules: []types.ArgRuleDef{
					{
						MinimumNArgs: 2,
					},
				},
			},
			args:        []string{"arg1", "arg2", "arg3"},
			expectError: false,
		},
		{
			name: "valid - range args",
			config: &types.ArgsConfig{
				Rules: []types.ArgRuleDef{
					{
						RangeArgs: []int{1, 3},
					},
				},
			},
			args:        []string{"arg1", "arg2"},
			expectError: false,
		},
		{
			name: "invalid - out of range args",
			config: &types.ArgsConfig{
				Rules: []types.ArgRuleDef{
					{
						RangeArgs: []int{1, 3},
					},
				},
			},
			args:        []string{"arg1", "arg2", "arg3", "arg4"},
			expectError: true,
			errorMsg:    "accepts between",
		},
		{
			name: "valid - no args rule",
			config: &types.ArgsConfig{
				Rules: []types.ArgRuleDef{
					{
						NoArgs: true,
					},
				},
			},
			args:        []string{},
			expectError: false,
		},
		{
			name: "invalid - args provided when none allowed",
			config: &types.ArgsConfig{
				Rules: []types.ArgRuleDef{
					{
						NoArgs: true,
					},
				},
			},
			args:        []string{"arg1"},
			expectError: true,
			errorMsg:    "accepts no arg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := createTestCommand()
			store := createTestArgsStore(cmd, tt.config, tt.args)
			err := ValidateArgs(cmd, tt.config, store)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateArgs_ComplexRules(t *testing.T) {
	tests := []struct {
		name        string
		config      *types.ArgsConfig
		args        []string
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid - AND rule combination",
			config: &types.ArgsConfig{
				Rules: []types.ArgRuleDef{
					{
						And: &[]types.ArgRuleDef{
							{MinimumNArgs: 2},
							{MaximumNArgs: 4},
						},
					},
				},
			},
			args:        []string{"arg1", "arg2", "arg3"},
			expectError: false,
		},
		{
			name: "invalid - AND rule combination",
			config: &types.ArgsConfig{
				Rules: []types.ArgRuleDef{
					{
						And: &[]types.ArgRuleDef{
							{MinimumNArgs: 2},
							{MaximumNArgs: 4},
						},
					},
				},
			},
			args:        []string{"arg1"},
			expectError: true,
			errorMsg:    "requires at least 2 arg",
		},
		{
			name: "valid - OR rule combination",
			config: &types.ArgsConfig{
				Rules: []types.ArgRuleDef{
					{
						Or: &[]types.ArgRuleDef{
							{ExactArgs: 2},
							{ExactArgs: 4},
						},
					},
				},
			},
			args:        []string{"arg1", "arg2"},
			expectError: false,
		},
		{
			name: "invalid - OR rule combination",
			config: &types.ArgsConfig{
				Rules: []types.ArgRuleDef{
					{
						Or: &[]types.ArgRuleDef{
							{ExactArgs: 2},
							{ExactArgs: 4},
						},
					},
				},
			},
			args:        []string{"arg1", "arg2", "arg3"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := createTestCommand()
			store := createTestArgsStore(cmd, tt.config, tt.args)
			err := ValidateArgs(cmd, tt.config, store)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Helper function to convert value to pointer
func toPtr[T any](v T) *T {
	return &v
}
