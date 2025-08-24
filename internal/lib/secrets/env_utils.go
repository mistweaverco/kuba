package secrets

import (
	"regexp"
	"strings"
)

// sanitizeEnvVarName sanitizes a string to be a valid POSIX environment variable name
// POSIX rules: must begin with a letter or underscore, and contain only letters, numbers, or underscores
func sanitizeEnvVarName(name string) string {
	if name == "" {
		return "_"
	}

	// Replace any non-alphanumeric characters with underscores
	re := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	sanitized := re.ReplaceAllString(name, "_")

	// Ensure it starts with a letter or underscore
	if len(sanitized) > 0 {
		firstChar := sanitized[0]
		if (firstChar < 'a' || firstChar > 'z') && (firstChar < 'A' || firstChar > 'Z') && firstChar != '_' {
			sanitized = "_" + sanitized
		}
	}

	// If the result is empty after sanitization, return underscore
	if sanitized == "" {
		return "_"
	}

	return sanitized
}

// extractSecretNameFromPath extracts the secret name from a full secret path
// This is useful for providers where the path contains additional metadata
func extractSecretNameFromPath(path string) string {
	// Split by common separators and take the last part
	parts := strings.FieldsFunc(path, func(r rune) bool {
		return r == '/' || r == '\\' || r == ':'
	})

	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return path
}
