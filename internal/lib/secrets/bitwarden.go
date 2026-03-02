package secrets

import (
	"context"
	"fmt"
	"os"

	sdk "github.com/bitwarden/sdk-go/v2"
)

// BitwardenManager handles Bitwarden Secrets Manager operations
type BitwardenManager struct {
	client         sdk.BitwardenClientInterface
	ctx            context.Context
	organizationID string
}

// NewBitwardenManager creates a new Bitwarden client.
// projectID is treated as the Bitwarden organization ID. If empty, the
// organization ID is read from the BITWARDEN_ORGANIZATION_ID environment
// variable.
//
// Authentication is performed using a Bitwarden secrets access token from
// BITWARDEN_ACCESS_TOKEN (or ACCESS_TOKEN as a fallback). Optional settings:
//   - BITWARDEN_API_URL / BITWARDEN_IDENTITY_URL for self-hosted instances
//   - BITWARDEN_STATE_FILE for a persisted state file path
func NewBitwardenManager(ctx context.Context, projectID string) (*BitwardenManager, error) {
	organizationID := projectID
	if organizationID == "" {
		organizationID = os.Getenv("BITWARDEN_ORGANIZATION_ID")
	}
	if organizationID == "" {
		return nil, fmt.Errorf("Bitwarden organization ID is required (set env BITWARDEN_ORGANIZATION_ID or configure 'project' for the environment)")
	}

	accessToken := os.Getenv("BITWARDEN_ACCESS_TOKEN")
	if accessToken == "" {
		accessToken = os.Getenv("ACCESS_TOKEN")
	}
	if accessToken == "" {
		return nil, fmt.Errorf("Bitwarden access token is required (set env BITWARDEN_ACCESS_TOKEN or ACCESS_TOKEN)")
	}

	var apiURLPtr *string
	if apiURL := os.Getenv("BITWARDEN_API_URL"); apiURL != "" {
		apiURLCopy := apiURL
		apiURLPtr = &apiURLCopy
	}

	var identityURLPtr *string
	if identityURL := os.Getenv("BITWARDEN_IDENTITY_URL"); identityURL != "" {
		identityURLCopy := identityURL
		identityURLPtr = &identityURLCopy
	}

	var stateFilePtr *string
	if stateFile := os.Getenv("BITWARDEN_STATE_FILE"); stateFile != "" {
		stateFileCopy := stateFile
		stateFilePtr = &stateFileCopy
	}

	client, err := sdk.NewBitwardenClient(apiURLPtr, identityURLPtr)
	if err != nil {
		return nil, fmt.Errorf("failed to create Bitwarden client: %w", err)
	}

	if err := client.AccessTokenLogin(accessToken, stateFilePtr); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to authenticate with Bitwarden: %w", err)
	}

	return &BitwardenManager{
		client:         client,
		ctx:            ctx,
		organizationID: organizationID,
	}, nil
}

// GetSecret retrieves a single secret from Bitwarden by its ID.
// projectID is ignored because Bitwarden secrets are addressed by ID.
func (b *BitwardenManager) GetSecret(projectID, secretID string) (string, error) {
	secret, err := b.client.Secrets().Get(secretID)
	if err != nil {
		return "", fmt.Errorf("failed to get Bitwarden secret '%s': %w", secretID, err)
	}
	return secret.Value, nil
}

// GetSecrets retrieves multiple secrets from Bitwarden by their IDs.
// projectID is ignored because Bitwarden secrets are addressed by ID.
func (b *BitwardenManager) GetSecrets(projectID string, secretIDs []string) (map[string]string, error) {
	secrets := make(map[string]string)

	for _, id := range secretIDs {
		value, err := b.GetSecret(projectID, id)
		if err != nil {
			return nil, err
		}
		secrets[id] = value
	}

	return secrets, nil
}

// GetSecretsByPath is currently not supported for Bitwarden.
// Bitwarden does not organize secrets by hierarchical paths, so we return
// an error to indicate that path-based mappings are not available.
func (b *BitwardenManager) GetSecretsByPath(projectID, secretPath string) (map[string]string, error) {
	return nil, fmt.Errorf("Bitwarden provider does not support 'secret-path' mappings")
}

// Close closes the Bitwarden client and releases native resources.
func (b *BitwardenManager) Close() error {
	if b.client != nil {
		b.client.Close()
	}
	return nil
}
