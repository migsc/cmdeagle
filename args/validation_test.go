package args

import (
	"testing"

	"github.com/migsc/cmdeagle/types"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestValidateArgs(t *testing.T) {
	cmd := &cobra.Command{
		Use: "testcmd",
	}

	t.Run("handles nil config and store", func(t *testing.T) {
		err := ValidateArgs(cmd, nil, nil)
		assert.NoError(t, err)
	})

	t.Run("validates dependencies", func(t *testing.T) {
		config := &types.ArgsConfig{
			Vars: []types.ArgVarDef{
				{Name: "dep"},
				{Name: "arg", DependsOn: []types.ParamDependency{{Name: "dep", When: &types.ParamConstraints{Eq: "value"}}}},
			},
		}
		store := CreateArgsStore(cmd, config, []string{"value", "test"})

		err := ValidateArgs(cmd, config, store)
		assert.NoError(t, err)
	})

	t.Run("validates multiple dependencies", func(t *testing.T) {
		config := &types.ArgsConfig{
			Vars: []types.ArgVarDef{
				{Name: "dep1"},
				{Name: "dep2"},
				{Name: "arg", DependsOn: []types.ParamDependency{
					{Name: "dep1", When: &types.ParamConstraints{Eq: "value1"}},
					{Name: "dep2", When: &types.ParamConstraints{Eq: "value2"}},
				}},
			},
		}
		store := CreateArgsStore(cmd, config, []string{"value1", "value2", "test"})

		err := ValidateArgs(cmd, config, store)
		assert.NoError(t, err)

		// Test failure case
		store = CreateArgsStore(cmd, config, []string{"value1", "wrong", "test"})
		err = ValidateArgs(cmd, config, store)
		assert.Error(t, err)
	})

	t.Run("validates conflicts", func(t *testing.T) {
		config := &types.ArgsConfig{
			Vars: []types.ArgVarDef{
				{Name: "arg1"},
				{Name: "arg2", ConflictsWith: []string{"arg1"}},
			},
		}
		store := CreateArgsStore(cmd, config, []string{"value", "test"})

		err := ValidateArgs(cmd, config, store)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "conflicts with")
	})

	t.Run("validates multiple conflicts", func(t *testing.T) {
		config := &types.ArgsConfig{
			Vars: []types.ArgVarDef{
				{Name: "arg1"},
				{Name: "arg2"},
				{Name: "arg3", ConflictsWith: []string{"arg1", "arg2"}},
			},
		}
		store := CreateArgsStore(cmd, config, []string{"value1", "value2", "test"})

		err := ValidateArgs(cmd, config, store)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "conflicts with")

		// Test no conflict when conflicting args are empty
		store = CreateArgsStore(cmd, config, []string{"", "", "test"})
		err = ValidateArgs(cmd, config, store)
		assert.NoError(t, err)
	})
}

func TestValidateRule(t *testing.T) {
	cmd := &cobra.Command{
		Use: "testcmd",
	}

	t.Run("validates NoArgs", func(t *testing.T) {
		rule := types.ArgRuleDef{NoArgs: true}
		err := validateRule(cmd, rule, []string{"arg"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "accepts no args")

		err = validateRule(cmd, rule, []string{})
		assert.NoError(t, err)
	})

	t.Run("validates OnlyValidArgs", func(t *testing.T) {
		cmd.ValidArgs = []string{"valid1", "valid2"}
		rule := types.ArgRuleDef{OnlyValidArgs: true}

		err := validateRule(cmd, rule, []string{"valid1"})
		assert.NoError(t, err)

		err = validateRule(cmd, rule, []string{"invalid"})
		assert.Error(t, err)
	})

	t.Run("validates ArbitraryArgs", func(t *testing.T) {
		rule := types.ArgRuleDef{ArbitraryArgs: true}
		err := validateRule(cmd, rule, []string{"anything", "goes", "here"})
		assert.NoError(t, err)
	})

	t.Run("validates MinimumNArgs", func(t *testing.T) {
		rule := types.ArgRuleDef{MinimumNArgs: 2}
		err := validateRule(cmd, rule, []string{"arg1"})
		assert.Error(t, err)

		err = validateRule(cmd, rule, []string{"arg1", "arg2"})
		assert.NoError(t, err)
	})

	t.Run("validates MaximumNArgs", func(t *testing.T) {
		rule := types.ArgRuleDef{MaximumNArgs: 2}
		err := validateRule(cmd, rule, []string{"arg1", "arg2", "arg3"})
		assert.Error(t, err)

		err = validateRule(cmd, rule, []string{"arg1", "arg2"})
		assert.NoError(t, err)
	})

	t.Run("validates ExactArgs", func(t *testing.T) {
		rule := types.ArgRuleDef{ExactArgs: 2}
		err := validateRule(cmd, rule, []string{"arg1"})
		assert.Error(t, err)

		err = validateRule(cmd, rule, []string{"arg1", "arg2"})
		assert.NoError(t, err)
	})

	t.Run("validates RangeArgs with min and max", func(t *testing.T) {
		rule := types.ArgRuleDef{RangeArgs: []int{1, 3}}

		err := validateRule(cmd, rule, []string{})
		assert.Error(t, err)

		err = validateRule(cmd, rule, []string{"arg1"})
		assert.NoError(t, err)

		err = validateRule(cmd, rule, []string{"arg1", "arg2", "arg3"})
		assert.NoError(t, err)

		err = validateRule(cmd, rule, []string{"arg1", "arg2", "arg3", "arg4"})
		assert.Error(t, err)
	})

	t.Run("validates RangeArgs with only min", func(t *testing.T) {
		rule := types.ArgRuleDef{RangeArgs: []int{2}}

		err := validateRule(cmd, rule, []string{"arg1"})
		assert.Error(t, err)

		err = validateRule(cmd, rule, []string{"arg1", "arg2"})
		assert.NoError(t, err)

		err = validateRule(cmd, rule, []string{"arg1", "arg2", "arg3"})
		assert.NoError(t, err)
	})

	t.Run("validates ExactValidArgs", func(t *testing.T) {
		cmd.ValidArgs = []string{"valid1", "valid2"}
		rule := types.ArgRuleDef{ExactValidArgs: 2}

		err := validateRule(cmd, rule, []string{"valid1"})
		assert.Error(t, err)

		err = validateRule(cmd, rule, []string{"valid1", "valid2"})
		assert.NoError(t, err)

		err = validateRule(cmd, rule, []string{"invalid1", "invalid2"})
		assert.Error(t, err)
	})

	t.Run("validates Not rule", func(t *testing.T) {
		rule := types.ArgRuleDef{
			Not: &types.ArgRuleDef{
				ExactArgs: 2,
			},
		}

		err := validateRule(cmd, rule, []string{"arg1", "arg2"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed on `not`")

		err = validateRule(cmd, rule, []string{"arg1"})
		assert.NoError(t, err)
	})
}

func TestLogicalOperations(t *testing.T) {
	cmd := &cobra.Command{
		Use: "testcmd",
	}

	t.Run("validates MatchAll/And", func(t *testing.T) {
		rules := []types.ArgRuleDef{
			{MinimumNArgs: 1},
			{MaximumNArgs: 2},
		}
		err := validateAll(cmd, rules, []string{"arg1"})
		assert.NoError(t, err)

		err = validateAll(cmd, rules, []string{"arg1", "arg2", "arg3"})
		assert.Error(t, err)
	})

	t.Run("validates MatchAny/Or", func(t *testing.T) {
		rules := []types.ArgRuleDef{
			{ExactArgs: 1},
			{ExactArgs: 2},
		}
		err := validateAny(cmd, rules, []string{"arg1"})
		assert.NoError(t, err)

		err = validateAny(cmd, rules, []string{"arg1", "arg2"})
		assert.NoError(t, err)

		err = validateAny(cmd, rules, []string{"arg1", "arg2", "arg3"})
		assert.Error(t, err)
	})

	t.Run("validates MatchNone/Nand", func(t *testing.T) {
		rules := []types.ArgRuleDef{
			{ExactArgs: 1},
			{ExactArgs: 2},
		}
		err := validateNone(cmd, rules, []string{"arg1", "arg2", "arg3"})
		assert.NoError(t, err)

		err = validateNone(cmd, rules, []string{"arg1"})
		assert.Error(t, err)
	})
}

func TestValidateArgTypes(t *testing.T) {
	cmd := &cobra.Command{
		Use: "testcmd",
	}

	t.Run("validates required arguments", func(t *testing.T) {
		config := &types.ArgsConfig{
			Vars: []types.ArgVarDef{
				{Name: "required", Required: true},
				{Name: "optional"},
			},
		}

		// Test valid case
		store := CreateArgsStore(cmd, config, []string{"value", "optional"})
		err := ValidateArgs(cmd, config, store)
		assert.NoError(t, err)

		// Test missing required argument
		store = CreateArgsStore(cmd, config, []string{})
		err = ValidateArgs(cmd, config, store)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing required argument")
	})

	t.Run("validates type conversion", func(t *testing.T) {
		config := &types.ArgsConfig{
			Vars: []types.ArgVarDef{
				{Name: "number", Type: "number"},
				{Name: "boolean", Type: "boolean"},
			},
		}

		// Test valid values
		store := CreateArgsStore(cmd, config, []string{"42", "true"})
		err := ValidateArgs(cmd, config, store)
		assert.NoError(t, err)

		// Test invalid number
		store = CreateArgsStore(cmd, config, []string{"not-a-number", "true"})
		err = ValidateArgs(cmd, config, store)
		assert.Error(t, err)

		// Test invalid boolean
		store = CreateArgsStore(cmd, config, []string{"42", "not-a-boolean"})
		err = ValidateArgs(cmd, config, store)
		assert.Error(t, err)
	})

	t.Run("validates numeric constraints", func(t *testing.T) {
		minVal := 0.0
		maxVal := 10.0
		config := &types.ArgsConfig{
			Vars: []types.ArgVarDef{
				{
					Name: "value",
					Type: "number",
					Constraints: types.ParamConstraints{
						MinValue: &minVal,
						MaxValue: &maxVal,
					},
				},
			},
		}

		// Test valid value
		store := CreateArgsStore(cmd, config, []string{"5"})
		err := ValidateArgs(cmd, config, store)
		assert.NoError(t, err)

		// Test too small
		store = CreateArgsStore(cmd, config, []string{"-1"})
		err = ValidateArgs(cmd, config, store)
		assert.Error(t, err)

		// Test too large
		store = CreateArgsStore(cmd, config, []string{"11"})
		err = ValidateArgs(cmd, config, store)
		assert.Error(t, err)
	})

	t.Run("validates string constraints", func(t *testing.T) {
		config := &types.ArgsConfig{
			Vars: []types.ArgVarDef{
				{
					Name: "value",
					Type: "string",
					Constraints: types.ParamConstraints{
						MinLength: toPtr(3),
						MaxLength: toPtr(5),
						Pattern:   "^[a-z]+$",
					},
				},
			},
		}

		// Test valid value
		store := CreateArgsStore(cmd, config, []string{"abcd"})
		err := ValidateArgs(cmd, config, store)
		assert.NoError(t, err)

		// Test too short
		store = CreateArgsStore(cmd, config, []string{"ab"})
		err = ValidateArgs(cmd, config, store)
		assert.Error(t, err)

		// Test too long
		store = CreateArgsStore(cmd, config, []string{"abcdef"})
		err = ValidateArgs(cmd, config, store)
		assert.Error(t, err)

		// Test invalid pattern
		store = CreateArgsStore(cmd, config, []string{"123"})
		err = ValidateArgs(cmd, config, store)
		assert.Error(t, err)
	})

	t.Run("validates default values", func(t *testing.T) {
		config := &types.ArgsConfig{
			Vars: []types.ArgVarDef{
				{Name: "arg1", Default: "default1"},
				{Name: "arg2", Default: "default2"},
			},
		}

		// Test defaults applied
		store := CreateArgsStore(cmd, config, []string{})
		err := ValidateArgs(cmd, config, store)
		assert.NoError(t, err)
		assert.Equal(t, "default1", store.GetVal("arg1"))
		assert.Equal(t, "default2", store.GetVal("arg2"))

		// Test override defaults
		store = CreateArgsStore(cmd, config, []string{"value1", "value2"})
		err = ValidateArgs(cmd, config, store)
		assert.NoError(t, err)
		assert.Equal(t, "value1", store.GetVal("arg1"))
		assert.Equal(t, "value2", store.GetVal("arg2"))
	})

	t.Run("validates complex logical constraints", func(t *testing.T) {
		config := &types.ArgsConfig{
			Vars: []types.ArgVarDef{
				{
					Name: "value",
					Type: "number",
					Constraints: types.ParamConstraints{
						And: []*types.ParamConstraints{
							{
								Gt:  0.0,
								Lt:  10.0,
								Not: &types.ParamConstraints{Eq: 5.0},
							},
							{
								Or: []*types.ParamConstraints{
									{In: []any{2.0, 3.0, 4.0}},
									{In: []any{6.0, 7.0, 8.0}},
								},
							},
						},
					},
				},
			},
		}

		// Test valid value
		store := CreateArgsStore(cmd, config, []string{"3"})
		err := ValidateArgs(cmd, config, store)
		assert.NoError(t, err)

		// Test invalid value (equals 5)
		store = CreateArgsStore(cmd, config, []string{"5"})
		err = ValidateArgs(cmd, config, store)
		assert.Error(t, err)

		// Test invalid value (not in allowed sets)
		store = CreateArgsStore(cmd, config, []string{"1"})
		err = ValidateArgs(cmd, config, store)
		assert.Error(t, err)
	})

	t.Run("validates comparison constraints", func(t *testing.T) {
		config := &types.ArgsConfig{
			Vars: []types.ArgVarDef{
				{
					Name: "num",
					Type: "number",
					Constraints: types.ParamConstraints{
						Gt:    5.0,
						Lte:   10.0,
						NotIn: []any{7.0, 8.0},
					},
				},
			},
		}

		// Test valid value
		store := CreateArgsStore(cmd, config, []string{"9"})
		err := ValidateArgs(cmd, config, store)
		assert.NoError(t, err)

		// Test too small
		store = CreateArgsStore(cmd, config, []string{"5"})
		err = ValidateArgs(cmd, config, store)
		assert.Error(t, err)

		// Test excluded value
		store = CreateArgsStore(cmd, config, []string{"7"})
		err = ValidateArgs(cmd, config, store)
		assert.Error(t, err)
	})

	t.Run("validates file and path constraints", func(t *testing.T) {
		config := &types.ArgsConfig{
			Vars: []types.ArgVarDef{
				{
					Name: "file",
					Type: "string",
					Constraints: types.ParamConstraints{
						FileExists: "true",
					},
				},
				{
					Name: "dir",
					Type: "string",
					Constraints: types.ParamConstraints{
						DirExists: "true",
					},
				},
			},
		}

		// Test valid paths
		store := CreateArgsStore(cmd, config, []string{"validation_test.go", "."})
		err := ValidateArgs(cmd, config, store)
		assert.NoError(t, err)

		// Test invalid file
		store = CreateArgsStore(cmd, config, []string{"nonexistent.txt", "."})
		err = ValidateArgs(cmd, config, store)
		assert.Error(t, err)

		// Test invalid directory
		store = CreateArgsStore(cmd, config, []string{"validation_test.go", "nonexistent"})
		err = ValidateArgs(cmd, config, store)
		assert.Error(t, err)
	})
}

// Helper function to convert int to pointer
func toPtr(v int) *int {
	return &v
}
