package auth

import (
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.User != "joaquim" {
		t.Errorf("expected default user 'joaquim', got '%s'", cfg.User)
	}

	if cfg.Realm != "amazon-server-logs.com" {
		t.Errorf("expected default realm 'amazon-server-logs.com', got '%s'", cfg.Realm)
	}

	if cfg.Hash == "" {
		t.Error("expected non-empty default hash")
	}
}

func TestDefaultConfigWithEnvVars(t *testing.T) {
	// Set env vars
	os.Setenv("AUTH_USER", "testuser")
	os.Setenv("AUTH_HASH", "$1$testhash")
	os.Setenv("AUTH_REALM", "test.realm.com")
	defer func() {
		os.Unsetenv("AUTH_USER")
		os.Unsetenv("AUTH_HASH")
		os.Unsetenv("AUTH_REALM")
	}()

	cfg := DefaultConfig()

	if cfg.User != "testuser" {
		t.Errorf("expected user 'testuser', got '%s'", cfg.User)
	}

	if cfg.Hash != "$1$testhash" {
		t.Errorf("expected hash '$1$testhash', got '%s'", cfg.Hash)
	}

	if cfg.Realm != "test.realm.com" {
		t.Errorf("expected realm 'test.realm.com', got '%s'", cfg.Realm)
	}
}

func TestSecretProvider(t *testing.T) {
	cfg := Config{
		User:  "testuser",
		Hash:  "$1$testhash",
		Realm: "test.realm.com",
	}

	secretFn := SecretProvider(cfg)

	// Valid user
	result := secretFn("testuser", "test.realm.com")
	if result != "$1$testhash" {
		t.Errorf("expected hash for valid user, got '%s'", result)
	}

	// Invalid user
	result = secretFn("invaliduser", "test.realm.com")
	if result != "" {
		t.Errorf("expected empty string for invalid user, got '%s'", result)
	}
}

func TestNewAuthenticator(t *testing.T) {
	cfg := DefaultConfig()
	auth := NewAuthenticator(cfg)

	if auth == nil {
		t.Error("expected non-nil authenticator")
	}
}
