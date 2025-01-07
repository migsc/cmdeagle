package flags

import (
	"embed"
	"fmt"

	"github.com/migsc/cmdeagle/types"

	"github.com/spf13/pflag"
)

//go:embed *
var PackageFS embed.FS

type FlagTypeDef struct {
	Bind func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any
}

func GetFlagType(name string) FlagTypeDef {
	flagType, ok := flagTypes[name]
	if !ok {
		panic(fmt.Sprintf("Flag type `%s` not found", name))
	}

	return flagType
}

// TODO: There's not really any error handling here. We should proably use the cast E functions to validate the values and return errors

var flagTypes = map[string]FlagTypeDef{
	"string": {
		Bind: func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
			var strVal string
			var flagVal any = &strVal
			defaultVal := ""
			if flagDef.Default != nil {
				defaultVal = flagDef.Default.(string)
			}
			flagSet.StringVarP(&strVal, flagDef.Name, flagDef.Shorthand, defaultVal, flagDef.Description)
			return &flagVal
		},
	},
	"boolean": {
		Bind: func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
			var boolVal bool
			var flagVal any = &boolVal
			defaultVal := false
			if flagDef.Default != nil {
				defaultVal = flagDef.Default.(bool)
			}
			flagSet.BoolVarP(&boolVal, flagDef.Name, flagDef.Shorthand, defaultVal, flagDef.Description)
			return &flagVal
		},
	},
	"bool": {
		Bind: func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
			var boolVal bool
			var flagVal any = &boolVal
			defaultVal := false
			if flagDef.Default != nil {
				defaultVal = flagDef.Default.(bool)
			}
			flagSet.BoolVarP(&boolVal, flagDef.Name, flagDef.Shorthand, defaultVal, flagDef.Description)
			return &flagVal
		},
	},
	"number": {
		Bind: func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
			var numVal float64
			var flagVal any = &numVal
			defaultVal := float64(0)
			if flagDef.Default != nil {
				defaultVal = flagDef.Default.(float64)
			}
			flagSet.Float64VarP(&numVal, flagDef.Name, flagDef.Shorthand, defaultVal, flagDef.Description)
			return &flagVal
		},
	},
	"float64": {
		Bind: func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
			var numVal float64
			var flagVal any = &numVal
			defaultVal := float64(0)
			if flagDef.Default != nil {
				defaultVal = flagDef.Default.(float64)
			}
			flagSet.Float64VarP(&numVal, flagDef.Name, flagDef.Shorthand, defaultVal, flagDef.Description)
			return &flagVal
		},
	},
	"float32": {
		Bind: func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
			var numVal float32
			var flagVal any = &numVal
			defaultVal := float32(0)
			if flagDef.Default != nil {
				defaultVal = flagDef.Default.(float32)
			}
			flagSet.Float32VarP(&numVal, flagDef.Name, flagDef.Shorthand, defaultVal, flagDef.Description)
			return &flagVal
		},
	},
	"int64": {
		Bind: func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
			var numVal int64
			var flagVal any = &numVal
			defaultVal := int64(0)
			if flagDef.Default != nil {
				defaultVal = flagDef.Default.(int64)
			}
			flagSet.Int64VarP(&numVal, flagDef.Name, flagDef.Shorthand, defaultVal, flagDef.Description)
			return &flagVal
		},
	},
	"int32": {
		Bind: func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
			var numVal int32
			var flagVal any = &numVal
			defaultVal := int32(0)
			if flagDef.Default != nil {
				defaultVal = flagDef.Default.(int32)
			}
			flagSet.Int32VarP(&numVal, flagDef.Name, flagDef.Shorthand, defaultVal, flagDef.Description)
			return &flagVal
		},
	},
	"int16": {
		Bind: func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
			var numVal int16
			var flagVal any = &numVal
			defaultVal := int16(0)
			if flagDef.Default != nil {
				defaultVal = flagDef.Default.(int16)
			}
			flagSet.Int16VarP(&numVal, flagDef.Name, flagDef.Shorthand, defaultVal, flagDef.Description)
			return &flagVal
		},
	},
	"int8": {
		Bind: func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
			var numVal int8
			var flagVal any = &numVal
			defaultVal := int8(0)
			if flagDef.Default != nil {
				defaultVal = flagDef.Default.(int8)
			}
			flagSet.Int8VarP(&numVal, flagDef.Name, flagDef.Shorthand, defaultVal, flagDef.Description)
			return &flagVal
		},
	},
	"int": {
		Bind: func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
			var numVal int
			var flagVal any = &numVal
			defaultVal := 0
			if flagDef.Default != nil {
				defaultVal = flagDef.Default.(int)
			}
			flagSet.IntVarP(&numVal, flagDef.Name, flagDef.Shorthand, defaultVal, flagDef.Description)
			return &flagVal
		},
	},
	"uint": {
		Bind: func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
			var numVal uint
			var flagVal any = &numVal
			defaultVal := uint(0)
			if flagDef.Default != nil {
				defaultVal = flagDef.Default.(uint)
			}
			flagSet.UintVarP(&numVal, flagDef.Name, flagDef.Shorthand, defaultVal, flagDef.Description)
			return &flagVal
		},
	},
	"uint64": {
		Bind: func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
			var numVal uint64
			var flagVal any = &numVal
			defaultVal := uint64(0)
			if flagDef.Default != nil {
				defaultVal = flagDef.Default.(uint64)
			}
			flagSet.Uint64VarP(&numVal, flagDef.Name, flagDef.Shorthand, defaultVal, flagDef.Description)
			return &flagVal
		},
	},
	"uint32": {
		Bind: func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
			var numVal uint32
			var flagVal any = &numVal
			defaultVal := uint32(0)
			if flagDef.Default != nil {
				defaultVal = flagDef.Default.(uint32)
			}
			flagSet.Uint32VarP(&numVal, flagDef.Name, flagDef.Shorthand, defaultVal, flagDef.Description)
			return &flagVal
		},
	},
	"uint16": {
		Bind: func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
			var numVal uint16
			var flagVal any = &numVal
			defaultVal := uint16(0)
			if flagDef.Default != nil {
				defaultVal = flagDef.Default.(uint16)
			}
			flagSet.Uint16VarP(&numVal, flagDef.Name, flagDef.Shorthand, defaultVal, flagDef.Description)
			return &flagVal
		},
	},
	"uint8": {
		Bind: func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
			var numVal uint8
			var flagVal any = &numVal
			defaultVal := uint8(0)
			if flagDef.Default != nil {
				defaultVal = flagDef.Default.(uint8)
			}
			flagSet.Uint8VarP(&numVal, flagDef.Name, flagDef.Shorthand, defaultVal, flagDef.Description)
			return &flagVal
		},
	},
}
