package secrets

import (
	"testing"
)

func TestSanitizeEnvVarName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal-name", "normal_name"},
		{"name_with_underscores", "name_with_underscores"},
		{"123-start-with-number", "_123_start_with_number"},
		{"start-with-dash", "start_with_dash"},
		{"has.dots", "has_dots"},
		{"has spaces", "has_spaces"},
		{"has@special#chars", "has_special_chars"},
		{"", "_"},
		{"a", "a"},
		{"A", "A"},
		{"_", "_"},
		{"1", "_1"},
		{"-", "_"},
		{"my-app-config", "my_app_config"},
		{"database.password", "database_password"},
		{"API_KEY_123", "API_KEY_123"},
	}

	for _, test := range tests {
		result := sanitizeEnvVarName(test.input)
		if result != test.expected {
			t.Errorf("sanitizeEnvVarName(%q) = %q, want %q", test.input, result, test.expected)
		}
	}
}

func TestExtractSecretNameFromPath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple-name", "simple-name"},
		{"path/to/secret", "secret"},
		{"path\\to\\secret", "secret"},
		{"path:to:secret", "secret"},
		{"projects/myproject/secrets/mysecret", "mysecret"},
		{"", ""},
		{"single", "single"},
	}

	for _, test := range tests {
		result := extractSecretNameFromPath(test.input)
		if result != test.expected {
			t.Errorf("extractSecretNameFromPath(%q) = %q, want %q", test.input, result, test.expected)
		}
	}
}
