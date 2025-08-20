package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// KubaConfig represents the structure of a kuba.yaml file
type KubaConfig struct {
	Environments map[string]Environment `yaml:",inline"`
}

// Environment represents a single environment configuration
type Environment struct {
	Provider string    `yaml:"provider"`
	Project  string    `yaml:"project"`
	Mappings []Mapping `yaml:"mappings"`
}

// Mapping represents a mapping between environment variable and secret key or value
type Mapping struct {
	EnvironmentVariable string      `yaml:"environment-variable"`
	SecretKey           string      `yaml:"secret-key,omitempty"`
	Value               interface{} `yaml:"value,omitempty"`
	Provider            string      `yaml:"provider,omitempty"`
	Project             string      `yaml:"project,omitempty"`
}

// LoadKubaConfig loads the kuba.yaml configuration file
func LoadKubaConfig(configPath string) (*KubaConfig, error) {
	if configPath == "" {
		configPath = "kuba.yaml"
	}

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("configuration file not found: %s", configPath)
	}

	// Read file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %w", err)
	}

	// Parse YAML
	var config KubaConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse configuration file: %w", err)
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// GetEnvironment returns the configuration for a specific environment
func (c *KubaConfig) GetEnvironment(envName string) (*Environment, error) {
	if envName == "" {
		envName = "default"
	}

	env, exists := c.Environments[envName]
	if !exists {
		return nil, fmt.Errorf("environment '%s' not found in configuration", envName)
	}

	return &env, nil
}

// validateConfig validates the configuration structure
func validateConfig(config *KubaConfig) error {
	if len(config.Environments) == 0 {
		return fmt.Errorf("no environments defined in configuration")
	}

	for envName, env := range config.Environments {
		if env.Provider == "" {
			return fmt.Errorf("environment '%s': provider is required", envName)
		}

		// Project is required for all providers except AWS and Azure
		if env.Project == "" && env.Provider != "aws" && env.Provider != "azure" {
			return fmt.Errorf("environment '%s': project is required for provider '%s'", envName, env.Provider)
		}

		if len(env.Mappings) == 0 {
			return fmt.Errorf("environment '%s': at least one mapping is required", envName)
		}

		for i, mapping := range env.Mappings {
			if mapping.EnvironmentVariable == "" {
				return fmt.Errorf("environment '%s': mapping %d: environment-variable is required", envName, i+1)
			}

			// Either secret-key or value must be provided, but not both
			if mapping.SecretKey == "" && mapping.Value == nil {
				return fmt.Errorf("environment '%s': mapping %d: either secret-key or value is required", envName, i+1)
			}

			if mapping.SecretKey != "" && mapping.Value != nil {
				return fmt.Errorf("environment '%s': mapping %d: cannot specify both secret-key and value", envName, i+1)
			}

			// Validate provider if specified and secret-key is used
			if mapping.Provider != "" && mapping.SecretKey != "" {
				if !isValidProvider(mapping.Provider) {
					return fmt.Errorf("environment '%s': mapping %d: invalid provider '%s'", envName, i+1, mapping.Provider)
				}
			}
		}

		// Validate main provider
		if !isValidProvider(env.Provider) {
			return fmt.Errorf("environment '%s': invalid provider '%s'", envName, env.Provider)
		}
	}

	return nil
}

// isValidProvider checks if the provider is supported
func isValidProvider(provider string) bool {
	validProviders := []string{"gcp", "aws", "azure"}
	for _, p := range validProviders {
		if p == provider {
			return true
		}
	}
	return false
}

// FindConfigFile searches for a kuba.yaml file in the current directory and parent directories
func FindConfigFile() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	// Search up the directory tree for kuba.yaml
	for {
		configPath := filepath.Join(currentDir, "kuba.yaml")
		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			break // Reached root directory
		}
		currentDir = parent
	}

	return "", fmt.Errorf("kuba.yaml not found in current directory or any parent directory")
}
