package bundle

import (
	"io/fs"
	"path/filepath"
)

// BundleManifest represents all files included in a command's bundle
type BundleManifest struct {
	// Command path (e.g., "mycli:subcommand")
	CommandPath string
	// Root directory for this bundle
	RootDir string
	// Map of relative paths to file entries
	Files map[string]*BundleFile
}

// BundleFile represents a single file in the bundle
type BundleFile struct {
	Path       string      // Relative path within bundle
	IsDir      bool        // Is this a directory?
	Mode       fs.FileMode // File permissions
	EmbeddedFS fs.FS       // For directories (using embed.FS)
	Content    []byte      // For single files (using embed.StaticFile)
}

func NewBundleManifest(commandPath string) *BundleManifest {
	return &BundleManifest{
		CommandPath: commandPath,
		Files:       make(map[string]*BundleFile),
	}
}

// AddFile adds a file or directory to the manifest
func (m *BundleManifest) AddFile(path string, isDir bool, fs fs.FS, content []byte) error {
	relPath, err := filepath.Rel(m.RootDir, path)
	if err != nil {
		return err
	}

	m.Files[relPath] = &BundleFile{
		Path:       relPath,
		IsDir:      isDir,
		EmbeddedFS: fs,
		Content:    content,
	}
	return nil
}
