package config

import (
	"os"
	"testing"
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
