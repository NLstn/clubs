package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&User{}, &APIKey{})
	assert.NoError(t, err)

	return db
}

func TestAPIKey_TableName(t *testing.T) {
	apiKey := APIKey{}
	assert.Equal(t, "api_keys", apiKey.TableName())
}

func TestAPIKey_IsExpired(t *testing.T) {
	tests := []struct {
		name      string
		expiresAt *time.Time
		want      bool
	}{
		{
			name:      "No expiration date",
			expiresAt: nil,
			want:      false,
		},
		{
			name: "Future expiration date",
			expiresAt: func() *time.Time {
				t := time.Now().Add(24 * time.Hour)
				return &t
			}(),
			want: false,
		},
		{
			name: "Past expiration date",
			expiresAt: func() *time.Time {
				t := time.Now().Add(-24 * time.Hour)
				return &t
			}(),
			want: true,
		},
		{
			name: "Expiration date is now",
			expiresAt: func() *time.Time {
				// Set expiration 1 second in the past to account for test execution time
				t := time.Now().Add(-1 * time.Second)
				return &t
			}(),
			want: true, // Expired since it's in the past
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiKey := &APIKey{
				ExpiresAt: tt.expiresAt,
			}
			assert.Equal(t, tt.want, apiKey.IsExpired())
		})
	}
}

func TestAPIKey_IsValid(t *testing.T) {
	tests := []struct {
		name      string
		isActive  bool
		expiresAt *time.Time
		want      bool
	}{
		{
			name:      "Active and not expired",
			isActive:  true,
			expiresAt: nil,
			want:      true,
		},
		{
			name:     "Inactive but not expired",
			isActive: false,
			expiresAt: func() *time.Time {
				t := time.Now().Add(24 * time.Hour)
				return &t
			}(),
			want: false,
		},
		{
			name:     "Active but expired",
			isActive: true,
			expiresAt: func() *time.Time {
				t := time.Now().Add(-24 * time.Hour)
				return &t
			}(),
			want: false,
		},
		{
			name:     "Inactive and expired",
			isActive: false,
			expiresAt: func() *time.Time {
				t := time.Now().Add(-24 * time.Hour)
				return &t
			}(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiKey := &APIKey{
				IsActive:  tt.isActive,
				ExpiresAt: tt.expiresAt,
			}
			assert.Equal(t, tt.want, apiKey.IsValid())
		})
	}
}

func TestAPIKey_BeforeCreate(t *testing.T) {
	t.Skip("Skipping SQLite test - PostgreSQL-specific UUID and text array types are tested in integration tests")
}

func TestAPIKey_Relationships(t *testing.T) {
	t.Skip("Skipping SQLite test - PostgreSQL-specific UUID types are tested in integration tests")
}

func TestAPIKey_Permissions(t *testing.T) {
	// Test permission encoding/decoding without database
	permissions := []string{"read:events", "write:members", "admin"}
	apiKey := &APIKey{}
	
	err := apiKey.SetPermissions(permissions)
	assert.NoError(t, err)
	assert.NotEmpty(t, apiKey.Permissions)
	
	retrieved := apiKey.GetPermissions()
	assert.Equal(t, permissions, retrieved)
}

func TestAPIKey_UniqueKeyHash(t *testing.T) {
	t.Skip("Skipping SQLite test - Unique constraint is tested in integration tests with PostgreSQL")
}
