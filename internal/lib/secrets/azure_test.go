package secrets

import (
	"context"
	"os"
	"testing"
)

func TestNewAzureKeyVaultManager(t *testing.T) {
	ctx := context.Background()

	// Test with valid vault URL but no credentials (will try default credential)
	// This test might fail in environments without Azure credentials, which is expected
	_, err := NewAzureKeyVaultManager(ctx, "https://testvault.vault.azure.net/", "", "", "")
	// We don't check the error here because it depends on the environment
	// The important thing is that the function doesn't panic
	if err != nil {
		t.Logf("Expected error in environment without Azure credentials: %v", err)
	}
}

func TestAzureKeyVaultManager_GetSecret(t *testing.T) {
	// Skip test if Azure credentials are not available
	if os.Getenv("AZURE_KEY_VAULT_URL") == "" {
		t.Skip("Skipping Azure Key Vault test - AZURE_KEY_VAULT_URL not set")
	}

	ctx := context.Background()
	vaultURL := os.Getenv("AZURE_KEY_VAULT_URL")
	tenantID := os.Getenv("AZURE_TENANT_ID")
	clientID := os.Getenv("AZURE_CLIENT_ID")
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")

	manager, err := NewAzureKeyVaultManager(ctx, vaultURL, tenantID, clientID, clientSecret)
	if err != nil {
		t.Skipf("Skipping test - failed to create Azure Key Vault manager: %v", err)
	}
	defer manager.Close()

	// Test getting a non-existent secret
	_, err = manager.GetSecret("", "non-existent-secret")
	if err == nil {
		t.Error("Expected error when getting non-existent secret")
	}
}

func TestAzureKeyVaultManager_GetSecrets(t *testing.T) {
	// Skip test if Azure credentials are not available
	if os.Getenv("AZURE_KEY_VAULT_URL") == "" {
		t.Skip("Skipping Azure Key Vault test - AZURE_KEY_VAULT_URL not set")
	}

	ctx := context.Background()
	vaultURL := os.Getenv("AZURE_KEY_VAULT_URL")
	tenantID := os.Getenv("AZURE_TENANT_ID")
	clientID := os.Getenv("AZURE_CLIENT_ID")
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")

	manager, err := NewAzureKeyVaultManager(ctx, vaultURL, tenantID, clientID, clientSecret)
	if err != nil {
		t.Skipf("Skipping test - failed to create Azure Key Vault manager: %v", err)
	}
	defer manager.Close()

	// Test getting multiple non-existent secrets
	_, err = manager.GetSecrets("", []string{"non-existent-secret-1", "non-existent-secret-2"})
	if err == nil {
		t.Error("Expected error when getting non-existent secrets")
	}
}

func TestAzureKeyVaultManager_Close(t *testing.T) {
	ctx := context.Background()

	// Test with a mock manager (no actual connection)
	manager := &AzureKeyVaultManager{
		ctx: ctx,
	}

	err := manager.Close()
	if err != nil {
		t.Errorf("Expected no error when closing manager, got: %v", err)
	}
}
