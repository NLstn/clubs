package frontend

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeLinks(t *testing.T) {
	originalURL := os.Getenv("FRONTEND_URL")
	os.Setenv("FRONTEND_URL", "http://example.com")
	defer os.Setenv("FRONTEND_URL", originalURL)

	err := Init()
	assert.NoError(t, err)

	t.Run("MakeMagicLink", func(t *testing.T) {
		link := MakeMagicLink("token123")
		assert.Equal(t, "http://example.com/auth/magic?token=token123", link)
	})

	t.Run("MakeClubLink", func(t *testing.T) {
		link := MakeClubLink("club1")
		assert.Equal(t, "http://example.com/clubs/club1", link)
	})

	t.Run("MakeEventLink", func(t *testing.T) {
		link := MakeEventLink("club1", "event2")
		assert.Equal(t, "http://example.com/clubs/club1/events/event2", link)
	})

	t.Run("MakeFineLink", func(t *testing.T) {
		link := MakeFineLink("club1", "fine3")
		assert.Equal(t, "http://example.com/clubs/club1/fines/fine3", link)
	})
}
