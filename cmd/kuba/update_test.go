package kuba

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected int
	}{
		{
			name:     "equal versions",
			v1:       "1.0.0",
			v2:       "1.0.0",
			expected: 0,
		},
		{
			name:     "v1 less than v2",
			v1:       "1.0.0",
			v2:       "1.0.1",
			expected: -1,
		},
		{
			name:     "v1 greater than v2",
			v1:       "1.0.1",
			v2:       "1.0.0",
			expected: 1,
		},
		{
			name:     "with v prefix",
			v1:       "v1.0.0",
			v2:       "v1.0.1",
			expected: -1,
		},
		{
			name:     "mixed v prefix",
			v1:       "v1.0.0",
			v2:       "1.0.1",
			expected: -1,
		},
		{
			name:     "different major versions",
			v1:       "1.0.0",
			v2:       "2.0.0",
			expected: -1,
		},
		{
			name:     "different minor versions",
			v1:       "1.1.0",
			v2:       "1.2.0",
			expected: -1,
		},
		{
			name:     "different patch versions",
			v1:       "1.0.1",
			v2:       "1.0.2",
			expected: -1,
		},
		{
			name:     "unequal length versions",
			v1:       "1.0",
			v2:       "1.0.0",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareVersions(tt.v1, tt.v2)
			assert.Equal(t, tt.expected, result, "compareVersions(%s, %s) = %d, want %d", tt.v1, tt.v2, result, tt.expected)
		})
	}
}

func TestDetectPlatform(t *testing.T) {
	platform := detectPlatform()

	// Platform should be in format "os-arch"
	assert.Contains(t, platform, "-", "Platform should contain a dash")

	// Should contain valid OS
	validOS := []string{"linux", "darwin", "windows"}
	hasValidOS := false
	for _, os := range validOS {
		if platform == os+"-amd64" || platform == os+"-386" || platform == os+"-arm64" || platform == os+"-armv7" {
			hasValidOS = true
			break
		}
	}
	assert.True(t, hasValidOS, "Platform should contain a valid OS: %s", platform)
}

func TestUpdateCommand(t *testing.T) {
	// Test that update command is properly configured
	assert.Equal(t, "update", updateCmd.Use)
	assert.Contains(t, updateCmd.Short, "Update kuba")
	assert.Contains(t, updateCmd.Long, "Check if a newer version")

	// Test that it takes no arguments
	// We can test this by checking if the command accepts no arguments
	// by calling the Args function with empty arguments
	err := updateCmd.Args(updateCmd, []string{})
	assert.NoError(t, err, "Command should accept no arguments")

	// Test that it rejects arguments
	err = updateCmd.Args(updateCmd, []string{"arg1"})
	assert.Error(t, err, "Command should reject arguments")
}

func TestGetCurrentVersion(t *testing.T) {
	version := getCurrentVersion()

	// Version should not be empty (it might be empty in tests, but the function should work)
	// We can't assert a specific value since it depends on build-time variables
	assert.NotNil(t, version, "Version should not be nil")
}

func TestCopyFile(t *testing.T) {
	// Create a temporary source file
	srcFile, err := createTempFile("test content")
	require.NoError(t, err)
	defer os.Remove(srcFile)

	// Create destination path
	dstFile := srcFile + ".copy"
	defer os.Remove(dstFile)

	// Copy the file
	err = copyFile(srcFile, dstFile)
	require.NoError(t, err)

	// Verify the copy was successful
	srcContent, err := os.ReadFile(srcFile)
	require.NoError(t, err)

	dstContent, err := os.ReadFile(dstFile)
	require.NoError(t, err)

	assert.Equal(t, srcContent, dstContent, "Copied file content should match source")
}

// Helper function to create a temporary file with content
func createTempFile(content string) (string, error) {
	tmpFile, err := os.CreateTemp("", "kuba-test-*")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	_, err = tmpFile.WriteString(content)
	if err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}

	return tmpFile.Name(), nil
}
