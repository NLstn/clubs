package models

import (
	"reflect"
	"testing"
)

// TestODataHookNamingConvention verifies that all OData hooks use the correct naming with 'OData' prefix
// This test ensures we don't accidentally add OData hooks without the required prefix
func TestODataHookNamingConvention(t *testing.T) {
	// List of all entities that are registered in OData
	// These should have OData hooks if they implement any lifecycle hooks
	entities := []interface{}{
		&User{},
		&UserSession{},
		&Club{},
		&Member{},
		&Team{},
		&TeamMember{},
		&Event{},
		&EventRSVP{},
		&Shift{},
		&ShiftMember{},
		&Fine{},
		&FineTemplate{},
		&Invite{},
		&JoinRequest{},
		&News{},
		&Notification{},
		&UserNotificationPreferences{},
		&ClubSettings{},
		&UserPrivacySettings{},
	}

	// OData hook signatures (these take context.Context and *http.Request)
	odataHookSignatures := []string{
		"BeforeCreate",
		"AfterCreate",
		"BeforeUpdate",
		"AfterUpdate",
		"BeforeDelete",
		"AfterDelete",
		"BeforeReadCollection",
		"BeforeReadEntity",
		"AfterReadCollection",
		"AfterReadEntity",
	}

	// Check each entity type
	for _, entity := range entities {
		entityType := reflect.TypeOf(entity)
		entityName := entityType.Elem().Name()

		// Check all methods on the entity
		for i := 0; i < entityType.NumMethod(); i++ {
			method := entityType.Method(i)
			methodName := method.Name

			// Check if this method name matches an OData hook name (without prefix)
			for _, hookName := range odataHookSignatures {
				if methodName == hookName {
					// This is a hook without OData prefix - check if it's a GORM hook or OData hook
					// GORM hooks take *gorm.DB as parameter, OData hooks take context.Context and *http.Request

					// Get the method signature
					methodType := method.Type
					if methodType.NumIn() > 1 { // First param is receiver
						firstParam := methodType.In(1)

						// Check if it's a GORM hook (takes *gorm.DB)
						isGormHook := firstParam.String() == "*gorm.DB"

						// Check if it might be an OData hook (takes context.Context)
						isODataHook := firstParam.String() == "context.Context" ||
							(methodType.NumIn() > 2 && methodType.In(2).String() == "*http.Request")

						if isODataHook && !isGormHook {
							t.Errorf("Entity %s has method %s without OData prefix, but it appears to be an OData hook (takes context.Context). Should be renamed to OData%s",
								entityName, methodName, methodName)
						}
					}
				}
			}
		}
	}
}

// TestGormHookNamingConvention verifies that GORM hooks (for UUID generation) don't have OData prefix
// and that they take *gorm.DB as parameter
func TestGormHookNamingConvention(t *testing.T) {
	// Entities that are known to have GORM hooks for UUID generation
	entitiesWithGormHooks := []interface{}{
		&Activity{},
		&Club{},
		&Notification{},
		&UserNotificationPreferences{},
		&Team{},
		&TeamMember{},
	}

	for _, entity := range entitiesWithGormHooks {
		entityType := reflect.TypeOf(entity)
		entityName := entityType.Elem().Name()

		// Look for BeforeCreate method
		method, found := entityType.MethodByName("BeforeCreate")
		if !found {
			t.Errorf("Entity %s should have BeforeCreate GORM hook for UUID generation", entityName)
			continue
		}

		// Verify it takes *gorm.DB as parameter
		methodType := method.Type
		if methodType.NumIn() < 2 {
			t.Errorf("Entity %s BeforeCreate should take a parameter", entityName)
			continue
		}

		firstParam := methodType.In(1)
		if firstParam.String() != "*gorm.DB" {
			t.Errorf("Entity %s BeforeCreate should take *gorm.DB as parameter, got %s", entityName, firstParam.String())
		}
	}
}

// TestODataHookSignatures verifies that all OData hooks have the correct signatures
func TestODataHookSignatures(t *testing.T) {
	// Sample entities with OData hooks
	entities := []interface{}{
		&Club{},
		&Member{},
		&Event{},
		&News{},
	}

	// Expected signature patterns for OData hooks
	expectedSignatures := map[string][]string{
		"ODataBeforeCreate": {"context.Context", "*http.Request"},
		"ODataAfterCreate":  {"context.Context", "*http.Request"},
		"ODataBeforeUpdate": {"context.Context", "*http.Request"},
		"ODataAfterUpdate":  {"context.Context", "*http.Request"},
		"ODataBeforeDelete": {"context.Context", "*http.Request"},
		"ODataAfterDelete":  {"context.Context", "*http.Request"},
		// Read hooks have an additional parameter for query options
		"ODataBeforeReadCollection": {"context.Context", "*http.Request", "interface {}"},
		"ODataBeforeReadEntity":     {"context.Context", "*http.Request", "interface {}"},
		"ODataAfterReadCollection":  {"context.Context", "*http.Request", "interface {}", "interface {}"},
		"ODataAfterReadEntity":      {"context.Context", "*http.Request", "interface {}", "interface {}"},
	}

	for _, entity := range entities {
		entityType := reflect.TypeOf(entity)
		entityName := entityType.Elem().Name()

		for hookName, expectedParams := range expectedSignatures {
			method, found := entityType.MethodByName(hookName)
			if !found {
				// Not all entities implement all hooks, so this is fine
				continue
			}

			methodType := method.Type
			// First param is receiver, so actual params start at index 1
			actualParamCount := methodType.NumIn() - 1

			if actualParamCount != len(expectedParams) {
				t.Errorf("Entity %s method %s has %d parameters, expected %d",
					entityName, hookName, actualParamCount, len(expectedParams))
				continue
			}

			// Verify parameter types
			for i, expectedParam := range expectedParams {
				actualParam := methodType.In(i + 1).String()

				if actualParam != expectedParam {
					t.Errorf("Entity %s method %s parameter %d: expected %s, got %s",
						entityName, hookName, i+1, expectedParam, actualParam)
				}
			}
		}
	}
}
