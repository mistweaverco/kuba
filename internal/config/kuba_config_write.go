package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// AddOrUpdateEnvSecretKeyMapping adds/updates an env mapping in kuba.yaml for the
// given environment and environment variable name.
//
// This is intentionally simple (yaml marshal/unmarshal) and may reorder keys.
func AddOrUpdateEnvSecretKeyMapping(configPath, envName, envVar, secretKey string) error {
	if configPath == "" {
		return fmt.Errorf("configPath is required")
	}
	if envName == "" {
		envName = "default"
	}
	if envVar == "" {
		return fmt.Errorf("envVar is required")
	}
	if secretKey == "" {
		return fmt.Errorf("secretKey is required")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg KubaConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	env, ok := cfg.Environments[envName]
	if !ok {
		return fmt.Errorf("environment '%s' not found in config", envName)
	}

	if env.Env == nil {
		env.Env = map[string]EnvItem{}
	}

	env.Env[envVar] = EnvItem{
		SecretKey: secretKey,
	}
	cfg.Environments[envName] = env

	out, err := yaml.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal updated config: %w", err)
	}

	if err := os.WriteFile(configPath, out, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// RemoveEnvMapping removes a single env mapping (by env var name) from the given
// environment in kuba.yaml.
//
// This is intentionally simple (yaml marshal/unmarshal) and may reorder keys.
func RemoveEnvMapping(configPath, envName, envVar string) error {
	if configPath == "" {
		return fmt.Errorf("configPath is required")
	}
	if envName == "" {
		envName = "default"
	}
	if envVar == "" {
		return fmt.Errorf("envVar is required")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg KubaConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	env, ok := cfg.Environments[envName]
	if !ok {
		return fmt.Errorf("environment '%s' not found in config", envName)
	}

	if env.Env == nil {
		return nil
	}

	delete(env.Env, envVar)
	cfg.Environments[envName] = env

	out, err := yaml.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal updated config: %w", err)
	}

	if err := os.WriteFile(configPath, out, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

