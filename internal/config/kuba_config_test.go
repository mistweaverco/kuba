package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadKubaConfig(t *testing.T) {
	// Create a temporary test config file
	testConfig := `---
default:
  provider: gcp
  project: "test-project"
  mappings:
    - environment-variable: "TEST_VAR"
      secret-key: "test_secret"
    - environment-variable: "ANOTHER_VAR"
      secret-key: "another_secret"
      provider: aws
      project: "aws-project"

development:
  provider: gcp
  project: "dev-project"
  mappings:
    - environment-variable: "DEV_VAR"
      secret-key: "dev_secret"
`

	tmpFile, err := os.CreateTemp("", "kuba-test-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(testConfig); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}
	tmpFile.Close()

	// Test loading the config
	config, err := LoadKubaConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify default environment
	defaultEnv, err := config.GetEnvironment("default")
	if err != nil {
		t.Fatalf("Failed to get default environment: %v", err)
	}

	if defaultEnv.Provider != "gcp" {
		t.Errorf("Expected provider 'gcp', got '%s'", defaultEnv.Provider)
	}

	if defaultEnv.Project != "test-project" {
		t.Errorf("Expected project 'test-project', got '%s'", defaultEnv.Project)
	}

	if len(defaultEnv.Mappings) != 2 {
		t.Errorf("Expected 2 mappings, got %d", len(defaultEnv.Mappings))
	}

	// Verify first mapping
	if defaultEnv.Mappings[0].EnvironmentVariable != "TEST_VAR" {
		t.Errorf("Expected environment variable 'TEST_VAR', got '%s'", defaultEnv.Mappings[0].EnvironmentVariable)
	}

	if defaultEnv.Mappings[0].SecretKey != "test_secret" {
		t.Errorf("Expected secret key 'test_secret', got '%s'", defaultEnv.Mappings[0].SecretKey)
	}

	// Verify second mapping with override
	if defaultEnv.Mappings[1].Provider != "aws" {
		t.Errorf("Expected provider 'aws', got '%s'", defaultEnv.Mappings[1].Provider)
	}

	if defaultEnv.Mappings[1].Project != "aws-project" {
		t.Errorf("Expected project 'aws-project', got '%s'", defaultEnv.Mappings[1].Project)
	}

	// Verify development environment
	devEnv, err := config.GetEnvironment("development")
	if err != nil {
		t.Fatalf("Failed to get development environment: %v", err)
	}

	if devEnv.Provider != "gcp" {
		t.Errorf("Expected provider 'gcp', got '%s'", devEnv.Provider)
	}

	if devEnv.Project != "dev-project" {
		t.Errorf("Expected project 'dev-project', got '%s'", devEnv.Project)
	}
}

func TestGetEnvironmentDefault(t *testing.T) {
	config := &KubaConfig{
		Environments: map[string]Environment{
			"default": {
				Provider: "gcp",
				Project:  "test-project",
				Mappings: []Mapping{
					{
						EnvironmentVariable: "TEST_VAR",
						SecretKey:           "test_secret",
					},
				},
			},
		},
	}

	// Test getting default environment
	env, err := config.GetEnvironment("")
	if err != nil {
		t.Fatalf("Failed to get default environment: %v", err)
	}

	if env.Provider != "gcp" {
		t.Errorf("Expected provider 'gcp', got '%s'", env.Provider)
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *KubaConfig
		wantErr bool
	}{
		{
			name: "valid config with secret-key",
			config: &KubaConfig{
				Environments: map[string]Environment{
					"default": {
						Provider: "gcp",
						Project:  "test-project",
						Mappings: []Mapping{
							{
								EnvironmentVariable: "TEST_VAR",
								SecretKey:           "test_secret",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid config with value",
			config: &KubaConfig{
				Environments: map[string]Environment{
					"default": {
						Provider: "gcp",
						Project:  "test-project",
						Mappings: []Mapping{
							{
								EnvironmentVariable: "TEST_VAR",
								Value:               "test_value",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid config with mixed mappings",
			config: &KubaConfig{
				Environments: map[string]Environment{
					"default": {
						Provider: "gcp",
						Project:  "test-project",
						Mappings: []Mapping{
							{
								EnvironmentVariable: "SECRET_VAR",
								SecretKey:           "test_secret",
							},
							{
								EnvironmentVariable: "VALUE_VAR",
								Value:               "test_value",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing provider",
			config: &KubaConfig{
				Environments: map[string]Environment{
					"default": {
						Project: "test-project",
						Mappings: []Mapping{
							{
								EnvironmentVariable: "TEST_VAR",
								SecretKey:           "test_secret",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid provider",
			config: &KubaConfig{
				Environments: map[string]Environment{
					"default": {
						Provider: "invalid",
						Project:  "test-project",
						Mappings: []Mapping{
							{
								EnvironmentVariable: "TEST_VAR",
								SecretKey:           "test_secret",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing both secret-key and value",
			config: &KubaConfig{
				Environments: map[string]Environment{
					"default": {
						Provider: "gcp",
						Project:  "test-project",
						Mappings: []Mapping{
							{
								EnvironmentVariable: "TEST_VAR",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "both secret-key and value specified",
			config: &KubaConfig{
				Environments: map[string]Environment{
					"default": {
						Provider: "gcp",
						Project:  "test-project",
						Mappings: []Mapping{
							{
								EnvironmentVariable: "TEST_VAR",
								SecretKey:           "test_secret",
								Value:               "test_value",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "valid AWS config without project",
			config: &KubaConfig{
				Environments: map[string]Environment{
					"default": {
						Provider: "aws",
						Project:  "", // Empty project for AWS should be valid
						Mappings: []Mapping{
							{
								EnvironmentVariable: "AWS_SECRET",
								SecretKey:           "aws-secret-key",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid GCP config without project",
			config: &KubaConfig{
				Environments: map[string]Environment{
					"default": {
						Provider: "gcp",
						Project:  "", // Empty project for GCP should be invalid
						Mappings: []Mapping{
							{
								EnvironmentVariable: "GCP_SECRET",
								SecretKey:           "gcp-secret-key",
							},
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInterpolation(t *testing.T) {
	// Test basic environment variable interpolation
	t.Run("basic env var interpolation", func(t *testing.T) {
		// Set test environment variable
		os.Setenv("TEST_VAR", "test_value")
		defer os.Unsetenv("TEST_VAR")

		config := &KubaConfig{
			Environments: map[string]Environment{
				"default": {
					Provider: "gcp",
					Project:  "test-project",
					Mappings: []Mapping{
						{
							EnvironmentVariable: "DB_PASSWORD",
							Value:               "secret123",
						},
						{
							EnvironmentVariable: "DB_CONNECTION_STRING",
							Value:               "postgresql://user:${DB_PASSWORD}@host:5432/db",
						},
					},
				},
			},
		}

		err := processValueInterpolations(config)
		require.NoError(t, err)

		// Check that interpolation worked
		env := config.Environments["default"]
		require.Equal(t, "postgresql://user:secret123@host:5432/db", env.Mappings[1].Value)
	})

	t.Run("environment variable interpolation", func(t *testing.T) {
		// Set test environment variable
		os.Setenv("DB_HOST", "mydbhost")
		defer os.Unsetenv("DB_HOST")

		config := &KubaConfig{
			Environments: map[string]Environment{
				"default": {
					Provider: "gcp",
					Project:  "test-project",
					Mappings: []Mapping{
						{
							EnvironmentVariable: "DB_CONNECTION_STRING",
							Value:               "postgresql://user:pass@${DB_HOST}:5432/mydb",
						},
					},
				},
			},
		}

		err := processValueInterpolations(config)
		require.NoError(t, err)

		// Check that interpolation worked
		env := config.Environments["default"]
		require.Equal(t, "postgresql://user:pass@mydbhost:5432/mydb", env.Mappings[0].Value)
	})

	t.Run("mixed interpolation", func(t *testing.T) {
		// Set test environment variable
		os.Setenv("DB_PORT", "5432")
		defer os.Unsetenv("DB_PORT")

		config := &KubaConfig{
			Environments: map[string]Environment{
				"default": {
					Provider: "gcp",
					Project:  "test-project",
					Mappings: []Mapping{
						{
							EnvironmentVariable: "DB_PASSWORD",
							Value:               "secret123",
						},
						{
							EnvironmentVariable: "DB_HOST",
							Value:               "mydbhost",
						},
						{
							EnvironmentVariable: "DB_CONNECTION_STRING",
							Value:               "postgresql://user:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/mydb",
						},
					},
				},
			},
		}

		err := processValueInterpolations(config)
		require.NoError(t, err)

		// Check that interpolation worked
		env := config.Environments["default"]
		require.Equal(t, "postgresql://user:secret123@mydbhost:5432/mydb", env.Mappings[2].Value)
	})

	t.Run("no interpolation needed", func(t *testing.T) {
		config := &KubaConfig{
			Environments: map[string]Environment{
				"default": {
					Provider: "gcp",
					Project:  "test-project",
					Mappings: []Mapping{
						{
							EnvironmentVariable: "SIMPLE_VALUE",
							Value:               "no interpolation here",
						},
					},
				},
			},
		}

		err := processValueInterpolations(config)
		require.NoError(t, err)

		// Check that value remains unchanged
		env := config.Environments["default"]
		require.Equal(t, "no interpolation here", env.Mappings[0].Value)
	})

	t.Run("unresolved variable remains unchanged", func(t *testing.T) {
		config := &KubaConfig{
			Environments: map[string]Environment{
				"default": {
					Provider: "gcp",
					Project:  "test-project",
					Mappings: []Mapping{
						{
							EnvironmentVariable: "UNRESOLVED",
							Value:               "value with ${UNKNOWN_VAR}",
						},
					},
				},
			},
		}

		err := processValueInterpolations(config)
		require.NoError(t, err)

		// Check that unresolved variable remains unchanged
		env := config.Environments["default"]
		require.Equal(t, "value with ${UNKNOWN_VAR}", env.Mappings[0].Value)
	})

	t.Run("numeric values are converted to string", func(t *testing.T) {
		config := &KubaConfig{
			Environments: map[string]Environment{
				"default": {
					Provider: "gcp",
					Project:  "test-project",
					Mappings: []Mapping{
						{
							EnvironmentVariable: "PORT",
							Value:               8080,
						},
						{
							EnvironmentVariable: "URL",
							Value:               "http://localhost:${PORT}",
						},
					},
				},
			},
		}

		err := processValueInterpolations(config)
		require.NoError(t, err)

		// Check that numeric value was converted and interpolation worked
		env := config.Environments["default"]
		require.Equal(t, "8080", env.Mappings[0].Value)
		require.Equal(t, "http://localhost:8080", env.Mappings[1].Value)
	})

	t.Run("shell-style default value syntax", func(t *testing.T) {
		// Test with default value when environment variable doesn't exist
		config := &KubaConfig{
			Environments: map[string]Environment{
				"default": {
					Provider: "gcp",
					Project:  "test-project",
					Mappings: []Mapping{
						{
							EnvironmentVariable: "APP_ENV",
							Value:               "${NODE_ENV:-development}",
						},
						{
							EnvironmentVariable: "REDIS_URL",
							Value:               "redis://${REDIS_HOST:-localhost}:${REDIS_PORT:-6379}/0",
						},
					},
				},
			},
		}

		err := processValueInterpolations(config)
		require.NoError(t, err)

		// Check that default values were used
		env := config.Environments["default"]
		require.Equal(t, "development", env.Mappings[0].Value)
		require.Equal(t, "redis://localhost:6379/0", env.Mappings[1].Value)
	})

	t.Run("shell-style default value with existing env var", func(t *testing.T) {
		// Set test environment variable
		os.Setenv("NODE_ENV", "production")
		defer os.Unsetenv("NODE_ENV")

		config := &KubaConfig{
			Environments: map[string]Environment{
				"default": {
					Provider: "gcp",
					Project:  "test-project",
					Mappings: []Mapping{
						{
							EnvironmentVariable: "APP_ENV",
							Value:               "${NODE_ENV:-development}",
						},
					},
				},
			},
		}

		err := processValueInterpolations(config)
		require.NoError(t, err)

		// Check that environment variable value was used instead of default
		env := config.Environments["default"]
		require.Equal(t, "production", env.Mappings[0].Value)
	})

	t.Run("mixed default value syntax", func(t *testing.T) {
		// Set some environment variables
		os.Setenv("DB_HOST", "mydbhost")
		defer os.Unsetenv("DB_HOST")

		config := &KubaConfig{
			Environments: map[string]Environment{
				"default": {
					Provider: "gcp",
					Project:  "test-project",
					Mappings: []Mapping{
						{
							EnvironmentVariable: "DB_PASSWORD",
							Value:               "secret123",
						},
						{
							EnvironmentVariable: "DB_CONNECTION_STRING",
							Value:               "postgresql://user:${DB_PASSWORD}@${DB_HOST:-localhost}:${DB_PORT:-5432}/mydb",
						},
					},
				},
			},
		}

		err := processValueInterpolations(config)
		require.NoError(t, err)

		// Check that interpolation worked with mixed syntax
		env := config.Environments["default"]
		require.Equal(t, "postgresql://user:secret123@mydbhost:5432/mydb", env.Mappings[1].Value)
	})
}

func TestLoadKubaConfigWithInterpolation(t *testing.T) {
	// Test loading a config file with interpolation
	t.Run("load config with interpolation", func(t *testing.T) {
		// Set test environment variable
		os.Setenv("DB_HOST", "mydbhost")
		defer os.Unsetenv("DB_HOST")

		// Create a temporary config file
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "kuba.yaml")

		configContent := `default:
  provider: gcp
  project: test-project
  mappings:
    - environment-variable: "DB_PASSWORD"
      value: "secret123"
    - environment-variable: "DB_CONNECTION_STRING"
      value: "postgresql://user:${DB_PASSWORD}@${DB_HOST}:5432/mydb"
`

		err := os.WriteFile(configPath, []byte(configContent), 0644)
		require.NoError(t, err)

		// Load the config
		config, err := LoadKubaConfig(configPath)
		require.NoError(t, err)

		// Check that interpolation worked
		env := config.Environments["default"]
		require.Equal(t, "postgresql://user:secret123@mydbhost:5432/mydb", env.Mappings[1].Value)
	})
}
