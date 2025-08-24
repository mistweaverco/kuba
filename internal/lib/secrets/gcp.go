package secrets

import (
	"context"
	"fmt"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// GCPSecretManager handles GCP Secret Manager operations
type GCPSecretManager struct {
	client *secretmanager.Client
	ctx    context.Context
}

// NewGCPSecretManager creates a new GCP Secret Manager client
func NewGCPSecretManager(ctx context.Context, credentialsFile string) (*GCPSecretManager, error) {
	var opts []option.ClientOption

	if credentialsFile != "" {
		opts = append(opts, option.WithCredentialsFile(credentialsFile))
	}

	client, err := secretmanager.NewClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret manager client: %w", err)
	}

	return &GCPSecretManager{
		client: client,
		ctx:    ctx,
	}, nil
}

// GetSecret retrieves a secret from GCP Secret Manager
func (g *GCPSecretManager) GetSecret(projectID, secretID string) (string, error) {
	// Build the resource name
	name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretID)

	// Access the secret version
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	result, err := g.client.AccessSecretVersion(g.ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %w", err)
	}

	return string(result.Payload.Data), nil
}

// Close closes the GCP Secret Manager client
func (g *GCPSecretManager) Close() error {
	return g.client.Close()
}

// GetSecrets retrieves multiple secrets from GCP Secret Manager
func (g *GCPSecretManager) GetSecrets(projectID string, secretIDs []string) (map[string]string, error) {
	secrets := make(map[string]string)

	for _, secretID := range secretIDs {
		secret, err := g.GetSecret(projectID, secretID)
		if err != nil {
			return nil, fmt.Errorf("failed to get secret '%s': %w", secretID, err)
		}
		secrets[secretID] = secret
	}

	return secrets, nil
}

// GetSecretsByPath retrieves all secrets that start with the given path prefix
func (g *GCPSecretManager) GetSecretsByPath(projectID, secretPath string) (map[string]string, error) {
	secrets := make(map[string]string)

	// List all secrets in the project
	req := &secretmanagerpb.ListSecretsRequest{
		Parent: fmt.Sprintf("projects/%s", projectID),
	}

	it := g.client.ListSecrets(g.ctx, req)
	for {
		secret, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate secrets: %w", err)
		}

		// Check if the secret name starts with the path prefix
		secretName := secret.Name
		if strings.HasPrefix(secretName, secretPath) {
			// Extract just the secret ID from the full path
			secretID := extractSecretNameFromPath(secretName)

			// Get the actual secret value
			secretValue, err := g.GetSecret(projectID, secretID)
			if err != nil {
				// Log warning but continue with other secrets
				fmt.Printf("Warning: failed to get secret '%s': %v\n", secretID, err)
				continue
			}

			// Sanitize the secret ID for use as an environment variable name
			envVarName := sanitizeEnvVarName(secretID)
			secrets[envVarName] = secretValue
		}
	}

	return secrets, nil
}
