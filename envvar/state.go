package envvar

import (
	"fmt"
	"strings"

	"github.com/migsc/cmdeagle/types"
)

type EnvStateStore struct {
	Entries map[string]string
}

func CreateEnvStore() *EnvStateStore {
	store := &EnvStateStore{
		Entries: map[string]string{},
	}

	return store
}

func (store *EnvStateStore) Set(key string, value string) {
	store.Entries[key] = value
}

func (store *EnvStateStore) Interpolate(script string) string {
	for key, val := range store.Entries {
		placeholder := fmt.Sprintf("{{%s}}", key)
		script = strings.ReplaceAll(script, placeholder, fmt.Sprint(val))
	}

	return script
}

func (store *EnvStateStore) GetEnvVariables() []types.EnvVar {
	envVars := make([]types.EnvVar, 0)

	for key, val := range store.Entries {
		envVars = append(envVars, types.EnvVar{Name: GetEnvVariableNameFromStateKey(key), Value: val})
	}

	return envVars
}

// func (store *EnvStateStore) ToJSON() map[string]any {
// 	return map[string]any{
// 		"args":  store.Args.ToJSON(),
// 		"flags": store.Flags.ToJSON(),
// 	}
// }

// func (store *ParamsStateStore) ToJSONString() string {
// 	jsonBytes, err := json.Marshal(store.ToJSON())
// 	if err != nil {
// 		return "{}"
// 	}
// 	return string(jsonBytes)
// }
