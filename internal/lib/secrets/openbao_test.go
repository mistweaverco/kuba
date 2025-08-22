package secrets

import (
	"context"
	"testing"
)

func TestNewOpenBaoManager(t *testing.T) {
	ctx := context.Background()

	// Test with minimal configuration
	manager, err := NewOpenBaoManager(ctx, "http://localhost:8200", "", "")
	if err != nil {
		t.Fatalf("Failed to create OpenBao manager: %v", err)
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

func TestNewOpenBaoManager_WithToken(t *testing.T) {
	ctx := context.Background()

	// Test with token
	manager, err := NewOpenBaoManager(ctx, "http://localhost:8200", "test-token", "")
	if err != nil {
		t.Fatalf("Failed to create OpenBao manager: %v", err)
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

func TestNewOpenBaoManager_WithNamespace(t *testing.T) {
	ctx := context.Background()

	// Test with namespace
	manager, err := NewOpenBaoManager(ctx, "http://localhost:8200", "test-token", "test-namespace")
	if err != nil {
		t.Fatalf("Failed to create OpenBao manager: %v", err)
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

func TestOpenBaoManager_GetSecret(t *testing.T) {
	ctx := context.Background()

	// This test requires an actual OpenBao server running
	// In a real test environment, you would mock the client or use a test server
	t.Skip("Skipping test that requires OpenBao server")

	manager, err := NewOpenBaoManager(ctx, "http://localhost:8200", "test-token", "")
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer manager.Close()

	// Test getting a secret
	secret, err := manager.GetSecret("", "secret/test")
	if err != nil {
		t.Errorf("Failed to get secret: %v", err)
		return
	}

	if secret == "" {
		t.Error("Expected secret to have a value")
	}
}

func TestOpenBaoManager_GetSecrets(t *testing.T) {
	ctx := context.Background()

	// This test requires an actual OpenBao server running
	// In a real test environment, you would mock the client or use a test server
	t.Skip("Skipping test that requires OpenBao server")

	manager, err := NewOpenBaoManager(ctx, "http://localhost:8200", "test-token", "")
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer manager.Close()

	// Test getting multiple secrets
	secretIDs := []string{"secret/test1", "secret/test2"}
	secrets, err := manager.GetSecrets("", secretIDs)
	if err != nil {
		t.Errorf("Failed to get secrets: %v", err)
		return
	}

	if len(secrets) != len(secretIDs) {
		t.Errorf("Expected %d secrets, got %d", len(secretIDs), len(secrets))
	}

	for _, secretID := range secretIDs {
		if _, exists := secrets[secretID]; !exists {
			t.Errorf("Expected secret '%s' to be in results", secretID)
		}
	}
}

func TestOpenBaoManager_ListSecrets(t *testing.T) {
	ctx := context.Background()

	// This test requires an actual OpenBao server running
	// In a real test environment, you would mock the client or use a test server
	t.Skip("Skipping test that requires OpenBao server")

	manager, err := NewOpenBaoManager(ctx, "http://localhost:8200", "test-token", "")
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer manager.Close()

	// Test listing secrets
	secrets, err := manager.ListSecrets("secret/")
	if err != nil {
		t.Errorf("Failed to list secrets: %v", err)
		return
	}

	// Should return a list (even if empty)
	if secrets == nil {
		t.Error("Expected secrets list to not be nil")
	}
}

func TestOpenBaoManager_CreateSecret(t *testing.T) {
	ctx := context.Background()

	// This test requires an actual OpenBao server running
	// In a real test environment, you would mock the client or use a test server
	t.Skip("Skipping test that requires OpenBao server")

	manager, err := NewOpenBaoManager(ctx, "http://localhost:8200", "test-token", "")
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer manager.Close()

	// Test creating a secret
	data := map[string]interface{}{
		"password": "test-password",
		"username": "test-user",
	}

	err = manager.CreateSecret("secret/test-create", data)
	if err != nil {
		t.Errorf("Failed to create secret: %v", err)
	}
}

func TestOpenBaoManager_UpdateSecret(t *testing.T) {
	ctx := context.Background()

	// This test requires an actual OpenBao server running
	// In a real test environment, you would mock the client or use a test server
	t.Skip("Skipping test that requires OpenBao server")

	manager, err := NewOpenBaoManager(ctx, "http://localhost:8200", "test-token", "")
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer manager.Close()

	// Test updating a secret
	data := map[string]interface{}{
		"password": "updated-password",
		"username": "updated-user",
	}

	err = manager.UpdateSecret("secret/test-update", data)
	if err != nil {
		t.Errorf("Failed to update secret: %v", err)
	}
}

func TestOpenBaoManager_DeleteSecret(t *testing.T) {
	ctx := context.Background()

	// This test requires an actual OpenBao server running
	// In a real test environment, you would mock the client or use a test server
	t.Skip("Skipping test that requires OpenBao server")

	manager, err := NewOpenBaoManager(ctx, "http://localhost:8200", "test-token", "")
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer manager.Close()

	// Test deleting a secret
	err = manager.DeleteSecret("secret/test-delete")
	if err != nil {
		t.Errorf("Failed to delete secret: %v", err)
	}
}
