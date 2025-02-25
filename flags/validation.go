package flags

import (
	"fmt"
	"regexp"

	"github.com/charmbracelet/log"
	"github.com/migsc/cmdeagle/params"
	"github.com/migsc/cmdeagle/types"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// TODO: We ened to do this differently because we depend on the pflag library to store the flags
// We must iterate through the flags with flagset's VisitAll function and then validate each one
func ValidateFlags(cobraCmd *cobra.Command, flagsConfigDefs []types.FlagDefinition, store *FlagsStateStore) error {
	// // First validate individual flag definitions
	// for _, flagDef := range flagsConfigDefs {
	// 	entry := store.Get(fmt.Sprintf("flags.%s", flagDef.Name))
	// 	if entry == nil {
	// 		if flagDef.Required {
	// 			return fmt.Errorf("required flag '%s' not provided", flagDef.Name)
	// 		}
	// 		continue
	// 	}

	// 	if entry.err != nil {
	// 		return entry.err
	// 	}
	// }

	// TODO: Then validate flag constraints

	for _, flagDef := range flagsConfigDefs {
		flag := store.Get(flagDef.Name)
		if flag == nil {
			continue
		}

		err := params.ValidateConstraint(flagDef.Constraints, flag.Value)
		if err != nil {
			return err
		}
	}

	var foundErr error
	store.VisitAll(func(flag *pflag.Flag) {

		flagDef := store.GetDef(flag.Name)

		// TODO Might want to handle extra flags on a setting that enables strictness
		if flagDef == nil {
			return
		}

		if flagDef.DependsOn != nil {
			for _, dependency := range flagDef.DependsOn {
				err := params.ValidateConstraint(dependency.When, store.GetVal(dependency.Name))
				if err != nil {
					foundErr = err
				}
			}
		}

		if flagDef.ConflictsWith != nil {
			for _, conflict := range flagDef.ConflictsWith {
				otherFlag := store.Get(conflict)

				// Only check for conflicts if both flags were explicitly set by the user
				if otherFlag != nil && otherFlag.Changed && flag.Changed {
					foundErr = fmt.Errorf("Argument %s conflicts with %s", flagDef.Name, conflict)
				}
			}
		}

		if flagDef.Constraints != nil {
			err := params.ValidateConstraint(flagDef.Constraints, flag.Value)
			if err != nil {
				foundErr = err
			}
		}

		// Validate pattern
		// Validate pattern
		if flagDef.Pattern != "" {
			var err error
			var match bool
			var pattern *regexp.Regexp

			pattern, err = regexp.Compile(flagDef.Pattern)

			if err != nil {
				foundErr = fmt.Errorf("invalid pattern for argument %s: %v", flagDef.Name, err)
			}
			log.Debug("Validating pattern for argument", "pattern", pattern, "value", flag.Value)
			match = pattern.MatchString(flag.Value.String())

			log.Debug("Validating pattern for argument", "pattern", pattern, "value", flag.Value, "match", match, "err", err)

			if !match {
				foundErr = fmt.Errorf("pattern validation failed for argument %s: %v", flagDef.Name, flag.Value)
			}
		}

	})

	return foundErr
}
