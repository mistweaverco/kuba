package secrets

import (
	"context"
	"testing"
)

func TestNewGCPSecretManager(t *testing.T) {
	ctx := context.Background()

	// Test with no credentials file (should use default credentials)
	manager, err := NewGCPSecretManager(ctx, "")
	if err != nil {
		// This might fail if no GCP credentials are available, which is expected in test environment
		t.Logf("Expected error when no GCP credentials available: %v", err)
		return
	}

	if manager == nil {
		t.Error("Expected manager to be created")
	}

	if manager.client == nil {
		t.Error("Expected client to be initialized")
	}

	// Clean up
	if err := manager.Close(); err != nil {
		t.Errorf("Failed to close manager: %v", err)
	}
}

func TestGCPSecretManager_GetSecret(t *testing.T) {
	ctx := context.Background()

	// This test requires actual GCP credentials and a real project
	// In a real test environment, you would mock the client or use test credentials
	t.Skip("Skipping test that requires GCP credentials")

	manager, err := NewGCPSecretManager(ctx, "")
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer manager.Close()

	// Test getting a secret
	secret, err := manager.GetSecret("test-project", "test-secret")
	if err != nil {
		t.Errorf("Failed to get secret: %v", err)
	}

	if secret == "" {
		t.Error("Expected non-empty secret")
	}
}

func TestGCPSecretManager_GetSecrets(t *testing.T) {
	ctx := context.Background()

	// This test requires actual GCP credentials and a real project
	// In a real test environment, you would mock the client or use test credentials
	t.Skip("Skipping test that requires GCP credentials")

	manager, err := NewGCPSecretManager(ctx, "")
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer manager.Close()

	// Test getting multiple secrets
	secretIDs := []string{"secret1", "secret2"}
	secrets, err := manager.GetSecrets("test-project", secretIDs)
	if err != nil {
		t.Errorf("Failed to get secrets: %v", err)
	}

	if len(secrets) != len(secretIDs) {
		t.Errorf("Expected %d secrets, got %d", len(secretIDs), len(secrets))
	}

	for _, secretID := range secretIDs {
		if _, exists := secrets[secretID]; !exists {
			t.Errorf("Expected secret '%s' to exist", secretID)
		}
	}
}
