package executable

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-version"
)

// DefaultVersionCommands is a predefined list of common flags and subcommands
var DefaultVersionCommands = [][]string{
	{"-v"}, {"--v"}, {"version"}, {"-version"}, {"--version"},
}

// ExtractSemver parses and validates a semantic version string (e.g., 1.2.3) from command output.
func ExtractSemver(output string) (string, error) {
	// Use go-version to parse and valifdate the semantic version
	v, err := version.NewVersion(output)
	if err != nil {
		return "", fmt.Errorf("no valid semantic version found in output: %s, error: %v", output, err)
	}
	return v.String(), nil
}

// GetVersionWithFallbacks attempts to retrieve a version using common flags and subcommands
func GetVersionWithFallbacks(binaryPath string) (string, error) {
	for _, cmdArgs := range DefaultVersionCommands {
		versionOutput, err := RunCommand(binaryPath, cmdArgs)
		if err != nil {
			continue // Swallow errors and try the next command
		}

		semver, err := ExtractSemver(versionOutput)
		if err == nil {
			return semver, nil // Return the first successfully parsed version
		}
	}

	return "", fmt.Errorf("unable to determine version for binary: %s", binaryPath)
}

// GetVersion retrieves the version of a binary using a user-specified command or falls back to defaults
func GetVersion(binaryPath string, customCommand []string) (string, error) {
	if len(customCommand) > 0 {
		// Use the custom version command if provided
		versionOutput, err := RunCommand(binaryPath, customCommand)
		if err != nil {
			return "", fmt.Errorf("custom version command failed: %v", err)
		}

		semver, err := ExtractSemver(versionOutput)
		if err == nil {
			return semver, nil
		}
		return "", fmt.Errorf("custom version command output could not be parsed")
	}

	// Fallback to default commands
	return GetVersionWithFallbacks(binaryPath)
}

// CheckVersionCompatibility determines if a found version satisfies version constraints
func CheckVersionCompatibility(foundVersion, declaredVersion string) (bool, string) {
	// Handle wildcard case first
	if declaredVersion == "*" {
		return true, ""
	}

	// Parse found version
	foundParts := strings.Split(strings.TrimPrefix(foundVersion, "v"), ".")
	if len(foundParts) < 3 {
		return false, "Invalid found version format"
	}
	foundMajor, foundMinor, foundPatch := foundParts[0], foundParts[1], foundParts[2]

	// Handle pre-release versions in found version
	foundPatch = strings.Split(foundPatch, "-")[0]

	// Parse declared version
	declaredVersion = strings.TrimSpace(strings.ToLower(declaredVersion))

	// Handle range operators
	switch {
	case strings.HasPrefix(declaredVersion, "^"):
		// Major version must match exactly, minor and patch can be greater
		declaredParts := strings.Split(strings.TrimPrefix(declaredVersion, "^"), ".")
		if len(declaredParts) < 3 {
			return false, "Invalid declared version format"
		}
		declaredMajor := declaredParts[0]

		if foundMajor != declaredMajor {
			return false, "Major version mismatch"
		}
		return true, ""

	case strings.HasPrefix(declaredVersion, "~"):
		// Major and minor must match exactly, only patch can be greater
		declaredParts := strings.Split(strings.TrimPrefix(declaredVersion, "~"), ".")
		if len(declaredParts) < 3 {
			return false, "Invalid declared version format"
		}
		declaredMajor, declaredMinor := declaredParts[0], declaredParts[1]

		if foundMajor != declaredMajor || foundMinor != declaredMinor {
			return false, "Major or minor version mismatch"
		}
		return true, ""

	case strings.HasPrefix(declaredVersion, ">="):
		// Version must be greater than or equal
		minVersion := strings.TrimPrefix(declaredVersion, ">=")
		minParts := strings.Split(minVersion, ".")
		if len(minParts) < 3 {
			return false, "Invalid minimum version format"
		}

		if !isVersionGreaterOrEqual(foundMajor, foundMinor, foundPatch, minParts[0], minParts[1], minParts[2]) {
			return false, "Version too low"
		}
		return true, ""

	case strings.HasPrefix(declaredVersion, ">"):
		// Version must be strictly greater
		minVersion := strings.TrimPrefix(declaredVersion, ">")
		minParts := strings.Split(minVersion, ".")
		if len(minParts) < 3 {
			return false, "Invalid minimum version format"
		}

		if !isVersionGreater(foundMajor, foundMinor, foundPatch, minParts[0], minParts[1], minParts[2]) {
			return false, "Version too low"
		}
		return true, ""

	case strings.HasPrefix(declaredVersion, "<="):
		// Version must be less than or equal
		maxVersion := strings.TrimPrefix(declaredVersion, "<=")
		maxParts := strings.Split(maxVersion, ".")
		if len(maxParts) < 3 {
			return false, "Invalid maximum version format"
		}

		if !isVersionLessOrEqual(foundMajor, foundMinor, foundPatch, maxParts[0], maxParts[1], maxParts[2]) {
			return false, "Version too high"
		}
		return true, ""

	case strings.HasPrefix(declaredVersion, "<"):
		// Version must be strictly less
		maxVersion := strings.TrimPrefix(declaredVersion, "<")
		maxParts := strings.Split(maxVersion, ".")
		if len(maxParts) < 3 {
			return false, "Invalid maximum version format"
		}

		if !isVersionLess(foundMajor, foundMinor, foundPatch, maxParts[0], maxParts[1], maxParts[2]) {
			return false, "Version too high"
		}
		return true, ""

	case strings.Contains(declaredVersion, " - "):
		// Handle version ranges
		parts := strings.Split(declaredVersion, " - ")
		if len(parts) != 2 {
			return false, "Invalid version range format"
		}

		minParts := strings.Split(parts[0], ".")
		maxParts := strings.Split(parts[1], ".")
		if len(minParts) < 3 || len(maxParts) < 3 {
			return false, "Invalid version range format"
		}

		if !isVersionGreaterOrEqual(foundMajor, foundMinor, foundPatch, minParts[0], minParts[1], minParts[2]) {
			return false, "Version below range"
		}
		if !isVersionLessOrEqual(foundMajor, foundMinor, foundPatch, maxParts[0], maxParts[1], maxParts[2]) {
			return false, "Version above range"
		}
		return true, ""

	default:
		// Exact version match
		declaredParts := strings.Split(declaredVersion, ".")
		if len(declaredParts) < 3 {
			return false, "Invalid declared version format"
		}

		if foundMajor != declaredParts[0] || foundMinor != declaredParts[1] || foundPatch != declaredParts[2] {
			return false, "Version does not match exactly"
		}
		return true, ""
	}
}

func isVersionGreaterOrEqual(foundMajor, foundMinor, foundPatch, targetMajor, targetMinor, targetPatch string) bool {
	// Convert to integers for proper numeric comparison
	fMajor, _ := strconv.Atoi(foundMajor)
	fMinor, _ := strconv.Atoi(foundMinor)
	fPatch, _ := strconv.Atoi(foundPatch)

	tMajor, _ := strconv.Atoi(targetMajor)
	tMinor, _ := strconv.Atoi(targetMinor)
	tPatch, _ := strconv.Atoi(targetPatch)

	if fMajor > tMajor {
		return true
	}
	if fMajor < tMajor {
		return false
	}
	if fMinor > tMinor {
		return true
	}
	if fMinor < tMinor {
		return false
	}
	return fPatch >= tPatch
}

func isVersionGreater(foundMajor, foundMinor, foundPatch, targetMajor, targetMinor, targetPatch string) bool {
	// Convert to integers for proper numeric comparison
	fMajor, _ := strconv.Atoi(foundMajor)
	fMinor, _ := strconv.Atoi(foundMinor)
	fPatch, _ := strconv.Atoi(foundPatch)

	tMajor, _ := strconv.Atoi(targetMajor)
	tMinor, _ := strconv.Atoi(targetMinor)
	tPatch, _ := strconv.Atoi(targetPatch)

	if fMajor > tMajor {
		return true
	}
	if fMajor < tMajor {
		return false
	}
	if fMinor > tMinor {
		return true
	}
	if fMinor < tMinor {
		return false
	}
	return fPatch > tPatch
}

func isVersionLessOrEqual(foundMajor, foundMinor, foundPatch, targetMajor, targetMinor, targetPatch string) bool {
	// Convert to integers for proper numeric comparison
	fMajor, _ := strconv.Atoi(foundMajor)
	fMinor, _ := strconv.Atoi(foundMinor)
	fPatch, _ := strconv.Atoi(foundPatch)

	tMajor, _ := strconv.Atoi(targetMajor)
	tMinor, _ := strconv.Atoi(targetMinor)
	tPatch, _ := strconv.Atoi(targetPatch)

	if fMajor < tMajor {
		return true
	}
	if fMajor > tMajor {
		return false
	}
	if fMinor < tMinor {
		return true
	}
	if fMinor > tMinor {
		return false
	}
	return fPatch <= tPatch
}

func isVersionLess(foundMajor, foundMinor, foundPatch, targetMajor, targetMinor, targetPatch string) bool {
	// Convert to integers for proper numeric comparison
	fMajor, _ := strconv.Atoi(foundMajor)
	fMinor, _ := strconv.Atoi(foundMinor)
	fPatch, _ := strconv.Atoi(foundPatch)

	tMajor, _ := strconv.Atoi(targetMajor)
	tMinor, _ := strconv.Atoi(targetMinor)
	tPatch, _ := strconv.Atoi(targetPatch)

	if fMajor < tMajor {
		return true
	}
	if fMajor > tMajor {
		return false
	}
	if fMinor < tMinor {
		return true
	}
	if fMinor > tMinor {
		return false
	}
	return fPatch < tPatch
}
