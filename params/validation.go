package params

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/migsc/cmdeagle/types"

	afero "github.com/spf13/afero"
	cast "github.com/spf13/cast"
)

// Returns a boolean and a reason in the case of failing the constraints
func ValidateConstraint(constraints *types.ParamConstraints, value any, useMemMapFs ...bool) error {

	if constraints == nil {
		return nil
	}

	if constraints.MinValue != nil {
		if cast.ToInt(value) < cast.ToInt(constraints.MinValue) {
			return fmt.Errorf("Value is less than the minimum value of %v", constraints.MinValue)
		}
	}

	if constraints.MaxValue != nil {
		if cast.ToInt(value) > cast.ToInt(constraints.MaxValue) {
			return fmt.Errorf("Value is greater than the maximum value of %v", constraints.MaxValue)
		}
	}

	if constraints.MultipleOf != nil {
		if cast.ToInt(value)%cast.ToInt(constraints.MultipleOf) != 0 {
			return fmt.Errorf("Value is not a multiple of %v", constraints.MultipleOf)
		}
	}

	if constraints.Eq != nil {
		if value != constraints.Eq {
			return fmt.Errorf("Value is not equal to %v", constraints.Eq)
		}
	}

	if constraints.Neq != nil {
		if value == constraints.Neq {
			return fmt.Errorf("Value is equal to %v", constraints.Neq)
		}
	}

	if constraints.Gt != nil {
		if cast.ToInt(value) <= cast.ToInt(constraints.Gt) {
			return fmt.Errorf("Value is not greater than %v", constraints.Gt)
		}
	}

	if constraints.Gte != nil {
		if cast.ToInt(value) < cast.ToInt(constraints.Gte) {
			return fmt.Errorf("Value is not greater than or equal to %v", constraints.Gte)
		}
	}

	if constraints.Lt != nil {
		if cast.ToInt(value) >= cast.ToInt(constraints.Lt) {
			return fmt.Errorf("Value is not less than %v", constraints.Lt)
		}
	}

	if constraints.Lte != nil {
		if cast.ToInt(value) > cast.ToInt(constraints.Lte) {
			return fmt.Errorf("Value is not less than or equal to %v", constraints.Lte)
		}
	}

	if constraints.In != nil {
		if !slices.Contains(constraints.In, value) {
			return fmt.Errorf("Value is not in the list of %v", constraints.In)
		}
	}

	if constraints.NotIn != nil {
		if slices.Contains(constraints.NotIn, value) {
			return fmt.Errorf("Value is in the list of %v", constraints.NotIn)
		}
	}

	if constraints.MinLength != nil {
		if len(cast.ToString(value)) < cast.ToInt(constraints.MinLength) {
			return fmt.Errorf("Value is less than the minimum character length of %v", constraints.MinLength)
		}
	}

	if constraints.MaxLength != nil {
		if len(cast.ToString(value)) > cast.ToInt(constraints.MaxLength) {
			return fmt.Errorf("Value is greater than the maximum character length of %v", constraints.MaxLength)
		}
	}

	if constraints.Pattern != "" {
		if !regexp.MustCompile(constraints.Pattern).MatchString(cast.ToString(value)) {
			return fmt.Errorf("Value does not match the pattern %v", constraints.Pattern)
		}
	}

	if HasFileConstraints(constraints) {
		var fs afero.Fs

		if useMemMapFs != nil && useMemMapFs[0] {
			fs = afero.NewMemMapFs()
		} else {
			fs = afero.NewOsFs()
		}

		err := validateFileConstraints(fs, constraints, value)
		if err != nil {
			return err
		}
	}

	if constraints.And != nil {
		for _, constraint := range constraints.And {
			err := ValidateConstraint(constraint, value)
			if err != nil {
				return err
			}
		}
	}

	if constraints.Nand != nil {
		for _, constraint := range constraints.Nand {
			err := ValidateConstraint(constraint, value)
			if err == nil {
				return fmt.Errorf("Validation failed on `nand` for constraint %v", constraint)
			}
		}
	}

	if constraints.Or != nil {
		var firstErrorFound error
		atLeastOneValid := false
		for _, constraint := range constraints.Or {

			err := ValidateConstraint(constraint, value)

			if err == nil {
				atLeastOneValid = true
			} else if firstErrorFound == nil {
				firstErrorFound = err
			}
		}

		if !atLeastOneValid {
			return firstErrorFound
		}
	}

	if constraints.Not != nil {
		err := ValidateConstraint(constraints.Not, value)
		if err == nil {
			return fmt.Errorf("Condition is true `not` constraint of %v", constraints.Not)
		}
	}

	return nil
}

func HasFileConstraints(constraints *types.ParamConstraints) bool {
	for _, constraint := range types.ConstraintFileKeys {
		if value := reflect.ValueOf(constraints).FieldByName(constraint).String(); value != "" {
			return true
		}
	}
	return false
}

// TODO delete this. I guess I didn't actually need it after all.
// func HasConditionalConstraints(constraints schema.ParamConstraints) bool {
// 	for _, constraint := range schema.ConstraintConditionals {
// 		if value := reflect.ValueOf(constraints).FieldByName(constraint).String(); value != "" {
// 			return true
// 		}
// 	}
// 	return false
// }

func validateFileConstraints(fs afero.Fs, constraint *types.ParamConstraints, value any) error {
	filePath := cast.ToString(value)

	// Check DirExists constraint
	if constraint.DirExists != "" {
		isDir, err := afero.IsDir(fs, filePath)
		if err != nil || !isDir {
			return fmt.Errorf("Value is not a directory: %v", value)
		}
	}

	// Check FileExists constraint
	if constraint.FileExists != "" {
		exists, err := afero.Exists(fs, filePath)
		if err != nil || !exists {
			return fmt.Errorf("File does not exist: %v", filePath)
		}

		// If we're only checking existence, we can return here
		if constraint.IsFileType == "" && constraint.HasPermissions == "" {
			return nil
		}

		// Verify it's not a directory
		isDir, err := afero.IsDir(fs, filePath)
		if err != nil || isDir {
			return fmt.Errorf("Path is a directory, not a file: %v", filePath)
		}
	}

	// Check HasPermissions constraint
	if constraint.HasPermissions != "" {
		info, err := fs.Stat(filePath)
		if err != nil {
			return fmt.Errorf("Cannot check file permissions: %v", err)
		}

		// Convert permission string to os.FileMode
		// Expected format is Unix-style octal (e.g., "0644")
		wantPerm, err := strconv.ParseInt(constraint.HasPermissions, 8, 32)
		if err != nil {
			return fmt.Errorf("Invalid permission format: %v", constraint.HasPermissions)
		}

		if info.Mode().Perm() != os.FileMode(wantPerm) {
			return fmt.Errorf("File has incorrect permissions. Want: %v, Got: %v",
				os.FileMode(wantPerm), info.Mode().Perm())
		}
	}

	// Check IsFileType constraint
	if constraint.IsFileType != "" {
		if err := validateFileType(fs, filePath, constraint.IsFileType); err != nil {
			return err
		}
	}

	return nil
}

func validateFileType(fs afero.Fs, filePath string, expectedType string) error {
	// Normalize expected type
	if !strings.HasPrefix(expectedType, ".") && !strings.Contains(expectedType, "/") {
		expectedType = "." + expectedType
	}

	// Open and read file for MIME type detection
	file, err := fs.Open(filePath)
	if err != nil {
		return fmt.Errorf("Cannot open file: %v", err)
	}
	defer file.Close()

	// Read first 512 bytes for MIME type detection
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("Cannot read file: %v", err)
	}

	detectedType := http.DetectContentType(buffer)

	// Handle MIME type constraints
	if strings.Contains(expectedType, "/") {
		// Special handling for binary files
		if detectedType == "application/octet-stream" {
			// For binary files, trust the file extension more than the detected type
			if ext := strings.ToLower(filepath.Ext(filePath)); ext != "" {
				if _, ok := MimeTypes[ext]; ok {
					return nil // Accept if extension matches expected type
				}
			}
		}

		if !strings.HasPrefix(strings.ToLower(detectedType), strings.ToLower(expectedType)) {
			return fmt.Errorf("File is not of type %v (detected: %v)", expectedType, detectedType)
		}
		return nil
	}

	// Handle extension constraints
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == "" {
		return fmt.Errorf("File has no extension")
	}

	// Check if extension matches expected MIME type
	if expectedMime, ok := MimeTypes[strings.ToLower(expectedType)]; ok {
		if strings.HasPrefix(strings.ToLower(detectedType), strings.ToLower(expectedMime)) {
			return nil
		}
		return fmt.Errorf("File extension %v doesn't match content type (detected: %v, expected: %v)",
			ext, detectedType, expectedMime)
	}

	return fmt.Errorf("Unknown file type: %v", expectedType)
}
