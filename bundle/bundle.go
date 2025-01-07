package bundle

import (
	"bytes"
	"embed"
	"runtime"

	"github.com/migsc/cmdeagle/file"
	"github.com/migsc/cmdeagle/types"

	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/spf13/afero"
)

//go:embed *
var PackageFS embed.FS

var StagingDirName = ".tmp"
var MainTemplateFileName = "main.template.go"
var MainTemplateReplacements = map[string][]byte{
	// The package switcharoo
	"package bundle": []byte("package main"),

	// Rename the template's intended main function
	"main_template": []byte("main"),
}

var packageSrcDirPath string
var bundleStagingDirPath string

func GetPackageSrcDirPath() string {
	log.Debug("Getting package src dir path:")
	if packageSrcDirPath != "" {
		log.Debug("Package src dir path already set:", "path", packageSrcDirPath)
		return packageSrcDirPath
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("failed to get current file path")
	}

	log.Debug("Package src dir path not set, setting it now:", "path", filepath.Dir(filename))

	packageSrcDirPath = filepath.Dir(filename)

	return packageSrcDirPath
}

func SetupStagingDir() (string, error) {
	// TODO possible to move into file package? Or is it reliant on executing from the bundle package?
	tempDirPath, err := afero.TempDir(afero.NewOsFs(), "", "cmdeagle")
	if err != nil {
		return "", err
	}

	log.Debug("Using bundle staging directory path for bundle of: %s\n", tempDirPath)

	bundleStagingDirPath = tempDirPath

	return tempDirPath, nil
}

func GetMainTemplateContent() ([]byte, error) {
	mainTemplateContent, err := os.ReadFile(filepath.Join(GetPackageSrcDirPath(), MainTemplateFileName))
	if err != nil {
		return nil, fmt.Errorf("error reading template: %w", err)
	}
	return mainTemplateContent, nil
}

func InterpolateMainContent(content []byte) []byte {
	for old, new := range MainTemplateReplacements {
		content = bytes.ReplaceAll(content, []byte(old), new)
	}

	// TODO: We need to embed a filesystem of the current working directory (where main.go will be)?

	return content
}

func SetupMainFile(path string) (string, error) {
	mainTemplateContent, err := GetMainTemplateContent()
	if err != nil {
		return "", fmt.Errorf("error reading template: %w", err)
	}

	for old, new := range MainTemplateReplacements {
		mainTemplateContent = bytes.ReplaceAll(mainTemplateContent, []byte(old), new)
	}

	// Write processed template to main.go
	mainFilePath := filepath.Join(path, "main.go")
	if err := os.WriteFile(mainFilePath, mainTemplateContent, 0644); err != nil {
		return "", fmt.Errorf("error writing main.go: %w", err)
	}

	return mainFilePath, nil
}

type Bundle struct {
	DirPath string
	DirName string
	Files   map[string]*BundleFile
}

// func CreateBundle(config *config.CmdeagleConfig) (*Bundle, error) {
// 	var bundle = Bundle{
// 		DirPath: config.DataDir,
// 		DirName: config.Name,
// 		Files:   make(map[string]*BundleFile),
// 	}

// 	return &bundle, nil
// }

// type BundleFile struct {
// 	path     string
// 	info     os.FileInfo
// 	children []*BundleFile
// 	parents  []*BundleFile
// }

// func GetDestDir(config *schema.CmdeagleConfig) (string, error) {
// 	if config.DataDir == "" {
// 		return "", fmt.Errorf("Directory path was not provided for data/bundle directory.")
// 	}

// 	return filepath.Join(config.DataDir, config.Name), nil
// }

func CopyIncludedFiles(config *types.CmdeagleConfig, command *types.CommandDefinition, namespace []string, targetDirPath string) error {
	if len(command.Includes) == 0 {
		return nil
	}

	log.Info("processing includes",
		"command", command.Name,
		"includes", command.Includes,
	)

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	for _, ns := range namespace {
		targetDirPath = filepath.Join(targetDirPath, ns)
	}

	for _, includePath := range command.Includes {
		if err := copyIncludedFile(filepath.Join(currentDir, includePath), targetDirPath); err != nil {
			return err

		}
	}

	return nil
}

func copyIncludedFile(includedFilePath string, targetDir string) error {
	log.Info("including bundle",
		"from", includedFilePath,
		"to", targetDir,
	)

	expandedPath, err := file.ExpandPath(includedFilePath)
	if err != nil {
		return fmt.Errorf("could not expand the path: %s\n%v", includedFilePath, err)
	}

	// We can't do this here because it would result in files being removed iteratively.
	// if err := file.SetupEmptyDir(targetDir); err != nil {
	// 	return err
	// }

	cmd := exec.Command("cp", "-r", expandedPath, targetDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// type FileManifest struct {
// 	files []os.FileInfo
// }

// func marshalFileManifest(manifest *FileManifest) ([]byte, error) {
// 	if manifest == nil {
// 		return nil, fmt.Errorf("manifest cannot be nil")
// 	}

// 	fileInfos := make([]map[string]interface{}, len(manifest.files))
// 	for i, f := range manifest.files {
// 		fileInfos[i] = map[string]interface{}{
// 			"name":    f.Name(),
// 			"size":    f.Size(),
// 			"mode":    f.Mode().String(),
// 			"modTime": f.ModTime(),
// 			"isDir":   f.IsDir(),
// 		}
// 	}

// 	return yaml.Marshal(map[string]interface{}{
// 		"files": fileInfos,
// 	})
// }
