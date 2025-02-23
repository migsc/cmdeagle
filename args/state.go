package args

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/migsc/cmdeagle/envvar"
	"github.com/migsc/cmdeagle/types"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

type ArgsStateStore struct {
	CobraCommand *cobra.Command
	Config       *[]types.ArgDefinition
	Entries      map[string]*ArgStateEntry
	Count        int
	RawList      []string
}

type ArgStateEntry struct {
	Position int
	Def      *types.ArgDefinition
	Val      any
	RawVal   string
	// TODO: Thought about adding multiple error handling and that might help. Not implementing for now.
	Err error
}

var DefaultArgType string = "string"

// func CreateArgsStore(cobraCommand *cobra.Command, commandDef *config.CommandDefinition, args []string) *ArgsStateStore {
func CreateArgsStore(cobraCommand *cobra.Command, argsConfigDef *[]types.ArgDefinition, args []string) *ArgsStateStore {
	log.Debug("Creating args store / A", "args", args)
	store := &ArgsStateStore{
		CobraCommand: cobraCommand,
		Config:       argsConfigDef,
		Entries:      make(map[string]*ArgStateEntry),
		Count:        len(args),
		RawList:      args,
	}

	log.Debug("Creating args store / B", "args", args)
	if argsConfigDef == nil {
		return store
	}

	log.Debug("Creating args store / C", "args", args)
	for index, def := range *argsConfigDef {
		// Set default type if not specified
		if def.Type == "" {
			def.Type = DefaultArgType
		}
		argType := GetArgType(def.Type)

		// Determine the raw and converted values
		var val, rawVal any
		var err error

		log.Debug("Creating entry", "index", index, "def", def, "args", args)
		if index < len(args) {
			log.Debug("Handling provided argument", "index", index, "arg", args[index])
			// Handle provided argument
			rawVal = args[index]
			val, err = argType.Convert(args[index])
		} else {
			log.Debug("Handling default value", "index", index, "def", def, "default", def.Default)
			if def.Default != nil {
				// Handle default value
				rawVal = def.Default
				val = rawVal // Use default value as-is
			} else {
				// Handle no value case
				rawVal = argType.DefaultVal
				val = rawVal
			}

			if def.Required {
				err = fmt.Errorf("missing required argument: %s", def.Name)
			}
		}
		// Create entry
		entry := &ArgStateEntry{
			Position: index,
			Def:      &def,
			RawVal:   fmt.Sprint(rawVal),
			Val:      val,
			Err:      err,
		}

		// Store by both name and position
		store.Set(def.Name, entry)
		store.Set(fmt.Sprintf("list[%d]", index), entry)
	}

	return store
}

func (store *ArgsStateStore) Get(key string) *ArgStateEntry {
	return store.Entries[key]
}

func (store *ArgsStateStore) GetAt(index int) *ArgStateEntry {
	return store.Entries[fmt.Sprintf("list[%d]", index)]
}

func (store *ArgsStateStore) GetRawValAt(index int) string {
	if index >= len(store.RawList) {
		return ""
	}

	return store.RawList[index]
}

func (store *ArgsStateStore) GetVal(key string) any {
	if _, ok := store.Entries[key]; !ok {
		return nil
	}

	return store.Entries[key].Val
}

func (store *ArgsStateStore) GetAllVal() []any {
	vals := make([]any, store.Count)

	for _, entry := range store.Entries {
		vals = append(vals, entry.Val)
	}

	return vals
}

func (store *ArgsStateStore) GetAllRawVal() []string {
	vals := make([]string, store.Count)

	for _, entry := range store.Entries {
		vals = append(vals, entry.RawVal)
	}

	return vals
}

func (store *ArgsStateStore) Set(key string, entry *ArgStateEntry) *ArgStateEntry {
	if _, ok := store.Entries[key]; !ok {
		store.Entries[key] = entry
		store.Count++
	} else {
		store.Entries[key] = entry
	}

	return store.Entries[key]
}

func (store *ArgsStateStore) SetVal(key string, val any) any {
	if _, ok := store.Entries[key]; !ok {
		store.Entries[key] = &ArgStateEntry{
			Val: val,
			Def: nil,
		}
		store.Count++
	} else {
		store.Entries[key].Val = val
	}

	return store.Entries[key].Val
}

func (store *ArgsStateStore) Interpolate(script string) string {
	log.Debug("Interpolating", "script", script)
	for key, entry := range store.Entries {
		log.Debug("Interpolating", "key", key, "val", entry.Val)
		placeholder := fmt.Sprintf("${args.%s}", key)
		script = strings.ReplaceAll(script, placeholder, fmt.Sprint(entry.Val))
	}

	return script
}

func (store *ArgsStateStore) GetEnvVariables() []types.EnvVar {
	envVars := make([]types.EnvVar, 0)

	for key, entry := range store.Entries {
		envVars = append(envVars, types.EnvVar{Name: "ARGS_" + envvar.GetEnvVariableNameFromStateKey(key), Value: fmt.Sprint(entry.Val)})
	}

	return envVars
}

func (store *ArgsStateStore) ToJSON() map[string]any {
	result := make(map[string]any)

	// Add positional arguments as a list
	list := make([]any, 0)
	for i := 0; i < len(store.RawList); i++ {
		if entry := store.GetAt(i); entry != nil {
			list = append(list, entry.Val)
		} else {
			list = append(list, store.GetRawValAt(i))
		}
	}
	result["list"] = list

	// Add named arguments
	for key, entry := range store.Entries {
		// Only include named arguments
		if !strings.HasPrefix(key, "list[") {
			result[key] = entry.Val
		}
	}

	return result
}

func (store *ArgsStateStore) ToJSONString() string {
	jsonBytes, err := json.Marshal(store.ToJSON())
	if err != nil {
		return "{}"
	}
	return string(jsonBytes)
}
