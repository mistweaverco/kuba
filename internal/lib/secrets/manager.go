package secrets

import (
	"context"
	"fmt"
	"os"

	"github.com/mistweaverco/kuba/internal/config"
)

// SecretManager defines the interface for secret management operations
type SecretManager interface {
	GetSecret(projectID, secretID string) (string, error)
	GetSecrets(projectID string, secretIDs []string) (map[string]string, error)
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
		// TODO: Implement AWS Secrets Manager
		return nil, fmt.Errorf("AWS secrets manager not yet implemented")
	case "azure":
		// TODO: Implement Azure Key Vault
		return nil, fmt.Errorf("Azure Key Vault not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported cloud provider: %s", provider)
	}
}

// GetSecretsForEnvironment retrieves all secrets and values for a given environment configuration
func (f *SecretManagerFactory) GetSecretsForEnvironment(ctx context.Context, env *config.Environment) (map[string]string, error) {
	// Group mappings by provider and project for secret-based mappings
	providerGroups := make(map[string]map[string][]string)

	// Process all mappings to separate secret-based and value-based ones
	for _, mapping := range env.Mappings {
		// Handle direct values first
		if mapping.Value != nil {
			continue // Skip secret processing for value-based mappings
		}

		// Process secret-based mappings
		if mapping.SecretKey != "" {
			provider := mapping.Provider
			if provider == "" {
				provider = env.Provider
			}

			project := mapping.Project
			if project == "" {
				project = env.Project
			}

			if providerGroups[provider] == nil {
				providerGroups[provider] = make(map[string][]string)
			}

			providerGroups[provider][project] = append(providerGroups[provider][project], mapping.SecretKey)
		}
	}

	// Fetch secrets from each provider
	allSecrets := make(map[string]string)

	for provider, projects := range providerGroups {
		for project, secretIDs := range projects {
			secretManager, err := f.CreateSecretManager(ctx, provider, project)
			if err != nil {
				// Log warning but continue with other providers
				fmt.Printf("Warning: failed to create secret manager for %s: %v\n", provider, err)
				continue
			}
			defer secretManager.Close()

			secrets, err := secretManager.GetSecrets(project, secretIDs)
			if err != nil {
				// Log warning but continue with other providers
				fmt.Printf("Warning: failed to get secrets from %s project %s: %v\n", provider, project, err)
				continue
			}

			// Map secrets to environment variables
			for _, mapping := range env.Mappings {
				if mapping.SecretKey != "" {
					if mapping.Provider == provider && mapping.Project == project {
						if secret, exists := secrets[mapping.SecretKey]; exists {
							allSecrets[mapping.EnvironmentVariable] = secret
						}
					} else if mapping.Provider == "" && mapping.Project == "" && env.Provider == provider && env.Project == project {
						if secret, exists := secrets[mapping.SecretKey]; exists {
							allSecrets[mapping.EnvironmentVariable] = secret
						}
					}
				}
			}
		}
	}

	// Process value-based mappings
	for _, mapping := range env.Mappings {
		if mapping.Value != nil {
			// Convert value to string
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
			allSecrets[mapping.EnvironmentVariable] = strValue
		}
	}

	return allSecrets, nil
}
