package config

import (
	"cmdeagle/args"
	"cmdeagle/envvar"
	"cmdeagle/flags"
	"cmdeagle/types"
	"encoding/json"
	"fmt"
	"strings"
)

type ParamsStateStore struct {
	Args    *args.ArgsStateStore
	Flags   *flags.FlagsStateStore
	Entries map[string]string
}

func CreateParamsStore(argsStore *args.ArgsStateStore, flagsStore *flags.FlagsStateStore) *ParamsStateStore {
	store := &ParamsStateStore{
		Args:    argsStore,
		Flags:   flagsStore,
		Entries: map[string]string{},
	}

	store.Entries["args.json"] = store.Args.ToJSONString()
	store.Entries["flags.json"] = store.Flags.ToJSONString()
	store.Entries["params.json"] = store.ToJSONString()

	return store
}

func (store *ParamsStateStore) Interpolate(script string) string {
	for key, val := range store.Entries {
		placeholder := fmt.Sprintf("${%s}", key)
		script = strings.ReplaceAll(script, placeholder, fmt.Sprint(val))
	}

	return script
}

func (store *ParamsStateStore) GetEnvVariables() []types.EnvVar {
	envVars := make([]types.EnvVar, 0)

	envVars = append(envVars, store.Args.GetEnvVariables()...)
	envVars = append(envVars, store.Flags.GetEnvVariables()...)

	for key, val := range store.Entries {
		envVars = append(envVars, types.EnvVar{Name: envvar.GetEnvVariableNameFromStateKey(key), Value: val})
	}

	return envVars
}

func (store *ParamsStateStore) ToJSON() map[string]any {
	return map[string]any{
		"args":  store.Args.ToJSON(),
		"flags": store.Flags.ToJSON(),
	}
}

func (store *ParamsStateStore) ToJSONString() string {
	jsonBytes, err := json.Marshal(store.ToJSON())
	if err != nil {
		return "{}"
	}
	return string(jsonBytes)
}
