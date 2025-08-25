package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mistweaverco/kuba/internal/lib/log"
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
	SecretPath          string      `yaml:"secret-path,omitempty"`
	Value               interface{} `yaml:"value,omitempty"`
	Provider            string      `yaml:"provider,omitempty"`
	Project             string      `yaml:"project,omitempty"`
}

// interpolateEnvVars replaces ${VAR_NAME} patterns with actual environment variable values
// It also supports previously resolved variables from the same configuration
// Supports both ${VAR_NAME} and ${VAR_NAME:-default} syntax
func interpolateEnvVars(value string, resolvedVars map[string]string) string {
	// Regex to match ${VAR_NAME} and ${VAR_NAME:-default} patterns
	re := regexp.MustCompile(`\$\{([^}]+)\}`)

	return re.ReplaceAllStringFunc(value, func(match string) string {
		// Extract the variable name and optional default from ${VAR_NAME} or ${VAR_NAME:-default}
		content := match[2 : len(match)-1]

		// Check if there's a default value specified
		if strings.Contains(content, ":-") {
			parts := strings.SplitN(content, ":-", 2)
			varName := parts[0]
			defaultValue := parts[1]

			// First check if we have this variable from previously resolved mappings
			if resolvedValue, exists := resolvedVars[varName]; exists {
				return resolvedValue
			}

			// Then check if it's an environment variable
			if envValue := os.Getenv(varName); envValue != "" {
				return envValue
			}

			// If not found, return the default value
			return defaultValue
		}

		// No default value specified, use original logic
		varName := content

		// First check if we have this variable from previously resolved mappings
		if resolvedValue, exists := resolvedVars[varName]; exists {
			return resolvedValue
		}

		// Then check if it's an environment variable
		if envValue := os.Getenv(varName); envValue != "" {
			return envValue
		}

		// If not found, return the original pattern (could be useful for debugging)
		return match
	})
}

// processValueInterpolations processes all value fields in mappings to resolve environment variable interpolations
func processValueInterpolations(config *KubaConfig) error {
	// Process environments in order to handle dependencies correctly
	// We'll process each environment multiple times until no more interpolations are possible
	// or until we detect a circular dependency

	for envName, env := range config.Environments {
		// Track resolved variables for this environment
		resolvedVars := make(map[string]string)

		// Process mappings multiple times to handle dependencies
		maxIterations := len(env.Mappings) * 2 // Allow for some dependency depth
		for iteration := 0; iteration < maxIterations; iteration++ {
			changed := false

			for i, mapping := range env.Mappings {
				if mapping.Value != nil {
					// Convert value to string for processing
					var strValue string
					switch v := mapping.Value.(type) {
					case string:
						strValue = v
					case int, int32, int64:
						strValue = fmt.Sprintf("%d", v)
					case float32, float64:
						strValue = fmt.Sprintf("%g", v)
					default:
						strValue = fmt.Sprintf("%v", v)
					}

					// Check if this value contains interpolation patterns
					if strings.Contains(strValue, "${") {
						// Interpolate the value
						interpolatedValue := interpolateEnvVars(strValue, resolvedVars)

						// If the value changed, update it
						if interpolatedValue != strValue {
							// Update the mapping value
							env.Mappings[i].Value = interpolatedValue
							// Update our resolved vars map
							resolvedVars[mapping.EnvironmentVariable] = interpolatedValue
							changed = true
						}
					} else {
						// No interpolation needed, but convert numeric values to strings for consistency
						if mapping.Value != strValue {
							env.Mappings[i].Value = strValue
							changed = true
						}
						// Store the value in resolved vars
						resolvedVars[mapping.EnvironmentVariable] = strValue
					}
				}
			}

			// If no changes were made in this iteration, we're done
			if !changed {
				break
			}
		}

		// Update the environment in the config
		config.Environments[envName] = env
	}

	return nil
}

// LoadKubaConfig loads the kuba.yaml configuration file
func LoadKubaConfig(configPath string) (*KubaConfig, error) {
	logger := log.NewLogger()

	if configPath == "" {
		configPath = "kuba.yaml"
	}

	logger.Debug("Loading configuration file", "path", configPath)

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		logger.Debug("Configuration file not found", "path", configPath)
		return nil, fmt.Errorf("configuration file not found: %s", configPath)
	}

	// Read file
	logger.Debug("Reading configuration file")
	data, err := os.ReadFile(configPath)
	if err != nil {
		logger.Debug("Failed to read configuration file", "path", configPath, "error", err)
		return nil, fmt.Errorf("failed to read configuration file: %w", err)
	}

	logger.Debug("Configuration file read successfully", "size_bytes", len(data))

	// Parse YAML
	logger.Debug("Parsing YAML configuration")
	var config KubaConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		logger.Debug("Failed to parse YAML configuration", "error", err)
		return nil, fmt.Errorf("failed to parse configuration file: %w", err)
	}

	logger.Debug("YAML parsed successfully", "environments_count", len(config.Environments))

	// Process environment variable interpolations
	logger.Debug("Processing environment variable interpolations")
	if err := processValueInterpolations(&config); err != nil {
		logger.Debug("Failed to process environment variable interpolations", "error", err)
		return nil, fmt.Errorf("failed to process environment variable interpolations: %w", err)
	}

	logger.Debug("Environment variable interpolations processed successfully")

	// Validate configuration
	logger.Debug("Validating configuration")
	if err := validateConfig(&config); err != nil {
		logger.Debug("Configuration validation failed", "error", err)
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	logger.Debug("Configuration validation passed")
	return &config, nil
}

// GetEnvironment returns the configuration for a specific environment
func (c *KubaConfig) GetEnvironment(envName string) (*Environment, error) {
	logger := log.NewLogger()

	if envName == "" {
		envName = "default"
		logger.Debug("No environment specified, using default")
	}

	logger.Debug("Getting environment configuration", "requested_env", envName, "available_environments", len(c.Environments))

	env, exists := c.Environments[envName]
	if !exists {
		logger.Debug("Environment not found in configuration", "requested_env", envName, "available_environments", getEnvironmentNames(c.Environments))
		return nil, fmt.Errorf("environment '%s' not found in configuration", envName)
	}

	logger.Debug("Environment configuration retrieved", "environment", envName, "provider", env.Provider, "project", env.Project, "mappings_count", len(env.Mappings))
	return &env, nil
}

// getEnvironmentNames returns a slice of available environment names
func getEnvironmentNames(environments map[string]Environment) []string {
	names := make([]string, 0, len(environments))
	for name := range environments {
		names = append(names, name)
	}
	return names
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

		// Project is required for all providers except AWS, Azure, and OpenBao
		if env.Project == "" && env.Provider != "aws" && env.Provider != "azure" && env.Provider != "openbao" {
			return fmt.Errorf("environment '%s': project is required for provider '%s'", envName, env.Provider)
		}

		if len(env.Mappings) == 0 {
			return fmt.Errorf("environment '%s': at least one mapping is required", envName)
		}

		for i, mapping := range env.Mappings {
			if mapping.EnvironmentVariable == "" {
				return fmt.Errorf("environment '%s': mapping %d: environment-variable is required", envName, i+1)
			}

			// Either secret-key, secret-path, or value must be provided, but not multiple
			secretFields := 0
			if mapping.SecretKey != "" {
				secretFields++
			}
			if mapping.SecretPath != "" {
				secretFields++
			}
			if mapping.Value != nil {
				secretFields++
			}

			if secretFields == 0 {
				return fmt.Errorf("environment '%s': mapping %d: either secret-key, secret-path, or value is required", envName, i+1)
			}

			if secretFields > 1 {
				return fmt.Errorf("environment '%s': mapping %d: cannot specify multiple of secret-key, secret-path, or value", envName, i+1)
			}

			// Validate provider if specified and secret fields are used
			if mapping.Provider != "" && (mapping.SecretKey != "" || mapping.SecretPath != "") {
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
	validProviders := []string{"gcp", "aws", "azure", "openbao"}
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
