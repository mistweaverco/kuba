package fileutils

import (
	"fmt"
	"os"

	"github.com/spf13/afero"
)

// FileSystem interface for filesystem operations
type FileSystem interface {
	Create(name string) (afero.File, error)
	MkdirAll(path string, perm os.FileMode) error
	OpenFile(name string, flag int, perm os.FileMode) (afero.File, error)
	Stat(name string) (os.FileInfo, error)
	UserConfigDir() (string, error)
	TempDir() string
	Getenv(key string) string
	WriteString(file afero.File, s string) (int, error)
	Close(file afero.File) error
}

// defaultFileSystem implements FileSystem using Afero
type defaultFileSystem struct {
	fs afero.Fs
}

func (d *defaultFileSystem) Create(name string) (afero.File, error) {
	return d.fs.Create(name)
}

func (d *defaultFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return d.fs.MkdirAll(path, perm)
}

func (d *defaultFileSystem) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	return d.fs.OpenFile(name, flag, perm)
}

func (d *defaultFileSystem) Stat(name string) (os.FileInfo, error) {
	return d.fs.Stat(name)
}

func (d *defaultFileSystem) UserConfigDir() (string, error) {
	return os.UserConfigDir()
}

func (d *defaultFileSystem) TempDir() string {
	return os.TempDir()
}

func (d *defaultFileSystem) Getenv(key string) string {
	return os.Getenv(key)
}

func (d *defaultFileSystem) WriteString(file afero.File, s string) (int, error) {
	return file.WriteString(s)
}

func (d *defaultFileSystem) Close(file afero.File) error {
	return file.Close()
}

// Global variables for dependency injection
var (
	fileSystem FileSystem = &defaultFileSystem{fs: afero.NewOsFs()}
)

// SetFileSystem sets the file system implementation
func SetFileSystem(fs FileSystem) {
	fileSystem = fs
}

// ResetDependencies resets all dependencies to their default implementations
func ResetDependencies() {
	fileSystem = &defaultFileSystem{fs: afero.NewOsFs()}
}

func FileExists(path string) bool {
	if path == "" {
		return false
	}
	_, err :=
		fileSystem.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

// GenerateDefaultKubaConfig creates a default kuba.yaml file
// in the current directory if it doesn't exist.
func GenerateDefaultKubaConfig() bool {
	fp := "kuba.yaml"

	if FileExists(fp) {
		return false // File already exists, no need to create it
	}

	contents := `# yaml-language-server: $schema=https://raw.githubusercontent.com/mistweaverco/kuba/refs/heads/main/kuba.schema.json
---
# Top-level sections for different environments.
default:
  provider: gcp
  project: "my-gcp-project-default"

  # Mapping of cloud projects to environment variables and secret keys.
  mappings:
    - environment-variable: "GCP_PROJECT_ID"
      secret-key: "gcp_project_secret"
    - environment-variable: "AWS_PROJECT_ID"
      secret-key: "aws_project_secret"
      provider: aws
      project: "my-aws-project-default"
    - environment-variable: "AZURE_PROJECT_ID"
      secret-key: "azure_project_secret"
      provider: azure
      project: "my-azure-project-default"
`

	file, err := fileSystem.Create(fp)
	if err != nil {
		fmt.Println("Error creating kuba.yaml:", err)
		return false
	}
	defer func() {
		if closeErr := fileSystem.Close(file); closeErr != nil {
			fmt.Printf("Warning: failed to close kuba.yaml file: %v\n", closeErr)
		}
	}()

	_, err = fileSystem.WriteString(file, contents)
	if err != nil {
		fmt.Println("Error writing to kuba.yaml:", err)
		return false
	}

	return true
}

// GetAppDataPath returns the path to the app data directory
// If the KUBA_HOME environment variable is set, it will use that path
// otherwise it will use the user's config directory
// e.g. /home/user/.config/kuba
func GetAppDataPath() string {
	if kubaHome := fileSystem.Getenv("KUBA_HOME"); kubaHome != "" {
		return EnsureDirExists(kubaHome)
	}
	userConfigDir, err := fileSystem.UserConfigDir()
	if err != nil {
		panic(err)
	}
	return EnsureDirExists(userConfigDir + string(os.PathSeparator) + "kuba")
}

// GetTempPath returns the path to the temp directory
// e.g. /tmp
func GetTempPath() string {
	return fileSystem.TempDir()
}

func EnsureDirExists(path string) string {
	if _, err := fileSystem.Stat(path); os.IsNotExist(err) {
		if err := fileSystem.MkdirAll(path, 0755); err != nil {
			// Log the error but don't fail the function
			fmt.Printf("Warning: failed to create directory %s: %v\n", path, err)
		}
	}
	return path
}
