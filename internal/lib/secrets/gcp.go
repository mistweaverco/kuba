package secrets

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
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
