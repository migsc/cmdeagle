package params

import (
	"testing"

	"github.com/migsc/cmdeagle/file"
	"github.com/migsc/cmdeagle/types"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestValidateFileConstraints(t *testing.T) {
	fs := afero.NewMemMapFs()

	// Create test directory and files
	err := fs.MkdirAll("/testdir", 0755)
	assert.NoError(t, err)

	// Create a text file with proper content
	textContent := "This is a text file.\nIt has multiple lines.\nAll of them are plain text.\n"
	err = afero.WriteFile(fs, "/testdir/test.txt", []byte(textContent), 0644)
	assert.NoError(t, err)

	// Create an image file with proper JPEG header
	jpegHeader := []byte{
		0xFF, 0xD8, 0xFF, 0xE0, // JPEG SOI and APP0 marker
		0x00, 0x10, // APP0 length
		0x4A, 0x46, 0x49, 0x46, 0x00, // "JFIF" marker
		0x01, 0x01, // version
		0x00,       // units
		0x00, 0x01, // X density
		0x00, 0x01, // Y density
		0x00, 0x00, // thumbnail
	}
	err = afero.WriteFile(fs, "/testdir/test.jpg", append(jpegHeader, []byte("test image data")...), 0644)
	assert.NoError(t, err)

	// Set permissions
	err = fs.Chmod("/testdir/test.txt", 0644)
	assert.NoError(t, err)

	t.Run("validates directory exists", func(t *testing.T) {
		constraint := &types.ParamConstraints{
			DirExists: "true",
		}

		// Test valid directory
		err := validateFileConstraints(fs, constraint, "/testdir")
		assert.NoError(t, err)

		// Test non-existent directory
		err = validateFileConstraints(fs, constraint, "/nonexistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not a directory")

		// Test file as directory
		err = validateFileConstraints(fs, constraint, "/testdir/test.txt")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not a directory")
	})

	t.Run("validates file exists", func(t *testing.T) {
		constraint := &types.ParamConstraints{
			FileExists: "true",
		}

		// Test valid file
		err := validateFileConstraints(fs, constraint, "/testdir/test.txt")
		assert.NoError(t, err)

		// Test non-existent file
		err = validateFileConstraints(fs, constraint, "/testdir/nonexistent.txt")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist")

		// Test directory as file
		err = validateFileConstraints(fs, constraint, "/testdir")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "is a directory")
	})

	t.Run("validates file permissions", func(t *testing.T) {
		constraint := &types.ParamConstraints{
			FileExists:     "true",
			HasPermissions: "0644",
		}

		// Test correct permissions
		err := validateFileConstraints(fs, constraint, "/testdir/test.txt")
		assert.NoError(t, err)

		// Create file with different permissions
		err = afero.WriteFile(fs, "/testdir/test_perm.txt", []byte("test"), 0600)
		assert.NoError(t, err)

		// Test incorrect permissions
		constraint.HasPermissions = "0755"
		err = validateFileConstraints(fs, constraint, "/testdir/test_perm.txt")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "incorrect permissions")
	})

	/* Commenting out MIME type tests for now
	t.Run("validates file type", func(t *testing.T) {
		// Test text file
		constraint := &types.ParamConstraints{
			FileExists: "true",
			IsFileType: "text/plain",
		}
		err := validateFileConstraints(fs, constraint, "/testdir/test.txt")
		assert.NoError(t, err)

		// Test image file
		constraint = &types.ParamConstraints{
			FileExists: "true",
			IsFileType: "image/jpeg",
		}
		err = validateFileConstraints(fs, constraint, "/testdir/test.jpg")
		assert.NoError(t, err)

		// Test wrong file type (text file with image type)
		constraint = &types.ParamConstraints{
			FileExists: "true",
			IsFileType: "image/jpeg",
		}
		err = validateFileConstraints(fs, constraint, "/testdir/test.txt")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not of type")

		// Test by extension
		constraint = &types.ParamConstraints{
			FileExists: "true",
			IsFileType: ".txt",
		}
		err = validateFileConstraints(fs, constraint, "/testdir/test.txt")
		assert.NoError(t, err)

		// Test wrong extension
		constraint = &types.ParamConstraints{
			FileExists: "true",
			IsFileType: ".jpg",
		}
		err = validateFileConstraints(fs, constraint, "/testdir/test.txt")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "doesn't match content type")

		// Test non-existent file
		constraint = &types.ParamConstraints{
			FileExists: "true",
			IsFileType: "text/plain",
		}
		err = validateFileConstraints(fs, constraint, "/testdir/nonexistent.txt")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist")
	})
	*/
}

func TestHasFileConstraints(t *testing.T) {
	t.Run("detects file constraints", func(t *testing.T) {
		// Test with file constraints
		constraint := &types.ParamConstraints{
			FileExists: "true",
		}
		assert.True(t, HasFileConstraints(constraint))

		constraint = &types.ParamConstraints{
			DirExists: "true",
		}
		assert.True(t, HasFileConstraints(constraint))

		constraint = &types.ParamConstraints{
			IsFileType: "text/plain",
		}
		assert.True(t, HasFileConstraints(constraint))

		// Test without file constraints
		minVal := 1.0
		maxVal := 10.0
		constraint = &types.ParamConstraints{
			MinValue: &minVal,
			MaxValue: &maxVal,
		}
		assert.False(t, HasFileConstraints(constraint))

		// Test nil constraints
		assert.False(t, HasFileConstraints(nil))
	})
}

func TestValidateFileType(t *testing.T) {
	t.Run("validates_by_MIME_type", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		textContent := "This is a text file.\nIt has multiple lines.\nAll of them are plain text.\n"
		afero.WriteFile(fs, "test.txt", []byte(textContent), 0644)

		// Test text file validation
		err := file.ValidateFileType(fs, "test.txt", "text/plain")
		assert.NoError(t, err)

		// Test wrong MIME type
		err = file.ValidateFileType(fs, "test.txt", "image/jpeg")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not of type image/jpeg")
	})

	t.Run("handles_errors", func(t *testing.T) {
		fs := afero.NewMemMapFs()

		// Test non-existent file
		err := file.ValidateFileType(fs, "nonexistent.txt", "text/plain")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file does not exist")
	})
}

// Helper function to convert int to pointer
func toPtr(v int) *int {
	return &v
}
