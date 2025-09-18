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
	Provider string             `yaml:"provider"`
	Project  string             `yaml:"project"`
	Env      map[string]EnvItem `yaml:"env"`
}

// EnvItem represents an environment variable configuration in the new simplified format
// It can be either a string (just the env var name) or a full mapping object
type EnvItem struct {
	// For string format: just the environment variable name
	EnvironmentVariable string      `yaml:"environment-variable,omitempty"`
	SecretKey           string      `yaml:"secret-key,omitempty"`
	SecretPath          string      `yaml:"secret-path,omitempty"`
	Value               interface{} `yaml:"value,omitempty"`
	Provider            string      `yaml:"provider,omitempty"`
	Project             string      `yaml:"project,omitempty"`
}

// UnmarshalYAML implements custom YAML unmarshaling for EnvItem
// This allows it to handle both string format (just env var name) and object format
func (e *EnvItem) UnmarshalYAML(value *yaml.Node) error {
	// For map syntax, the env var name is the map key; object holds fields only
	var temp struct {
		SecretKey  string      `yaml:"secret-key,omitempty"`
		SecretPath string      `yaml:"secret-path,omitempty"`
		Value      interface{} `yaml:"value,omitempty"`
		Provider   string      `yaml:"provider,omitempty"`
		Project    string      `yaml:"project,omitempty"`
	}
	if err := value.Decode(&temp); err != nil {
		return err
	}
	e.SecretKey = temp.SecretKey
	e.SecretPath = temp.SecretPath
	e.Value = temp.Value
	e.Provider = temp.Provider
	e.Project = temp.Project
	return nil
}

// GetEnvItems returns all env items for an environment
func (e *Environment) GetEnvItems() []EnvItem {
	items := make([]EnvItem, 0, len(e.Env))
	for name, item := range e.Env {
		item.EnvironmentVariable = name
		items = append(items, item)
	}
	return items
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

// processValueInterpolations processes all value fields in env items to resolve environment variable interpolations
func processValueInterpolations(config *KubaConfig) error {
	// Process environments in order to handle dependencies correctly
	// We'll process each environment multiple times until no more interpolations are possible
	// or until we detect a circular dependency

	for envName, env := range config.Environments {
		// Track resolved variables for this environment
		resolvedVars := make(map[string]string)

		// Process env items multiple times to handle dependencies
		maxIterations := len(env.Env) * 2 // Allow for some dependency depth
		for iteration := 0; iteration < maxIterations; iteration++ {
			changed := false

			// Process env items
			for name, envItem := range env.Env {
				if envItem.Value != nil {
					// Convert value to string for processing
					var strValue string
					switch v := envItem.Value.(type) {
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
							// Update the env item value
							tmp := env.Env[name]
							tmp.Value = interpolatedValue
							env.Env[name] = tmp
							// Update our resolved vars map
							resolvedVars[name] = interpolatedValue
							changed = true
						}
					} else {
						// No interpolation needed, but convert numeric values to strings for consistency
						if envItem.Value != strValue {
							tmp := env.Env[name]
							tmp.Value = strValue
							env.Env[name] = tmp
							changed = true
						}
						// Store the value in resolved vars
						resolvedVars[name] = strValue
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

	logger.Debug("Environment configuration retrieved", "environment", envName, "provider", env.Provider, "project", env.Project, "env_count", len(env.Env))
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

		// Project is required for all providers except AWS, Azure, OpenBao, and local
		if env.Project == "" && env.Provider != "aws" && env.Provider != "azure" && env.Provider != "openbao" && env.Provider != "local" {
			return fmt.Errorf("environment '%s': project is required for provider '%s'", envName, env.Provider)
		}

		// At least one env item must be provided
		if len(env.Env) == 0 {
			return fmt.Errorf("environment '%s': at least one env item is required", envName)
		}

		// Validate env items
		idx := 0
		for _, envItem := range env.Env {
			idx++
			// name is the environment variable

			// Either secret-key, secret-path, or value must be provided (no bare items)
			// Special case: for local provider (env-level or item-level), only value is allowed
			secretFields := 0
			if envItem.SecretKey != "" {
				secretFields++
			}
			if envItem.SecretPath != "" {
				secretFields++
			}
			if envItem.Value != nil {
				secretFields++
			}

			if secretFields == 0 {
				return fmt.Errorf("environment '%s': env item %d: either secret-key, secret-path, or value is required", envName, idx)
			}

			if secretFields > 1 {
				return fmt.Errorf("environment '%s': env item %d: cannot specify multiple of secret-key, secret-path, or value", envName, idx)
			}

			// Determine effective provider for this item
			effectiveProvider := env.Provider
			if envItem.Provider != "" {
				effectiveProvider = envItem.Provider
			}

			// Validate provider value if set on item
			if envItem.Provider != "" && !isValidProvider(envItem.Provider) {
				return fmt.Errorf("environment '%s': env item %d: invalid provider '%s'", envName, idx, envItem.Provider)
			}

			// Local provider rules: only value is allowed
			if effectiveProvider == "local" {
				if envItem.Value == nil {
					return fmt.Errorf("environment '%s': env item %d: provider 'local' requires 'value'", envName, idx)
				}
				if envItem.SecretKey != "" || envItem.SecretPath != "" {
					return fmt.Errorf("environment '%s': env item %d: provider 'local' does not support 'secret-key' or 'secret-path'", envName, idx)
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
	validProviders := []string{"gcp", "aws", "azure", "openbao", "local"}
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
