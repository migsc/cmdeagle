package envvar

import (
	"embed"
	"strings"
)

//go:embed *
var PackageFS embed.FS

func GetEnvVariableNameFromStateKey(key string) string {
	var envKey string = key

	envKey = strings.ReplaceAll(envKey, ".", "_")
	envKey = strings.ReplaceAll(envKey, "[", "_")
	envKey = strings.ReplaceAll(envKey, "]", "")
	envKey = strings.ToUpper(envKey)

	return envKey
}
