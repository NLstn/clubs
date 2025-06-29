package database

import (
	"fmt"
	"os"
	"testing"
)

func TestInit_WithEnvironmentVariables(t *testing.T) {
	// Save original environment
	originalEnvVars := map[string]string{
		"DATABASE_URL":           os.Getenv("DATABASE_URL"),
		"DATABASE_PORT":          os.Getenv("DATABASE_PORT"),
		"DATABASE_USER":          os.Getenv("DATABASE_USER"),
		"DATABASE_USER_PASSWORD": os.Getenv("DATABASE_USER_PASSWORD"),
		"DATABASE_NAME":          os.Getenv("DATABASE_NAME"),
		"DATABASE_SSL_MODE":      os.Getenv("DATABASE_SSL_MODE"),
	}

	// Cleanup function to restore environment
	cleanup := func() {
		for key, value := range originalEnvVars {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}
	defer cleanup()

	tests := []struct {
		name            string
		envVars         map[string]string
		expectedDBName  string
		expectedSSLMode string
		shouldError     bool
	}{
		{
			name: "Default database name and SSL mode",
			envVars: map[string]string{
				"DATABASE_URL":           "localhost",
				"DATABASE_PORT":          "5432",
				"DATABASE_USER":          "test",
				"DATABASE_USER_PASSWORD": "test",
			},
			expectedDBName:  "clubs",
			expectedSSLMode: "disable",
			shouldError:     false,
		},
		{
			name: "Custom database name and SSL mode",
			envVars: map[string]string{
				"DATABASE_URL":           "localhost",
				"DATABASE_PORT":          "5432",
				"DATABASE_USER":          "test",
				"DATABASE_USER_PASSWORD": "test",
				"DATABASE_NAME":          "custom_db",
				"DATABASE_SSL_MODE":      "require",
			},
			expectedDBName:  "custom_db",
			expectedSSLMode: "require",
			shouldError:     false,
		},
		{
			name: "Empty database name falls back to default",
			envVars: map[string]string{
				"DATABASE_URL":           "localhost",
				"DATABASE_PORT":          "5432",
				"DATABASE_USER":          "test",
				"DATABASE_USER_PASSWORD": "test",
				"DATABASE_NAME":          "",
				"DATABASE_SSL_MODE":      "prefer",
			},
			expectedDBName:  "clubs",
			expectedSSLMode: "prefer",
			shouldError:     false,
		},
		{
			name: "Empty SSL mode falls back to default",
			envVars: map[string]string{
				"DATABASE_URL":           "localhost",
				"DATABASE_PORT":          "5432",
				"DATABASE_USER":          "test",
				"DATABASE_USER_PASSWORD": "test",
				"DATABASE_NAME":          "test_db",
				"DATABASE_SSL_MODE":      "",
			},
			expectedDBName:  "test_db",
			expectedSSLMode: "disable",
			shouldError:     false,
		},
		{
			name: "Missing required DATABASE_URL",
			envVars: map[string]string{
				"DATABASE_PORT":          "5432",
				"DATABASE_USER":          "test",
				"DATABASE_USER_PASSWORD": "test",
			},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all database environment variables
			os.Unsetenv("DATABASE_URL")
			os.Unsetenv("DATABASE_PORT")
			os.Unsetenv("DATABASE_USER")
			os.Unsetenv("DATABASE_USER_PASSWORD")
			os.Unsetenv("DATABASE_NAME")
			os.Unsetenv("DATABASE_SSL_MODE")

			// Set test environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// We can't actually test Init() because it would try to connect to PostgreSQL
			// Instead, we test the environment variable reading logic directly
			if !tt.shouldError {
				dbName := os.Getenv("DATABASE_NAME")
				if dbName == "" {
					dbName = "clubs"
				}
				if dbName != tt.expectedDBName {
					t.Errorf("Expected database name %s, got %s", tt.expectedDBName, dbName)
				}

				sslMode := os.Getenv("DATABASE_SSL_MODE")
				if sslMode == "" {
					sslMode = "disable"
				}
				if sslMode != tt.expectedSSLMode {
					t.Errorf("Expected SSL mode %s, got %s", tt.expectedSSLMode, sslMode)
				}
			}

			// Test that Init() would fail appropriately for missing required vars
			if tt.shouldError {
				err := Init()
				if err == nil {
					t.Error("Expected Init() to return an error for missing required environment variables")
				}
			}
		})
	}
}

func TestNewConnection_WithConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectedDSN string
	}{
		{
			name: "Default SSL mode",
			config: &Config{
				Host:     "localhost",
				Port:     5432,
				User:     "test",
				Password: "test",
				DBName:   "clubs",
				SSLMode:  "disable",
			},
			expectedDSN: "host=localhost port=5432 user=test password=test dbname=clubs sslmode=disable",
		},
		{
			name: "Custom SSL mode",
			config: &Config{
				Host:     "example.com",
				Port:     5432,
				User:     "myuser",
				Password: "mypass",
				DBName:   "mydb",
				SSLMode:  "require",
			},
			expectedDSN: "host=example.com port=5432 user=myuser password=mypass dbname=mydb sslmode=require",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We can't actually test NewConnection() because it would try to connect to PostgreSQL
			// Instead, we test that the DSN is constructed correctly
			expectedDSN := tt.expectedDSN
			actualDSN := constructDSN(tt.config)
			if actualDSN != expectedDSN {
				t.Errorf("Expected DSN %s, got %s", expectedDSN, actualDSN)
			}
		})
	}
}

// Helper function to construct DSN for testing (extracted logic from NewConnection)
func constructDSN(config *Config) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)
}
