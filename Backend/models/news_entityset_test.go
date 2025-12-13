package models

import (
"testing"
)

// TestNewsEntitySetName verifies that the News entity has the EntitySetName method
// and that it returns "News" instead of the default pluralized "Newses"
func TestNewsEntitySetName(t *testing.T) {
news := News{}
entitySetName := news.EntitySetName()

if entitySetName != "News" {
t.Errorf("Expected EntitySetName to be 'News', got '%s'", entitySetName)
}

// Verify it doesn't return the incorrectly pluralized form
if entitySetName == "Newses" {
t.Error("EntitySetName should not return 'Newses' (incorrect pluralization)")
}

t.Logf("✅ News.EntitySetName() correctly returns '%s'", entitySetName)
}
