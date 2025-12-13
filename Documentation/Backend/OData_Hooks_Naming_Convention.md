# OData Hooks Naming Convention

## Overview

The Clubs backend uses two different types of hooks with distinct naming conventions:

1. **GORM Hooks** - Database-level operations (no prefix)
2. **OData Hooks** - API-level operations (with `OData` prefix)

## GORM Hooks (No Prefix)

### Purpose
GORM hooks are used for database-level operations such as:
- Generating UUIDs for new records
- Setting default values
- Database-level validations

### Naming Convention
GORM hooks use the standard GORM naming **WITHOUT** any prefix:
- `BeforeCreate(tx *gorm.DB) error`
- `AfterCreate(tx *gorm.DB) error`
- `BeforeUpdate(tx *gorm.DB) error`
- etc.

### Signature
```go
func (e *Entity) BeforeCreate(tx *gorm.DB) error {
    // Database-level operations
    return nil
}
```

### Example
```go
// Activity entity - GORM hook for UUID generation
func (a *Activity) BeforeCreate(tx *gorm.DB) error {
    if a.ID == "" {
        a.ID = uuid.New().String()
    }
    return nil
}
```

### Entities Using GORM Hooks
The following entities have GORM hooks for UUID generation:
- `Activity`
- `Club`
- `Notification`
- `UserNotificationPreferences`
- `Team`
- `TeamMember`

## OData Hooks (With OData Prefix)

### Purpose
OData hooks are used for API-level operations such as:
- Authorization checks (verifying user permissions)
- Setting audit fields from HTTP context (CreatedBy, UpdatedBy)
- Filtering data based on user context
- Validating business rules that require HTTP context

### Naming Convention
OData hooks **MUST** use the `OData` prefix:
- `ODataBeforeCreate(ctx context.Context, r *http.Request) error`
- `ODataAfterCreate(ctx context.Context, r *http.Request) error`
- `ODataBeforeUpdate(ctx context.Context, r *http.Request) error`
- `ODataBeforeDelete(ctx context.Context, r *http.Request) error`
- `ODataBeforeReadCollection(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error)`
- `ODataBeforeReadEntity(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error)`
- `ODataAfterReadCollection(ctx context.Context, r *http.Request, opts interface{}, results interface{}) (interface{}, error)`
- `ODataAfterReadEntity(ctx context.Context, r *http.Request, opts interface{}, entity interface{}) (interface{}, error)`

### Signature Patterns

#### Write Hooks (Create/Update/Delete)
```go
func (e *Entity) ODataBeforeCreate(ctx context.Context, r *http.Request) error {
    // Authorization and validation
    userID, ok := ctx.Value(auth.UserIDKey).(string)
    if !ok || userID == "" {
        return fmt.Errorf("unauthorized")
    }
    
    // Set audit fields
    e.CreatedBy = userID
    e.UpdatedBy = userID
    
    return nil
}
```

#### Read Hooks (Collection/Entity)
```go
func (e Entity) ODataBeforeReadCollection(ctx context.Context, r *http.Request, opts interface{}) ([]func(*gorm.DB) *gorm.DB, error) {
    userID, ok := ctx.Value(auth.UserIDKey).(string)
    if !ok || userID == "" {
        return nil, fmt.Errorf("unauthorized")
    }
    
    // Return GORM scopes to filter results
    scope := func(db *gorm.DB) *gorm.DB {
        return db.Where("user_id = ?", userID)
    }
    
    return []func(*gorm.DB) *gorm.DB{scope}, nil
}
```

### Example
```go
// Club entity - OData hook for authorization and audit fields
func (c *Club) ODataBeforeCreate(ctx context.Context, r *http.Request) error {
    // Extract user ID from context
    userID, ok := ctx.Value(auth.UserIDKey).(string)
    if !ok || userID == "" {
        return fmt.Errorf("unauthorized: user ID not found in context")
    }

    // Set audit fields from HTTP context
    now := time.Now()
    c.CreatedAt = now
    c.CreatedBy = userID
    c.UpdatedAt = now
    c.UpdatedBy = userID

    return nil
}
```

## Why Two Different Naming Conventions?

### Reason 1: Clear Separation of Concerns
- **GORM hooks** handle database-level operations that are always executed regardless of how the entity is created
- **OData hooks** handle API-level operations that only apply to OData API requests

### Reason 2: Execution Order
The two types of hooks run at different times:
1. **OData hooks** run first (when processing HTTP request)
2. **GORM hooks** run second (when executing database operation)

This is noted in the Club entity:
```go
// BeforeCreate GORM hook - sets UUID if not provided
func (c *Club) BeforeCreate(tx *gorm.DB) error {
    if c.ID == "" {
        c.ID = uuid.New().String()
    }
    
    // Note: CreatedBy and UpdatedBy are set by OData hooks from HTTP context
    // GORM BeforeCreate runs after OData BeforeCreate
    
    return nil
}
```

### Reason 3: Different Parameters
- **GORM hooks** receive `*gorm.DB` to access the database transaction
- **OData hooks** receive `context.Context` and `*http.Request` to access HTTP context and user information

## Enforcement

### Automated Testing
The `hooks_naming_test.go` file contains tests to verify:
1. All OData hooks use the `OData` prefix
2. All GORM hooks do NOT use the `OData` prefix
3. Hook signatures are correct

### Migration from Old Naming
If you find any hooks that take `context.Context` and `*http.Request` but don't have the `OData` prefix, they should be renamed:

```go
// OLD (incorrect) ❌
func (e *Entity) BeforeCreate(ctx context.Context, r *http.Request) error {
    // ...
}

// NEW (correct) ✅
func (e *Entity) ODataBeforeCreate(ctx context.Context, r *http.Request) error {
    // ...
}
```

## Complete Example: Entity with Both Hook Types

```go
type MyEntity struct {
    ID        string    `json:"ID" gorm:"primaryKey;type:char(36)"`
    Name      string    `json:"Name"`
    CreatedBy string    `json:"CreatedBy" gorm:"type:char(36)"`
    CreatedAt time.Time `json:"CreatedAt"`
}

// GORM hook - for UUID generation (no OData prefix)
func (e *MyEntity) BeforeCreate(tx *gorm.DB) error {
    if e.ID == "" {
        e.ID = uuid.New().String()
    }
    return nil
}

// OData hook - for authorization and audit fields (with OData prefix)
func (e *MyEntity) ODataBeforeCreate(ctx context.Context, r *http.Request) error {
    userID, ok := ctx.Value(auth.UserIDKey).(string)
    if !ok || userID == "" {
        return fmt.Errorf("unauthorized")
    }
    
    e.CreatedBy = userID
    e.CreatedAt = time.Now()
    
    return nil
}
```

## Summary

| Hook Type | Prefix | Parameters | Purpose | Example |
|-----------|--------|------------|---------|---------|
| GORM | None | `tx *gorm.DB` | Database operations | UUID generation |
| OData | `OData` | `ctx context.Context, r *http.Request` | API operations | Authorization, audit fields |

**Remember:** When implementing lifecycle hooks that need access to HTTP context or user information, always use the `OData` prefix!
