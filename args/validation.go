package args

import (
	"errors"
	"fmt"

	"github.com/migsc/cmdeagle/params"
	"github.com/migsc/cmdeagle/types"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

func ValidateArgs(cobraCmd *cobra.Command, argsConfigDef *types.ArgsConfig, store *ArgsStateStore) error {
	// log.Debug("Validating args", "argsConfigDef", argsConfigDef, "store", store)
	if argsConfigDef == nil || store == nil {
		return nil
	}

	rawArgs := store.GetAllRawVal()

	if argsConfigDef.Vars != nil {

		for index := range argsConfigDef.Vars {

			entry := store.GetAt(index)

			log.Debug("Validating arg", "index", index, "name", entry.Def.Name, "rawVal", entry.RawVal, "val", entry.Val, "err", entry.Err)

			if entry == nil {
				continue
			}

			if entry.Err != nil {
				return entry.Err
			}

			// TODO: Gonna need to refactor the constraints to be pointer based for this to work
			// if entry.Def != nil && entry.Def.Constraints != nil {
			// 	for _, dependency := range entry.Def.DependsOn {
			// 		err := ValidateConstraint(dependency.When, store.GetVal(dependency.Name))
			// 		if err != nil {
			// 			return err
			// 		}
			// 	}
			// }

			if entry.Def != nil && entry.Def.DependsOn != nil {
				for _, dependency := range entry.Def.DependsOn {
					err := params.ValidateConstraint(dependency.When, store.GetVal(dependency.Name))
					if err != nil {
						return err
					}
				}
			}

			if entry.Def != nil && entry.Def.ConflictsWith != nil {
				for _, conflict := range entry.Def.ConflictsWith {
					conflictVal := store.GetVal(conflict)
					if conflictVal != nil && conflictVal != "" {
						return fmt.Errorf("argument %s conflicts with %s", entry.Def.Name, conflict)
					}
				}
			}
		}
	}

	if argsConfigDef.Rules != nil {
		for _, ruleDef := range argsConfigDef.Rules {
			err := validateRule(cobraCmd, ruleDef, rawArgs)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func validateRule(cobraCmd *cobra.Command, ruleDef types.ArgRuleDef, args []string) error {
	if ruleDef.NoArgs {
		if len(args) > 0 {
			return fmt.Errorf("%s accepts no args", cobraCmd.Name())
		}
	}

	if ruleDef.OnlyValidArgs {
		err := cobra.OnlyValidArgs(cobraCmd, args)
		if err != nil {
			return err
		}
	}

	if ruleDef.ArbitraryArgs {
		err := cobra.ArbitraryArgs(cobraCmd, args)
		if err != nil {
			return err
		}
	}

	if ruleDef.MinimumNArgs > 0 {
		err := cobra.MinimumNArgs(ruleDef.MinimumNArgs)(cobraCmd, args)
		if err != nil {
			return err
		}
	}

	if ruleDef.MaximumNArgs > 0 {
		err := cobra.MaximumNArgs(ruleDef.MaximumNArgs)(cobraCmd, args)
		if err != nil {
			return err
		}
	}

	if ruleDef.ExactArgs > 0 {
		err := cobra.ExactArgs(ruleDef.ExactArgs)(cobraCmd, args)
		if err != nil {
			return err
		}
	}

	if len(ruleDef.RangeArgs) > 1 {
		err := cobra.RangeArgs(ruleDef.RangeArgs[0], ruleDef.RangeArgs[1])(cobraCmd, args)
		if err != nil {
			return err
		}
	} else if len(ruleDef.RangeArgs) == 1 {
		err := cobra.MinimumNArgs(ruleDef.RangeArgs[0])(cobraCmd, args)
		if err != nil {
			return err
		}
	}

	if ruleDef.ExactValidArgs > 0 {
		err := cobra.ExactArgs(ruleDef.ExactValidArgs)(cobraCmd, args)
		if err != nil {
			return err
		}
	}

	if ruleDef.MatchAll != nil {
		err := validateAll(cobraCmd, *ruleDef.MatchAll, args)
		if err != nil {
			return err
		}
	}

	if ruleDef.And != nil {
		err := validateAll(cobraCmd, *ruleDef.And, args)
		if err != nil {
			return err
		}
	}

	if ruleDef.MatchAny != nil {
		err := validateAny(cobraCmd, *ruleDef.MatchAny, args)
		if err != nil {
			return err
		}
	}

	if ruleDef.Or != nil {
		err := validateAny(cobraCmd, *ruleDef.Or, args)
		if err != nil {
			return err
		}
	}

	if ruleDef.MatchNone != nil {
		err := validateNone(cobraCmd, *ruleDef.MatchNone, args)
		if err != nil {
			return err
		}
	}

	if ruleDef.Nand != nil {
		err := validateNone(cobraCmd, *ruleDef.Nand, args)
		if err != nil {
			return err
		}
	}

	if ruleDef.Not != nil {
		err := validateRule(cobraCmd, ruleDef, args)
		if err == nil {
			return errors.New(fmt.Sprintf("Validation failed on `not` for rule: %v", ruleDef))
		}

	}

	return nil
}

func validateAll(cobraCmd *cobra.Command, ruleDefs []types.ArgRuleDef, args []string) error {
	for _, rule := range ruleDefs {
		err := validateRule(cobraCmd, rule, args)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateAny(cobraCmd *cobra.Command, ruleDefs []types.ArgRuleDef, args []string) error {
	var firstErrorFound error
	atLeastOneValid := false
	for _, rule := range ruleDefs {

		err := validateRule(cobraCmd, rule, args)
		if err == nil {
			atLeastOneValid = true
		} else if firstErrorFound == nil {
			firstErrorFound = err
		}
	}

	if !atLeastOneValid {
		return firstErrorFound
	}

	return nil
}

func validateNone(cobraCmd *cobra.Command, ruleDefs []types.ArgRuleDef, args []string) error {
	atLeastOneValid := false
	for _, rule := range ruleDefs {

		err := validateRule(cobraCmd, rule, args)
		if err == nil {
			atLeastOneValid = true
		}
	}

	if atLeastOneValid {
		return fmt.Errorf("Validation failed when asserting none were valid for rule: %v", ruleDefs)
	}

	return nil
}
