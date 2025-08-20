package secrets

import (
	"context"
	"testing"
)

func TestNewAWSSecretsManager(t *testing.T) {
	ctx := context.Background()

	// Test with no region or profile (should use default config)
	manager, err := NewAWSSecretsManager(ctx, "", "")
	if err != nil {
		// This might fail if no AWS credentials are available, which is expected in test environment
		t.Logf("Expected error when no AWS credentials available: %v", err)
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

func TestNewAWSSecretsManager_WithRegion(t *testing.T) {
	ctx := context.Background()

	// Test with specific region
	manager, err := NewAWSSecretsManager(ctx, "us-east-1", "")
	if err != nil {
		// This might fail if no AWS credentials are available, which is expected in test environment
		t.Logf("Expected error when no AWS credentials available: %v", err)
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

func TestNewAWSSecretsManager_WithProfile(t *testing.T) {
	ctx := context.Background()

	// Test with specific profile
	manager, err := NewAWSSecretsManager(ctx, "", "test-profile")
	if err != nil {
		// This might fail if no AWS credentials are available, which is expected in test environment
		t.Logf("Expected error when no AWS credentials available: %v", err)
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

func TestAWSSecretsManager_GetSecret(t *testing.T) {
	ctx := context.Background()

	// This test requires actual AWS credentials and a real secret
	// In a real test environment, you would mock the client or use test credentials
	t.Skip("Skipping test that requires AWS credentials")

	manager, err := NewAWSSecretsManager(ctx, "", "")
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer manager.Close()

	// Test getting a secret
	// Note: In AWS, projectID is not used, but we keep the interface consistent
	secret, err := manager.GetSecret("", "test-secret")
	if err != nil {
		t.Errorf("Failed to get secret: %v", err)
	}

	if secret == "" {
		t.Error("Expected non-empty secret")
	}
}

func TestAWSSecretsManager_GetSecrets(t *testing.T) {
	ctx := context.Background()

	// This test requires actual AWS credentials and real secrets
	// In a real test environment, you would mock the client or use test credentials
	t.Skip("Skipping test that requires AWS credentials")

	manager, err := NewAWSSecretsManager(ctx, "", "")
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer manager.Close()

	// Test getting multiple secrets
	// Note: In AWS, projectID is not used, but we keep the interface consistent
	secretIDs := []string{"secret1", "secret2"}
	secrets, err := manager.GetSecrets("", secretIDs)
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

func TestAWSSecretsManager_ListSecrets(t *testing.T) {
	ctx := context.Background()

	// This test requires actual AWS credentials
	// In a real test environment, you would mock the client or use test credentials
	t.Skip("Skipping test that requires AWS credentials")

	manager, err := NewAWSSecretsManager(ctx, "", "")
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer manager.Close()

	// Test listing secrets
	secrets, err := manager.ListSecrets()
	if err != nil {
		t.Errorf("Failed to list secrets: %v", err)
	}

	// We can't assert on the exact number since it depends on the AWS account
	// but we can check that the function doesn't error
	_ = secrets
}

func TestAWSSecretsManager_CreateUpdateDeleteSecret(t *testing.T) {
	ctx := context.Background()

	// This test requires actual AWS credentials and permissions to create/delete secrets
	// In a real test environment, you would mock the client or use test credentials
	t.Skip("Skipping test that requires AWS credentials and permissions")

	manager, err := NewAWSSecretsManager(ctx, "", "")
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer manager.Close()

	secretName := "test-secret-for-crud"
	secretValue := "test-value"
	description := "Test secret for CRUD operations"

	// Test creating a secret
	err = manager.CreateSecret(secretName, secretValue, description)
	if err != nil {
		t.Errorf("Failed to create secret: %v", err)
	}

	// Test updating the secret
	newValue := "updated-test-value"
	err = manager.UpdateSecret(secretName, newValue)
	if err != nil {
		t.Errorf("Failed to update secret: %v", err)
	}

	// Test deleting the secret
	err = manager.DeleteSecret(secretName, true)
	if err != nil {
		t.Errorf("Failed to delete secret: %v", err)
	}
}
