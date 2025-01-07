package args

import (
	"embed"
	"fmt"
	"time"

	cast "github.com/spf13/cast"
)

//go:embed *
var PackageFS embed.FS

type ArgTypeDef struct {
	DefaultVal any
	Convert    func(val string) (any, error)
}

var ArgRuleConditionals = []string{"MatchAll", "MatchAny", "MatchNone", "And", "Or", "Nand", "Not"}

func AddArgType(name string, defaultVal any, convert func(val string) (any, error)) {
	if _, ok := argTypes[name]; ok {
		panic(fmt.Sprintf("Arg type `%s` already exists", name))
	}

	argTypes[name] = ArgTypeDef{
		Convert:    convert,
		DefaultVal: defaultVal,
	}
}

func GetArgType(name string) ArgTypeDef {
	argType, ok := argTypes[name]
	if !ok {
		panic(fmt.Sprintf("Arg type `%s` not found", name))
	}

	return argType
}

var argTypes = map[string]ArgTypeDef{
	"string": {
		DefaultVal: "",
		Convert:    func(val string) (any, error) { return cast.ToStringE(val) },
	},
	"date": {
		DefaultVal: time.Now(),
		Convert:    func(val string) (any, error) { return cast.StringToDate(val) },
	},
	"time": {
		DefaultVal: time.Now(),
		Convert:    func(val string) (any, error) { return cast.ToTimeE(val) },
	},
	"duration": {
		DefaultVal: time.Duration(0),
		Convert:    func(val string) (any, error) { return cast.ToDurationE(val) },
	},
	"boolean": {
		DefaultVal: false,
		Convert:    func(val string) (any, error) { return cast.ToBoolE(val) },
	},
	"bool": {
		DefaultVal: false,
		Convert:    func(val string) (any, error) { return cast.ToBoolE(val) },
	},
	"number": {
		DefaultVal: 0.0,
		Convert:    func(val string) (any, error) { return cast.ToFloat64E(val) },
	},
	"float64": {
		DefaultVal: 0.0,
		Convert:    func(val string) (any, error) { return cast.ToFloat64E(val) },
	},
	"float32": {
		DefaultVal: 0.0,
		Convert:    func(val string) (any, error) { return cast.ToFloat32E(val) },
	},
	"int64": {
		DefaultVal: int64(0),
		Convert:    func(val string) (any, error) { return cast.ToInt64E(val) },
	},
	"int32": {
		DefaultVal: int32(0),
		Convert:    func(val string) (any, error) { return cast.ToInt32E(val) },
	},
	"int16": {
		DefaultVal: int16(0),
		Convert:    func(val string) (any, error) { return cast.ToInt16E(val) },
	},
	"int8": {
		DefaultVal: int8(0),
		Convert:    func(val string) (any, error) { return cast.ToInt8E(val) },
	},
	"int": {
		DefaultVal: int(0),
		Convert:    func(val string) (any, error) { return cast.ToIntE(val) },
	},
	"uint": {
		DefaultVal: uint(0),
		Convert:    func(val string) (any, error) { return cast.ToUintE(val) },
	},
	"uint64": {
		DefaultVal: uint64(0),
		Convert:    func(val string) (any, error) { return cast.ToUint64E(val) },
	},
	"uint32": {
		DefaultVal: uint32(0),
		Convert:    func(val string) (any, error) { return cast.ToUint32E(val) },
	},
	"uint16": {
		DefaultVal: uint16(0),
		Convert:    func(val string) (any, error) { return cast.ToUint16E(val) },
	},
	"uint8": {
		DefaultVal: uint8(0),
		Convert:    func(val string) (any, error) { return cast.ToUint8E(val) },
	},
}
