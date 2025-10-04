package secrets

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mistweaverco/kuba/internal/config"
	"github.com/mistweaverco/kuba/internal/lib/log"
)

// SecretManager defines the interface for secret management operations
type SecretManager interface {
	GetSecret(projectID, secretID string) (string, error)
	GetSecrets(projectID string, secretIDs []string) (map[string]string, error)
	GetSecretsByPath(projectID, secretPath string) (map[string]string, error)
	Close() error
}

// SecretManagerFactory creates secret managers for different cloud providers
type SecretManagerFactory struct{}

// NewSecretManagerFactory creates a new secret manager factory
func NewSecretManagerFactory() *SecretManagerFactory {
	return &SecretManagerFactory{}
}

// CreateSecretManager creates a secret manager for the specified provider
func (f *SecretManagerFactory) CreateSecretManager(ctx context.Context, provider string, projectID string) (SecretManager, error) {
	switch provider {
	case "gcp":
		// Check for GCP credentials
		credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
		return NewGCPSecretManager(ctx, credentialsFile)
	case "aws":
		// Check for AWS region and profile
		region := os.Getenv("AWS_REGION")
		profile := os.Getenv("AWS_PROFILE")
		return NewAWSSecretsManager(ctx, region, profile)
	case "azure":
		// Check for Azure Key Vault configuration
		vaultURL := os.Getenv("AZURE_KEY_VAULT_URL")
		if vaultURL == "" {
			return nil, fmt.Errorf("AZURE_KEY_VAULT_URL environment variable is required for Azure Key Vault")
		}

		// Optional: tenant ID, client ID, and client secret for service principal auth
		tenantID := os.Getenv("AZURE_TENANT_ID")
		clientID := os.Getenv("AZURE_CLIENT_ID")
		clientSecret := os.Getenv("AZURE_CLIENT_SECRET")

		return NewAzureKeyVaultManager(ctx, vaultURL, tenantID, clientID, clientSecret)
	case "openbao":
		// Check for OpenBao configuration
		address := os.Getenv("OPENBAO_ADDR")
		if address == "" {
			return nil, fmt.Errorf("OPENBAO_ADDR environment variable is required for OpenBao")
		}

		// Optional: token and namespace
		token := os.Getenv("OPENBAO_TOKEN")
		namespace := os.Getenv("OPENBAO_NAMESPACE")

		return NewOpenBaoManager(ctx, address, token, namespace)
	case "local":
		// Local provider doesn't require any external configuration
		return NewLocalManager(ctx)
	default:
		return nil, fmt.Errorf("unsupported cloud provider: %s", provider)
	}
}

// GetSecretsForEnvironment retrieves all secrets and values for a given environment configuration
func (f *SecretManagerFactory) GetSecretsForEnvironment(ctx context.Context, env *config.Environment) (map[string]string, error) {
	logger := log.NewLogger()

	// Group mappings by provider and project for secret-based mappings
	providerGroups := make(map[string]map[string][]string)

	// Group mappings by provider and project for path-based mappings
	pathGroups := make(map[string]map[string]string)

	// Get all env items (from map)
	envItems := env.GetEnvItems()
	logger.Debug("Processing environment mappings", "total_mappings", len(envItems))

	// Process all env items to separate secret-based and value-based ones
	for i, envItem := range envItems {
		logger.Debug("Processing mapping", "index", i, "env_var", envItem.EnvironmentVariable, "has_secret_key", envItem.SecretKey != "", "has_secret_path", envItem.SecretPath != "", "has_value", envItem.Value != nil)

		// Handle direct values first
		if envItem.Value != nil {
			logger.Debug("Skipping secret processing for value-based mapping", "env_var", envItem.EnvironmentVariable)
			continue // Skip secret processing for value-based mappings
		}

		// Process secret-based mappings (single key)
		if envItem.SecretKey != "" {
			provider := envItem.Provider
			if provider == "" {
				provider = env.Provider
			}

			project := envItem.Project
			if project == "" {
				project = env.Project
			}

			// For AWS, Azure, OpenBao, and local, we use a default project key since they don't use projects in the same way as GCP
			if (provider == "aws" || provider == "azure" || provider == "openbao" || provider == "local") && project == "" {
				project = "default"
			}

			logger.Debug("Adding secret-based mapping to provider group", "provider", provider, "project", project, "secret_key", envItem.SecretKey)

			if providerGroups[provider] == nil {
				providerGroups[provider] = make(map[string][]string)
			}

			providerGroups[provider][project] = append(providerGroups[provider][project], envItem.SecretKey)
		}

		// Process path-based mappings
		if envItem.SecretPath != "" {
			provider := envItem.Provider
			if provider == "" {
				provider = env.Provider
			}

			project := envItem.Project
			if project == "" {
				project = env.Project
			}

			// For AWS, Azure, OpenBao, and local, we use a default project key since they don't use projects in the same way as GCP
			if (provider == "aws" || provider == "azure" || provider == "openbao" || provider == "local") && project == "" {
				project = "default"
			}

			logger.Debug("Adding path-based mapping to provider group", "provider", provider, "project", project, "secret_path", envItem.SecretPath)

			// Create a separate group for path-based lookups
			pathKey := fmt.Sprintf("%s:%s", provider, project)
			if pathGroups[pathKey] == nil {
				pathGroups[pathKey] = make(map[string]string)
			}
			pathGroups[pathKey][envItem.EnvironmentVariable] = envItem.SecretPath
		}
	}

	logger.Debug("Provider groups created", "secret_providers", len(providerGroups), "path_providers", len(pathGroups))

	// Fetch secrets from each provider
	allSecrets := make(map[string]string)

	for provider, projects := range providerGroups {
		for project, secretIDs := range projects {
			logger.Debug("Creating secret manager", "provider", provider, "project", project, "secret_count", len(secretIDs))

			secretManager, err := f.CreateSecretManager(ctx, provider, project)
			if err != nil {
				logger.Debug("Failed to create secret manager", "provider", provider, "project", project, "error", err)
				// Log warning but continue with other providers
				fmt.Printf("Warning: failed to create secret manager for %s: %v\n", provider, err)
				continue
			}
			defer secretManager.Close()

			logger.Debug("Fetching secrets from provider", "provider", provider, "project", project, "secret_ids", secretIDs)
			secrets, err := secretManager.GetSecrets(project, secretIDs)
			if err != nil {
				logger.Debug("Failed to get secrets from provider", "provider", provider, "project", project, "error", err)
				// Log warning but continue with other providers
				fmt.Printf("Warning: failed to get secrets from %s project %s: %v\n", provider, project, err)
				continue
			}

			logger.Debug("Successfully retrieved secrets from provider", "provider", provider, "project", project, "retrieved_count", len(secrets))

			// Map secrets to environment variables
			for _, envItem := range envItems {
				if envItem.SecretKey != "" {
					envItemProvider := envItem.Provider
					if envItemProvider == "" {
						envItemProvider = env.Provider
					}

					envItemProject := envItem.Project
					if envItemProject == "" {
						envItemProject = env.Project
					}

					// For AWS, Azure, OpenBao, and local, we use a default project key since they don't use projects in the same way as GCP
					if (envItemProvider == "aws" || envItemProvider == "azure" || envItemProvider == "openbao" || envItemProvider == "local") && envItemProject == "" {
						envItemProject = "default"
					}

					// Only process mappings that match the current provider and project
					if envItemProvider == provider && envItemProject == project {
						if secretValue, exists := secrets[envItem.SecretKey]; exists {
							allSecrets[envItem.EnvironmentVariable] = secretValue
							logger.Debug("Mapped secret to environment variable", "env_var", envItem.EnvironmentVariable, "secret_key", envItem.SecretKey, "provider", provider, "project", project)
						} else {
							logger.Debug("Secret key not found in provider response", "env_var", envItem.EnvironmentVariable, "secret_key", envItem.SecretKey, "provider", provider, "project", project)
						}
					}
				}
			}
		}
	}

	// Process path-based mappings
	for pathKey, pathMappings := range pathGroups {
		// Parse the path key to get provider and project
		parts := strings.Split(pathKey, ":")
		if len(parts) != 2 {
			fmt.Printf("Warning: invalid path key format: %s\n", pathKey)
			continue
		}

		provider := parts[0]
		project := parts[1]

		secretManager, err := f.CreateSecretManager(ctx, provider, project)
		if err != nil {
			// Log warning but continue with other providers
			fmt.Printf("Warning: failed to create secret manager for %s: %v\n", provider, err)
			continue
		}
		defer secretManager.Close()

		// Process each path mapping
		for envVar, secretPath := range pathMappings {
			secrets, err := secretManager.GetSecretsByPath(project, secretPath)
			if err != nil {
				// Log warning but continue with other paths
				fmt.Printf("Warning: failed to get secrets from path '%s': %v\n", secretPath, err)
				continue
			}

			// Add all secrets from this path to the result
			// The environment variable name from the mapping is used as a prefix
			for secretName, secretValue := range secrets {
				// Create a unique environment variable name by combining the mapping's env var and the secret name
				finalEnvVarName := envVar + "_" + secretName
				allSecrets[finalEnvVarName] = secretValue
			}
		}
	}

	// Process value-based mappings (no bare items allowed anymore)
	for _, envItem := range envItems {
		if envItem.Value != nil {
			// Convert value to string
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
			allSecrets[envItem.EnvironmentVariable] = strValue
		}
	}

	// Perform interpolation on all values now that we have all secrets and values
	// This allows values to reference other environment variables that were just resolved
	for key, value := range allSecrets {
		if strings.Contains(value, "${") {
			interpolatedValue := config.InterpolateEnvVars(value, allSecrets)
			allSecrets[key] = interpolatedValue
		}
	}

	return allSecrets, nil
}
