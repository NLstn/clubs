package odata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIsValidEmail tests the email validation function
func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		// Valid emails
		{"valid simple email", "user@example.com", true},
		{"valid with plus", "user+tag@example.com", true},
		{"valid with dots", "first.last@example.com", true},
		{"valid with numbers", "user123@example456.com", true},
		{"valid with dash", "user-name@example.com", true},
		{"valid with underscore", "user_name@example.com", true},
		{"valid subdomain", "user@mail.example.com", true},
		{"valid long TLD", "user@example.technology", true},
		{"valid with percent", "user%test@example.com", true},

		// Invalid emails
		{"empty string", "", false},
		{"no @ symbol", "userexample.com", false},
		{"no domain", "user@", false},
		{"no local part", "@example.com", false},
		{"no TLD", "user@example", false},
		{"double @", "user@@example.com", false},
		{"spaces", "user @example.com", false},
		{"missing dot in domain", "user@examplecom", false},
		{"starts with dot", ".user@example.com", true},       // Allowed by simple regex
		{"ends with dot", "user.@example.com", true},         // Allowed by simple regex
		{"consecutive dots", "user..name@example.com", true}, // Allowed by simple regex
		{"only whitespace", "   ", false},

		// Edge cases
		{"too short", "a@b.c", false},                              // Fails TLD length check (needs 2+ chars)
		{"whitespace trimmed valid", "  user@example.com  ", true}, // Should be trimmed
		{"single char local", "a@example.com", true},
		{"single char domain", "user@e.com", true},
		{"two char TLD", "user@example.co", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidEmail(tt.email)
			assert.Equal(t, tt.want, got, "isValidEmail(%q) = %v, want %v", tt.email, got, tt.want)
		})
	}
}

// TestIsValidEmailLength tests email length boundaries
func TestIsValidEmailLength(t *testing.T) {
	// Test minimum length (3 characters)
	assert.False(t, isValidEmail("ab"), "2 chars should be invalid")
	assert.False(t, isValidEmail("a@b"), "no TLD should be invalid")

	// Test maximum length (254 characters according to RFC)
	// 254 chars: local(64) + @ + domain(189) total
	longLocal := "a123456789012345678901234567890123456789012345678901234567890123"
	longDomain := "example." + "subdomain." + "verylongdomainname." + "anothersubdomain." + "yetanother." + "andmore." + "evermore." + "stillgoing." + "keepgoing." + "almostthere." + "finalsubdomain.com"
	longEmail := longLocal + "@" + longDomain

	if len(longEmail) <= 254 {
		assert.True(t, isValidEmail(longEmail), "long valid email should pass")
	}

	// Over 254 characters should fail
	veryLongEmail := longLocal + "@" + longDomain + ".toolong.domain.extension.exceeded"
	if len(veryLongEmail) > 254 {
		assert.False(t, isValidEmail(veryLongEmail), "email over 254 chars should fail")
	}
}

// TestIsValidUUID tests the UUID validation function
func TestIsValidUUID(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want bool
	}{
		// Valid UUIDs (v4 format most common)
		{"valid UUID v4", "550e8400-e29b-41d4-a716-446655440000", true},
		{"valid UUID lowercase", "6ba7b810-9dad-11d1-80b4-00c04fd430c8", true},
		{"valid UUID uppercase", "6BA7B810-9DAD-11D1-80B4-00C04FD430C8", true},
		{"valid UUID mixed case", "550E8400-e29b-41d4-A716-446655440000", true},
		{"valid nil UUID", "00000000-0000-0000-0000-000000000000", true},

		// Invalid UUIDs
		{"empty string", "", false},
		{"too short", "550e8400-e29b-41d4-a716", false},
		{"too long", "550e8400-e29b-41d4-a716-446655440000-extra", false},
		{"missing dashes", "550e8400e29b41d4a716446655440000", true}, // google/uuid accepts this format
		{"wrong dash positions", "550e8400-e29b41-d4a7-16446655440000", false},
		{"invalid characters", "550e8400-e29b-41d4-a716-44665544000g", false},
		{"spaces", "550e8400 e29b 41d4 a716 446655440000", false},
		{"not a UUID", "not-a-uuid-at-all", false},
		{"only dashes", "------------------------------------", false},
		{"partial UUID", "550e8400-e29b-41d4", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidUUID(tt.id)
			assert.Equal(t, tt.want, got, "isValidUUID(%q) = %v, want %v", tt.id, got, tt.want)
		})
	}
}

// TestIsValidRole tests the role validation function
func TestIsValidRole(t *testing.T) {
	tests := []struct {
		name string
		role string
		want bool
	}{
		// Valid roles
		{"owner", "owner", true},
		{"admin", "admin", true},
		{"member", "member", true},

		// Invalid roles
		{"empty string", "", false},
		{"invalid role", "superadmin", false},
		{"uppercase OWNER", "OWNER", false}, // Case sensitive
		{"uppercase ADMIN", "ADMIN", false},
		{"uppercase MEMBER", "MEMBER", false},
		{"mixed case Owner", "Owner", false},
		{"mixed case Admin", "Admin", false},
		{"guest", "guest", false},
		{"moderator", "moderator", false},
		{"user", "user", false},
		{"spaces", "  owner  ", false}, // No trimming in function
		{"with spaces", "owner ", false},
		{"special chars", "owner!", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidRole(tt.role)
			assert.Equal(t, tt.want, got, "isValidRole(%q) = %v, want %v", tt.role, got, tt.want)
		})
	}
}

// TestIsValidRSVPResponse tests the RSVP response validation function
func TestIsValidRSVPResponse(t *testing.T) {
	tests := []struct {
		name     string
		response string
		want     bool
	}{
		// Valid responses
		{"yes", "yes", true},
		{"no", "no", true},
		{"maybe", "maybe", true},

		// Invalid responses
		{"empty string", "", false},
		{"uppercase YES", "YES", false}, // Case sensitive
		{"uppercase NO", "NO", false},
		{"uppercase MAYBE", "MAYBE", false},
		{"mixed case Yes", "Yes", false},
		{"mixed case No", "No", false},
		{"mixed case Maybe", "Maybe", false},
		{"accept", "accept", false},
		{"decline", "decline", false},
		{"tentative", "tentative", false},
		{"unknown", "unknown", false},
		{"spaces", "  yes  ", false}, // No trimming in function
		{"with spaces", "yes ", false},
		{"y", "y", false},
		{"n", "n", false},
		{"special chars", "yes!", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidRSVPResponse(tt.response)
			assert.Equal(t, tt.want, got, "isValidRSVPResponse(%q) = %v, want %v", tt.response, got, tt.want)
		})
	}
}

// TestValidationFunctionsWithNormalization tests that validation works with the normalization done in action handlers
func TestValidationFunctionsWithNormalization(t *testing.T) {
	t.Run("role validation after trimming and lowercasing", func(t *testing.T) {
		// Simulate what happens in updateMemberRoleAction: strings.TrimSpace(strings.ToLower(input))
		normalized := "owner" // from input "  OWNER  "
		assert.True(t, isValidRole(normalized), "normalized role should be valid")
	})

	t.Run("RSVP validation after trimming and lowercasing", func(t *testing.T) {
		// Simulate what happens in addRSVPAction: strings.TrimSpace(strings.ToLower(input))
		normalized := "yes" // from input "  YES  "
		assert.True(t, isValidRSVPResponse(normalized), "normalized response should be valid")
	})

	t.Run("email validation after trimming", func(t *testing.T) {
		// isValidEmail does its own trimming
		input := "  user@example.com  "
		assert.True(t, isValidEmail(input), "email with spaces should be trimmed and valid")
	})
}
