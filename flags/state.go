package flags

import (
	"cmdeagle/envvar"
	"cmdeagle/types"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type FlagsStateStore struct {
	// TODO implement global flags
	// TODO implement grouped flag configuration
	cobraCommand *cobra.Command
	pFlagSet     *pflag.FlagSet
	flagDefMap   map[string]*types.FlagDefinition
}

func CreateFlagsStore(cobraCommand *cobra.Command, commandDef *types.CommandDefinition) *FlagsStateStore {
	// TODO handle persistent flags
	// https://cobra.dev/#persistent-flags

	store := &FlagsStateStore{
		cobraCommand: cobraCommand,
		pFlagSet:     cobraCommand.Flags(),
		flagDefMap:   make(map[string]*types.FlagDefinition),
	}

	for _, flagDef := range commandDef.Flags {
		log.Debug("\tGetting flag definition", "name", flagDef.Name)
		flagType := GetFlagType(flagDef.Type)

		log.Debug("\tBinding flag", "name", flagDef.Name)
		var flagVal *any
		flagType.Bind(flagVal, store.pFlagSet, &flagDef)

		store.flagDefMap[flagDef.Name] = &flagDef
	}

	return store
}

func (store *FlagsStateStore) Get(key string) *pflag.Flag {
	return store.pFlagSet.Lookup(key)
}

func (store *FlagsStateStore) GetVal(key string) any {
	flag := store.pFlagSet.Lookup(key)
	if flag == nil {
		return nil
	}

	return flag.Value
}

func (store *FlagsStateStore) GetDef(key string) *types.FlagDefinition {
	return store.flagDefMap[key]
}

func (store *FlagsStateStore) VisitAll(fn func(flag *pflag.Flag)) {
	store.pFlagSet.VisitAll(fn)
}

func (store *FlagsStateStore) Interpolate(script string) string {
	store.pFlagSet.VisitAll(func(flag *pflag.Flag) {
		placeholder := fmt.Sprintf("${flags.%s}", flag.Name)
		script = strings.ReplaceAll(script, placeholder, fmt.Sprint(flag.Value))
	})

	return script
}

func (store *FlagsStateStore) GetEnvVariables() []types.EnvVar {
	envVars := make([]types.EnvVar, 0)

	store.pFlagSet.VisitAll(func(flag *pflag.Flag) {
		envVars = append(envVars, types.EnvVar{Name: "FLAGS_" + envvar.GetEnvVariableNameFromStateKey(flag.Name), Value: fmt.Sprint(flag.Value)})
	})

	return envVars
}

func (store *FlagsStateStore) ToJSON() map[string]any {
	result := make(map[string]any)

	store.pFlagSet.VisitAll(func(flag *pflag.Flag) {
		result[flag.Name] = flag.Value.String()
	})

	return result
}

func (store *FlagsStateStore) ToJSONString() string {
	jsonBytes, err := json.Marshal(store.ToJSON())
	if err != nil {
		return "{}"
	}
	return string(jsonBytes)
}
