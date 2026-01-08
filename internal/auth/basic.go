package auth

import (
	"os"

	goauth "github.com/abbot/go-http-auth"
)

// Config holds authentication configuration
type Config struct {
	User  string
	Hash  string
	Realm string
}

// DefaultConfig returns configuration from environment variables with fallbacks
func DefaultConfig() Config {
	return Config{
		User:  getEnvOrDefault("AUTH_USER", "joaquim"),
		Hash:  getEnvOrDefault("AUTH_HASH", "$1$neD1XEAG$WylfbCkcn9psU0o467.AM1"), // default: amazon
		Realm: getEnvOrDefault("AUTH_REALM", "amazon-server-logs.com"),
	}
}

// SecretProvider creates a secret function for the given config
func SecretProvider(cfg Config) func(user, realm string) string {
	return func(user, realm string) string {
		if user == cfg.User {
			return cfg.Hash
		}
		return ""
	}
}

// NewAuthenticator creates a new BasicAuth with the given config
func NewAuthenticator(cfg Config) *goauth.BasicAuth {
	return goauth.NewBasicAuthenticator(cfg.Realm, SecretProvider(cfg))
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
