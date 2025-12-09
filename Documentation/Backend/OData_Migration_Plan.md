# OData API Migration Plan

## Executive Summary

This document outlines the comprehensive plan to migrate the Clubs backend from custom REST APIs to OData v4 APIs using the `github.com/NLstn/go-odata` package. The migration will maintain **100% feature parity** while gaining the benefits of standardized querying, filtering, expansion, and pagination capabilities.

## Current Architecture Analysis

### Existing API Structure

The current backend has 15 main API domains:

1. **Authentication** - Magic link & Keycloak OAuth
2. **Clubs** - Club CRUD and logo management
3. **Members** - Member management and roles
4. **Teams** - Team organization within clubs
5. **Events** - Event scheduling with RSVP
6. **Recurring Events** - Complex recurring event patterns
7. **Shifts** - Shift scheduling and assignments
8. **Fines** - Fine tracking and management
9. **Fine Templates** - Predefined fine configurations
10. **Invites** - User invitation system
11. **Join Requests** - Club membership requests
12. **News** - Club announcements
13. **Notifications** - User notification system
14. **User Profile** - User account management
15. **Dashboard** - Aggregated data views
16. **Search** - Global search functionality
17. **Club Settings** - Feature toggles per club
18. **Privacy** - User privacy preferences

### Core Data Models

#### Primary Entities
- **Club** - Club with soft delete, logo, creator tracking
- **User** - User accounts with email/name
- **Member** - Club membership with roles (owner/admin/member)
- **Team** - Sub-groups within clubs with soft delete
- **Event** - Events with date/time, RSVP, recurring patterns
- **EventRSVP** - RSVP responses (yes/no/maybe)
- **Shift** - Work shifts for events
- **ShiftMember** - Shift assignments
- **Fine** - Financial penalties with amount/reason
- **FineTemplate** - Reusable fine definitions
- **Invite** - Email-based club invitations
- **JoinRequest** - User requests to join clubs
- **News** - Club news posts
- **Notification** - User notifications with types
- **ClubSettings** - Per-club feature flags
- **UserPrivacySettings** - Privacy preferences per club/global
- **MagicLink** - Authentication tokens

## OData Migration Strategy

### Phase 1: Foundation & Entity Registration

#### 1.1 OData Service Setup

**Goal:** Initialize go-odata service alongside existing REST API

**Tasks:**
- Install `github.com/NLstn/go-odata` dependency
- Create new `odata` package structure
- Initialize OData service with GORM connection
- Mount OData service at `/api/v2/` path
- Keep existing `/api/v1/` endpoints running

**Implementation:**
```go
// Backend/odata/service.go
package odata

import (
    "github.com/nlstn/go-odata"
    "gorm.io/gorm"
)

func NewODataService(db *gorm.DB) (*odata.Service, error) {
    service, err := odata.NewServiceWithConfig(db, odata.ServiceConfig{
        PersistentChangeTracking: false, // Enable later if needed
    })
    if err != nil {
        return nil, err
    }
    
    // Set namespace
    if err := service.SetNamespace("ClubsService"); err != nil {
        return nil, err
    }
    
    return service, nil
}
```

#### 1.2 Model Annotation & Mapping

**Goal:** Annotate existing GORM models with OData tags

**Strategy:**
- Add `odata:` tags to all model fields
- Mark primary keys with `odata:"key"`
- Mark required fields with `odata:"required"`
- Mark nullable fields with `odata:"nullable"`
- Define navigation properties with `odata:"nav"`
- Add computed properties where needed

**Example:**
```go
type Club struct {
    ID          string     `json:"id" gorm:"type:uuid;primary_key" odata:"key"`
    Name        string     `json:"name" odata:"required"`
    Description string     `json:"description" odata:"nullable"`
    LogoURL     *string    `json:"logo_url,omitempty" odata:"nullable"`
    CreatedAt   time.Time  `json:"created_at" odata:"immutable"`
    CreatedBy   string     `json:"created_by" gorm:"type:uuid" odata:"required"`
    UpdatedAt   time.Time  `json:"updated_at"`
    UpdatedBy   string     `json:"updated_by" gorm:"type:uuid" odata:"required"`
    Deleted     bool       `json:"deleted" gorm:"default:false" odata:"required"`
    DeletedAt   *time.Time `json:"deleted_at,omitempty" odata:"nullable"`
    DeletedBy   *string    `json:"deleted_by,omitempty" gorm:"type:uuid" odata:"nullable"`
    
    // Navigation properties
    Members     []Member   `gorm:"foreignKey:ClubID" odata:"nav"`
    Events      []Event    `gorm:"foreignKey:ClubID" odata:"nav"`
    Teams       []Team     `gorm:"foreignKey:ClubID" odata:"nav"`
    News        []News     `gorm:"foreignKey:ClubID" odata:"nav"`
    Fines       []Fine     `gorm:"foreignKey:ClubID" odata:"nav"`
}
```

#### 1.3 Entity Registration

**Goal:** Register all entities with the OData service

**Entities to Register:**
1. Users (base entity)
2. Clubs (with nav to Members, Events, Teams, News, Fines)
3. Members (with nav to User, Club, Teams)
4. Teams (with nav to Club, Members, Events, Fines)
5. Events (with nav to Club, Team, RSVPs, Shifts)
6. EventRSVPs (with nav to Event, User)
7. Shifts (with nav to Event, ShiftMembers)
8. ShiftMembers (with nav to Shift, User)
9. Fines (with nav to Club, User, Team)
10. FineTemplates (with nav to Club)
11. Invites (with nav to Club, User)
12. JoinRequests (with nav to Club, User)
13. News (with nav to Club)
14. Notifications (with nav to User, Club, Event, Fine, Invite, JoinRequest)
15. ClubSettings (with nav to Club)
16. UserPrivacySettings (with nav to User, Club)

### Phase 2: Authentication & Authorization

#### 2.1 OData Authentication Middleware

**Goal:** Integrate JWT authentication with OData service

**Strategy:**
- Create OData-compatible auth middleware
- Extract user from JWT token
- Inject user context into OData request pipeline
- Leverage go-odata's lifecycle hooks for authorization

**Implementation:**
```go
// Backend/odata/middleware/auth.go
func ODataAuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract JWT token from Authorization header
        token := extractToken(r)
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        // Validate token and extract user ID
        userID, err := validateJWT(token)
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        // Add user to context
        ctx := context.WithValue(r.Context(), "userID", userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

#### 2.2 Read Hooks for Authorization

**Goal:** Implement row-level security using OData read hooks

**Strategy:**
- Use `RegisterReadHook` to filter queries based on user permissions
- Implement club membership checks
- Filter soft-deleted items for non-owners
- Apply privacy settings

**Example:**
```go
// Register read hook for Clubs entity
service.RegisterReadHook("Clubs", func(ctx context.Context, query *gorm.DB) (*gorm.DB, error) {
    userID := ctx.Value("userID").(string)
    
    // Only return clubs where user is a member
    return query.
        Joins("INNER JOIN members ON members.club_id = clubs.id").
        Where("members.user_id = ?", userID).
        Where("clubs.deleted = false OR clubs.created_by = ?", userID), nil
})

// Register read hook for Members entity
service.RegisterReadHook("Members", func(ctx context.Context, query *gorm.DB) (*gorm.DB, error) {
    userID := ctx.Value("userID").(string)
    
    // Check club membership and settings
    return query.
        Joins("INNER JOIN members m2 ON m2.club_id = members.club_id").
        Where("m2.user_id = ?", userID), nil
})
```

#### 2.3 Lifecycle Hooks for Write Authorization

**Goal:** Implement write permission checks using BeforeCreate/BeforeUpdate hooks

**Example:**
```go
// Register before create hook for Events
service.RegisterBeforeCreate("Events", func(ctx context.Context, entity interface{}) error {
    userID := ctx.Value("userID").(string)
    event := entity.(*Event)
    
    // Check if user is admin/owner of the club
    isAdmin, err := checkAdminRights(event.ClubID, userID)
    if err != nil {
        return err
    }
    if !isAdmin {
        return fmt.Errorf("only admins can create events")
    }
    
    event.CreatedBy = userID
    event.UpdatedBy = userID
    return nil
})
```

### Phase 3: Core CRUD Operations

#### 3.1 Entity Collections (GET /api/v2/EntitySet)

**Automatic OData Features:**
- `$filter` - Filter results: `?$filter=name eq 'Soccer Club'`
- `$select` - Choose fields: `?$select=id,name,description`
- `$expand` - Include related: `?$expand=Members,Events`
- `$orderby` - Sort: `?$orderby=name asc`
- `$top` / `$skip` - Pagination: `?$top=10&$skip=20`
- `$count` - Total count: `?$count=true`
- `$search` - Full-text: `?$search=soccer`

**REST Equivalent Mapping:**
- `GET /api/v1/clubs` → `GET /api/v2/Clubs`
- `GET /api/v1/clubs/{clubid}/members` → `GET /api/v2/Members?$filter=clubId eq '{clubid}'`
- `GET /api/v1/clubs/{clubid}/events` → `GET /api/v2/Events?$filter=clubId eq '{clubid}'`

#### 3.2 Single Entity (GET /api/v2/EntitySet(key))

**Examples:**
- `GET /api/v1/clubs/{clubid}` → `GET /api/v2/Clubs('{clubid}')`
- `GET /api/v1/clubs/{clubid}/events/{eventid}` → `GET /api/v2/Events('{eventid}')`

#### 3.3 Create (POST /api/v2/EntitySet)

**Example:**
```http
POST /api/v2/Clubs
Content-Type: application/json

{
  "name": "Soccer Club",
  "description": "A club for soccer enthusiasts"
}
```

#### 3.4 Update (PATCH /api/v2/EntitySet(key))

**Example:**
```http
PATCH /api/v2/Clubs('{clubid}')
Content-Type: application/json

{
  "name": "Updated Soccer Club",
  "description": "New description"
}
```

#### 3.5 Delete (DELETE /api/v2/EntitySet(key))

**Strategy:**
- Implement soft delete using BeforeDelete hook
- Set `Deleted = true` and `DeletedAt = now`
- Actual DELETE performs soft delete for applicable entities

### Phase 4: Custom Operations (Actions & Functions)

#### 4.1 OData Actions (State-Changing Operations)

**Actions to Implement:**

1. **AcceptInvite** - Bound to Invite entity
```go
service.RegisterAction("Invite", "Accept", func(ctx context.Context, entityKey interface{}, params map[string]interface{}) (interface{}, error) {
    userID := ctx.Value("userID").(string)
    inviteID := entityKey.(string)
    
    // Business logic for accepting invite
    return acceptInvite(inviteID, userID)
})
```

2. **RejectInvite** - Bound to Invite entity
3. **AcceptJoinRequest** - Bound to JoinRequest entity
4. **RejectJoinRequest** - Bound to JoinRequest entity
5. **LeaveClub** - Bound to Club entity
6. **UploadLogo** - Bound to Club entity (special handling for multipart)
7. **DeleteLogo** - Bound to Club entity
8. **MarkNotificationRead** - Bound to Notification entity
9. **MarkAllNotificationsRead** - Unbound action

**Example Registration:**
```go
// Bound action: LeaveClub
service.RegisterAction("Club", "Leave", func(ctx context.Context, entityKey interface{}, params map[string]interface{}) (interface{}, error) {
    userID := ctx.Value("userID").(string)
    clubID := entityKey.(string)
    
    return leaveClub(clubID, userID)
})
```

#### 4.2 OData Functions (Non-State-Changing Queries)

**Functions to Implement:**

1. **CheckAdminRights** - Bound to Club entity
```go
service.RegisterFunction("Club", "IsAdmin", func(ctx context.Context, entityKey interface{}, params map[string]interface{}) (interface{}, error) {
    userID := ctx.Value("userID").(string)
    clubID := entityKey.(string)
    
    isAdmin, err := checkAdminRights(clubID, userID)
    return map[string]bool{"isAdmin": isAdmin}, err
})
```

2. **GetOwnerCount** - Bound to Club entity
3. **GetInviteLink** - Bound to Club entity
4. **GetUpcomingEvents** - Bound to Club entity
5. **GetDashboardNews** - Unbound function
6. **GetDashboardEvents** - Unbound function
7. **GetDashboardActivities** - Unbound function
8. **SearchGlobal** - Unbound function

**Example:**
```go
// Unbound function: GetDashboardNews
service.RegisterFunction("", "GetDashboardNews", func(ctx context.Context, entityKey interface{}, params map[string]interface{}) (interface{}, error) {
    userID := ctx.Value("userID").(string)
    
    return getDashboardNews(userID)
})
```

### Phase 5: Complex Scenarios

#### 5.1 File Upload (Club Logo)

**Strategy:**
- OData doesn't natively support multipart/form-data
- Implement custom action with manual handler
- Use Azure Blob Storage as before

**Implementation:**
```go
// Custom handler for logo upload
mux.Handle("/api/v2/Clubs({clubId})/UploadLogo", func(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // Extract club ID from path
    clubID := extractClubID(r.URL.Path)
    
    // Handle multipart form
    file, header, err := r.FormFile("logo")
    if err != nil {
        http.Error(w, "No file provided", http.StatusBadRequest)
        return
    }
    defer file.Close()
    
    // Upload to Azure Blob Storage
    logoURL, err := azure.UploadClubLogo(clubID, file, header)
    if err != nil {
        http.Error(w, "Upload failed", http.StatusInternalServerError)
        return
    }
    
    // Update club in database
    // Return OData-compatible response
})
```

#### 5.2 Magic Link Authentication

**Strategy:**
- Keep authentication endpoints as custom REST API
- Authentication is not resource-oriented, so REST is appropriate
- Mount at `/api/v1/auth/` alongside OData service

**Endpoints to Keep:**
- `POST /api/v1/auth/requestMagicLink`
- `GET /api/v1/auth/verifyMagicLink`
- `POST /api/v1/auth/refreshToken`
- `POST /api/v1/auth/logout`
- `POST /api/v1/auth/keycloak/*` (all Keycloak endpoints)

#### 5.3 Recurring Events

**Strategy:**
- Store recurring event patterns in Event entity
- Use OData function to generate instances
- Clients can query base pattern or expanded instances

**Implementation:**
```go
// Function to expand recurring events
service.RegisterFunction("Event", "ExpandRecurrence", func(ctx context.Context, entityKey interface{}, params map[string]interface{}) (interface{}, error) {
    eventID := entityKey.(string)
    startDate := params["startDate"].(time.Time)
    endDate := params["endDate"].(time.Time)
    
    // Generate event instances based on recurrence pattern
    return generateRecurringInstances(eventID, startDate, endDate)
})
```

#### 5.4 Soft Delete Visibility

**Strategy:**
- Use read hooks to filter deleted items
- Owners can see deleted items in their clubs
- Use `$filter` parameter to explicitly request deleted items

**Implementation:**
```go
service.RegisterReadHook("Clubs", func(ctx context.Context, query *gorm.DB) (*gorm.DB, error) {
    userID := ctx.Value("userID").(string)
    
    // Check if user explicitly wants to see deleted items
    if ctx.Value("includeDeleted") == true {
        // Only owners can see deleted clubs
        return query.Where("deleted = false OR created_by = ?", userID), nil
    }
    
    // Default: hide deleted items
    return query.Where("deleted = false"), nil
})
```

#### 5.5 Privacy Settings & Member Visibility

**Strategy:**
- Use read hooks to apply privacy settings
- Check club settings for member list visibility
- Filter based on user role and preferences

### Phase 6: Advanced OData Features

#### 6.1 Navigation Property Expansion

**Examples:**
```http
# Get club with all members
GET /api/v2/Clubs('{clubId}')?$expand=Members

# Get club with members and their user details
GET /api/v2/Clubs('{clubId}')?$expand=Members($expand=User)

# Get events with RSVPs and user details
GET /api/v2/Events?$expand=RSVPs($expand=User)&$filter=clubId eq '{clubId}'

# Complex expansion
GET /api/v2/Clubs('{clubId}')?$expand=Events($expand=RSVPs),Teams($expand=Members)
```

#### 6.2 Advanced Filtering

**Examples:**
```http
# Filter events by date range
GET /api/v2/Events?$filter=startTime ge 2024-01-01T00:00:00Z and startTime le 2024-12-31T23:59:59Z

# Filter members by role
GET /api/v2/Members?$filter=clubId eq '{clubId}' and role eq 'admin'

# Filter unpaid fines for a user
GET /api/v2/Fines?$filter=userId eq '{userId}' and paid eq false

# Complex filter with navigation properties
GET /api/v2/Events?$filter=Club/name eq 'Soccer Club' and startTime gt now()
```

#### 6.3 Aggregation & Grouping

**Examples:**
```http
# Count members per club
GET /api/v2/Members?$apply=groupby((clubId),aggregate($count as memberCount))

# Sum fines by club
GET /api/v2/Fines?$apply=groupby((clubId),aggregate(amount with sum as totalFines))

# Average event attendance
GET /api/v2/EventRSVPs?$apply=filter(response eq 'yes')/groupby((eventId),aggregate($count as attendees))
```

#### 6.4 Full-Text Search

**Strategy:**
- Leverage go-odata's built-in `$search` parameter
- Configure searchable fields per entity

**Example:**
```go
// During entity registration, specify searchable fields
// This is typically done via model tags or service configuration
```

#### 6.5 Change Tracking (Delta Queries)

**Strategy:**
- Enable change tracking for high-traffic entities
- Clients can request delta updates since last sync

**Implementation:**
```go
// Enable change tracking for specific entities
service.EnableChangeTracking("Events")
service.EnableChangeTracking("News")
service.EnableChangeTracking("Notifications")
```

**Client Usage:**
```http
# Initial request
GET /api/v2/Events?$deltatoken

# Response includes deltatoken
# Later requests
GET /api/v2/Events?$deltatoken=abc123

# Returns only changes since last sync
```

### Phase 7: Testing Strategy

#### 7.1 Unit Tests

**Test Coverage:**
- Entity registration validation
- OData tag correctness
- Navigation property setup
- Authentication middleware
- Authorization hooks
- Custom actions/functions

**Example:**
```go
func TestClubEntityRegistration(t *testing.T) {
    service := setupTestODataService(t)
    
    // Verify Club entity is registered
    metadata := service.GetMetadata()
    clubEntity := metadata.Entities["Clubs"]
    assert.NotNil(t, clubEntity)
    
    // Verify key property
    assert.Equal(t, "ID", clubEntity.KeyProperty)
    
    // Verify navigation properties
    assert.Contains(t, clubEntity.NavigationProperties, "Members")
    assert.Contains(t, clubEntity.NavigationProperties, "Events")
}
```

#### 7.2 Integration Tests

**Test Scenarios:**
- CRUD operations for each entity
- Complex queries with $filter, $expand, $select
- Custom actions and functions
- Authentication flow
- Authorization enforcement
- Soft delete behavior
- Privacy settings

**Example:**
```go
func TestGetClubWithMembers(t *testing.T) {
    service, db := setupTestODataService(t)
    token := createTestUser(t, db)
    
    // Create test data
    club := createTestClub(t, db)
    createTestMembers(t, db, club.ID, 5)
    
    // Execute OData query
    req := httptest.NewRequest("GET", "/api/v2/Clubs('"+club.ID+"')?$expand=Members", nil)
    req.Header.Set("Authorization", "Bearer "+token)
    
    w := httptest.NewRecorder()
    service.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
    
    // Verify response structure
    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    
    assert.Equal(t, club.ID, response["id"])
    assert.NotNil(t, response["Members"])
    assert.Len(t, response["Members"], 5)
}
```

#### 7.3 OData Compliance Tests

**Strategy:**
- Use go-odata's compliance test suite
- Ensure standard OData v4 compatibility
- Verify metadata document correctness

```bash
cd Backend/odata
go test -v ./compliance/...
```

#### 7.4 Performance Tests

**Test Scenarios:**
- Large result sets with pagination
- Complex expand queries
- Aggregation queries
- Concurrent request handling

### Phase 8: Migration Execution Plan

#### 8.1 Parallel Development (Weeks 1-2)

**Tasks:**
1. Set up OData service infrastructure
2. Annotate all models with OData tags
3. Register all entities
4. Implement authentication middleware
5. Test basic CRUD operations

**Success Criteria:**
- OData service running at `/api/v2/`
- All entities registered and queryable
- Authentication working
- Basic CRUD operations functional

#### 8.2 Authorization & Hooks (Weeks 3-4)

**Tasks:**
1. Implement read hooks for all entities
2. Implement write hooks (BeforeCreate, BeforeUpdate, BeforeDelete)
3. Test authorization rules
4. Implement soft delete logic

**Success Criteria:**
- All authorization rules enforced
- Users can only access authorized data
- Soft delete working correctly
- Privacy settings applied

#### 8.3 Custom Operations (Weeks 5-6)

**Tasks:**
1. Implement all custom actions
2. Implement all custom functions
3. Handle special cases (file upload, magic link)
4. Test all operations

**Success Criteria:**
- All actions callable and working
- All functions returning correct data
- File upload functioning
- Authentication endpoints operational

#### 8.4 Advanced Features (Week 7)

**Tasks:**
1. Enable change tracking
2. Configure full-text search
3. Test complex expansions
4. Optimize query performance

**Success Criteria:**
- Delta queries working
- Search returning relevant results
- Complex expansions performant
- Query response times acceptable

#### 8.5 Testing & Validation (Week 8)

**Tasks:**
1. Comprehensive integration testing
2. OData compliance testing
3. Performance testing
4. Security audit

**Success Criteria:**
- All tests passing
- OData compliance verified
- Performance benchmarks met
- Security review approved

#### 8.6 Frontend Migration (Weeks 9-12)

**Tasks:**
1. Update frontend API client to use OData endpoints
2. Leverage OData features ($expand, $filter, etc.)
3. Update UI components for OData response format
4. Test all frontend workflows

**Success Criteria:**
- All frontend features working with OData
- User experience maintained or improved
- No breaking changes for users

#### 8.7 Deprecation & Cleanup (Week 13)

**Tasks:**
1. Add deprecation warnings to old REST endpoints
2. Monitor usage of old endpoints
3. Communicate migration to API consumers
4. Plan sunset date for old API

**Success Criteria:**
- Migration timeline communicated
- Old endpoints marked deprecated
- Monitoring in place

### Phase 9: API Endpoint Mapping

#### Complete REST to OData Mapping

| REST Endpoint | OData Equivalent | Notes |
|--------------|------------------|-------|
| **Authentication** | | |
| `POST /api/v1/auth/requestMagicLink` | `POST /api/v1/auth/requestMagicLink` | Keep as REST |
| `GET /api/v1/auth/verifyMagicLink` | `GET /api/v1/auth/verifyMagicLink` | Keep as REST |
| `POST /api/v1/auth/refreshToken` | `POST /api/v1/auth/refreshToken` | Keep as REST |
| `POST /api/v1/auth/logout` | `POST /api/v1/auth/logout` | Keep as REST |
| `POST /api/v1/auth/keycloak/login` | `POST /api/v1/auth/keycloak/login` | Keep as REST |
| `GET /api/v1/auth/keycloak/callback` | `GET /api/v1/auth/keycloak/callback` | Keep as REST |
| `POST /api/v1/auth/keycloak/logout` | `POST /api/v1/auth/keycloak/logout` | Keep as REST |
| **Clubs** | | |
| `GET /api/v1/clubs` | `GET /api/v2/Clubs` | Auto-filtered by membership |
| `GET /api/v1/clubs/{id}` | `GET /api/v2/Clubs('{id}')` | |
| `POST /api/v1/clubs` | `POST /api/v2/Clubs` | |
| `PATCH /api/v1/clubs/{id}` | `PATCH /api/v2/Clubs('{id}')` | |
| `DELETE /api/v1/clubs/{id}` | `DELETE /api/v2/Clubs('{id}')` | Soft delete via hook |
| `POST /api/v1/clubs/{id}/logo` | `POST /api/v2/Clubs('{id}')/UploadLogo` | Custom action |
| `DELETE /api/v1/clubs/{id}/logo` | `POST /api/v2/Clubs('{id}')/DeleteLogo` | Custom action |
| `DELETE /api/v1/clubs/{id}/hard-delete` | `POST /api/v2/Clubs('{id}')/HardDelete` | Custom action (admin only) |
| **Members** | | |
| `GET /api/v1/clubs/{clubid}/members` | `GET /api/v2/Members?$filter=clubId eq '{clubid}'&$expand=User` | |
| `PATCH /api/v1/clubs/{clubid}/members/{id}` | `PATCH /api/v2/Members('{id}')` | Update role |
| `DELETE /api/v1/clubs/{clubid}/members/{id}` | `DELETE /api/v2/Members('{id}')` | |
| `POST /api/v1/clubs/{clubid}/leave` | `POST /api/v2/Clubs('{clubid}')/Leave` | Custom action |
| `GET /api/v1/clubs/{clubid}/isAdmin` | `GET /api/v2/Clubs('{clubid}')/IsAdmin()` | Custom function |
| `GET /api/v1/clubs/{clubid}/ownerCount` | `GET /api/v2/Clubs('{clubid}')/GetOwnerCount()` | Custom function |
| **Teams** | | |
| `GET /api/v1/clubs/{clubid}/teams` | `GET /api/v2/Teams?$filter=clubId eq '{clubid}'` | |
| `POST /api/v1/clubs/{clubid}/teams` | `POST /api/v2/Teams` | |
| `GET /api/v1/clubs/{clubid}/teams/{id}` | `GET /api/v2/Teams('{id}')` | |
| `PATCH /api/v1/clubs/{clubid}/teams/{id}` | `PATCH /api/v2/Teams('{id}')` | |
| `DELETE /api/v1/clubs/{clubid}/teams/{id}` | `DELETE /api/v2/Teams('{id}')` | Soft delete |
| `GET /api/v1/clubs/{clubid}/teams/{id}/members` | `GET /api/v2/TeamMembers?$filter=teamId eq '{id}'&$expand=User` | |
| `POST /api/v1/clubs/{clubid}/teams/{id}/members` | `POST /api/v2/TeamMembers` | |
| `DELETE /api/v1/clubs/{clubid}/teams/{teamid}/members/{id}` | `DELETE /api/v2/TeamMembers('{id}')` | |
| **Events** | | |
| `GET /api/v1/clubs/{clubid}/events` | `GET /api/v2/Events?$filter=clubId eq '{clubid}'` | |
| `POST /api/v1/clubs/{clubid}/events` | `POST /api/v2/Events` | |
| `POST /api/v1/clubs/{clubid}/events/recurring` | `POST /api/v2/Events` | Handle via isRecurring field |
| `GET /api/v1/clubs/{clubid}/events/{id}` | `GET /api/v2/Events('{id}')` | |
| `PUT /api/v1/clubs/{clubid}/events/{id}` | `PATCH /api/v2/Events('{id}')` | OData prefers PATCH |
| `DELETE /api/v1/clubs/{clubid}/events/{id}` | `DELETE /api/v2/Events('{id}')` | |
| `GET /api/v1/clubs/{clubid}/events/upcoming` | `GET /api/v2/Clubs('{clubid}')/GetUpcomingEvents()` | Custom function |
| `POST /api/v1/clubs/{clubid}/events/{id}/rsvp` | `POST /api/v2/EventRSVPs` | Create/update |
| `GET /api/v1/clubs/{clubid}/events/{id}/rsvps` | `GET /api/v2/EventRSVPs?$filter=eventId eq '{id}'&$expand=User` | |
| **Shifts** | | |
| `GET /api/v1/clubs/{clubid}/shifts` | `GET /api/v2/Shifts?$filter=clubId eq '{clubid}'` | |
| `POST /api/v1/clubs/{clubid}/shifts` | `POST /api/v2/Shifts` | |
| `GET /api/v1/clubs/{clubid}/events/{eventid}/shifts` | `GET /api/v2/Shifts?$filter=eventId eq '{eventid}'` | |
| `GET /api/v1/clubs/{clubid}/shifts/{id}/members` | `GET /api/v2/ShiftMembers?$filter=shiftId eq '{id}'&$expand=User` | |
| `POST /api/v1/clubs/{clubid}/shifts/{id}/members` | `POST /api/v2/ShiftMembers` | |
| `DELETE /api/v1/clubs/{clubid}/shifts/{shiftid}/members/{id}` | `DELETE /api/v2/ShiftMembers('{id}')` | |
| **Fines** | | |
| `GET /api/v1/me/fines` | `GET /api/v2/Fines?$filter=userId eq '{currentUserId}' and paid eq false&$expand=Club` | |
| `GET /api/v1/clubs/{clubid}/fines` | `GET /api/v2/Fines?$filter=clubId eq '{clubid}'` | |
| `POST /api/v1/clubs/{clubid}/fines` | `POST /api/v2/Fines` | |
| `DELETE /api/v1/clubs/{clubid}/fines/{id}` | `DELETE /api/v2/Fines('{id}')` | |
| **Fine Templates** | | |
| `GET /api/v1/clubs/{clubid}/fine-templates` | `GET /api/v2/FineTemplates?$filter=clubId eq '{clubid}'` | |
| `POST /api/v1/clubs/{clubid}/fine-templates` | `POST /api/v2/FineTemplates` | |
| `PUT /api/v1/clubs/{clubid}/fine-templates/{id}` | `PATCH /api/v2/FineTemplates('{id}')` | |
| `DELETE /api/v1/clubs/{clubid}/fine-templates/{id}` | `DELETE /api/v2/FineTemplates('{id}')` | |
| **Invites** | | |
| `GET /api/v1/clubs/{clubid}/invites` | `GET /api/v2/Invites?$filter=clubId eq '{clubid}'` | Admin only |
| `POST /api/v1/clubs/{clubid}/invites` | `POST /api/v2/Invites` | |
| `GET /api/v1/invites` | `GET /api/v2/Invites?$filter=email eq '{currentUserEmail}'&$expand=Club` | |
| `POST /api/v1/invites/{id}/accept` | `POST /api/v2/Invites('{id}')/Accept` | Custom action |
| `POST /api/v1/invites/{id}/reject` | `POST /api/v2/Invites('{id}')/Reject` | Custom action |
| **Join Requests** | | |
| `GET /api/v1/clubs/{clubid}/joinRequests` | `GET /api/v2/JoinRequests?$filter=clubId eq '{clubid}'&$expand=User` | Admin only |
| `GET /api/v1/clubs/{clubid}/inviteLink` | `GET /api/v2/Clubs('{clubid}')/GetInviteLink()` | Custom function |
| `POST /api/v1/clubs/{clubid}/join` | `POST /api/v2/JoinRequests` | |
| `GET /api/v1/clubs/{clubid}/info` | `GET /api/v2/Clubs('{clubid}')?$select=id,name,description` | Add computed fields |
| `POST /api/v1/joinRequests/{id}/accept` | `POST /api/v2/JoinRequests('{id}')/Accept` | Custom action |
| `POST /api/v1/joinRequests/{id}/reject` | `POST /api/v2/JoinRequests('{id}')/Reject` | Custom action |
| **News** | | |
| `GET /api/v1/clubs/{clubid}/news` | `GET /api/v2/News?$filter=clubId eq '{clubid}'&$orderby=createdAt desc` | |
| `POST /api/v1/clubs/{clubid}/news` | `POST /api/v2/News` | |
| `GET /api/v1/clubs/{clubid}/news/{id}` | `GET /api/v2/News('{id}')` | |
| `PATCH /api/v1/clubs/{clubid}/news/{id}` | `PATCH /api/v2/News('{id}')` | |
| `DELETE /api/v1/clubs/{clubid}/news/{id}` | `DELETE /api/v2/News('{id}')` | |
| **Notifications** | | |
| `GET /api/v1/notifications` | `GET /api/v2/Notifications?$filter=userId eq '{currentUserId}'&$orderby=createdAt desc` | |
| `GET /api/v1/notifications/count` | `GET /api/v2/Notifications/$count?$filter=userId eq '{currentUserId}' and isRead eq false` | |
| `GET /api/v1/notifications/{id}` | `GET /api/v2/Notifications('{id}')` | |
| `PATCH /api/v1/notifications/{id}` | `PATCH /api/v2/Notifications('{id}')` | Mark as read |
| `POST /api/v1/notifications/mark-all-read` | `POST /api/v2/MarkAllNotificationsRead` | Unbound action |
| `GET /api/v1/notification-preferences` | `GET /api/v2/UserNotificationPreferences?$filter=userId eq '{currentUserId}'` | |
| `POST /api/v1/notification-preferences` | `POST /api/v2/UserNotificationPreferences` or `PATCH` | |
| **User Profile** | | |
| `GET /api/v1/me` | `GET /api/v2/Users('{currentUserId}')` | |
| `PUT /api/v1/me` | `PATCH /api/v2/Users('{currentUserId}')` | |
| **Privacy** | | |
| `GET /api/v1/me/privacy` | `GET /api/v2/UserPrivacySettings?$filter=userId eq '{currentUserId}' and clubId eq null` | Global settings |
| `POST /api/v1/me/privacy` | `POST /api/v2/UserPrivacySettings` or `PATCH` | |
| `GET /api/v1/me/privacy/clubs` | `GET /api/v2/UserPrivacySettings?$filter=userId eq '{currentUserId}' and clubId ne null&$expand=Club` | |
| **Club Settings** | | |
| `GET /api/v1/clubs/{clubid}/settings` | `GET /api/v2/ClubSettings?$filter=clubId eq '{clubid}'` | |
| `POST /api/v1/clubs/{clubid}/settings` | `POST /api/v2/ClubSettings` or `PATCH` | |
| **Dashboard** | | |
| `GET /api/v1/dashboard/news` | `GET /api/v2/GetDashboardNews()` | Unbound function |
| `GET /api/v1/dashboard/events` | `GET /api/v2/GetDashboardEvents()` | Unbound function |
| `GET /api/v1/dashboard/activities` | `GET /api/v2/GetDashboardActivities()` | Unbound function |
| **Search** | | |
| `GET /api/v1/search?q={query}` | `GET /api/v2/SearchGlobal(query='{query}')` | Unbound function |
| **Health** | | |
| `GET /api/v1/health` | `GET /api/v1/health` | Keep as REST |

## Benefits of OData Migration

### 1. Standardization
- Industry-standard protocol with extensive tooling
- Consistent query syntax across all entities
- Self-documenting via metadata endpoint

### 2. Flexibility
- Clients can request exactly what they need ($select, $expand)
- Powerful filtering without custom endpoints
- Built-in pagination and sorting

### 3. Reduced Backend Code
- Eliminate custom query endpoints
- Automatic handling of common operations
- Less boilerplate code to maintain

### 4. Enhanced Frontend Capabilities
- Complex queries without backend changes
- Reduce number of API calls via $expand
- Better performance with selective field retrieval

### 5. Future-Proofing
- Add new fields without breaking clients
- Support for advanced features (delta queries, batch operations)
- Easier to generate client libraries

## Risks & Mitigations

### Risk 1: Learning Curve
**Mitigation:**
- Comprehensive documentation
- Training sessions for team
- Gradual migration with parallel operation

### Risk 2: OData Complexity
**Mitigation:**
- Start with simple queries
- Create helper utilities for common patterns
- Provide examples and templates

### Risk 3: Performance
**Mitigation:**
- Implement proper indexing
- Use read hooks for query optimization
- Monitor and profile query performance
- Set reasonable $top limits

### Risk 4: File Upload Complexity
**Mitigation:**
- Document custom file upload endpoints
- Provide clear examples
- Consider OData media entities in future

### Risk 5: Breaking Changes
**Mitigation:**
- Run both APIs in parallel
- Gradual deprecation of old endpoints
- Version management strategy

## Success Metrics

1. **API Response Time:** < 200ms for simple queries
2. **Test Coverage:** > 90% for OData endpoints
3. **Frontend Migration:** 100% of features migrated
4. **OData Compliance:** Pass all go-odata compliance tests
5. **Developer Satisfaction:** Positive feedback on developer experience

## Conclusion

This migration plan provides a comprehensive roadmap for transitioning the Clubs backend to OData v4 APIs while maintaining complete feature parity. The phased approach allows for incremental development and testing, minimizing risk and ensuring a smooth transition. The use of `github.com/NLstn/go-odata` provides a robust, standards-compliant foundation that will serve the project well into the future.

## Next Steps

1. **Week 1:** Review and approval of migration plan
2. **Week 2:** Begin Phase 1 implementation
3. **Weekly Reviews:** Progress check-ins and adjustments
4. **Week 8:** Complete backend migration
5. **Week 12:** Complete frontend migration
6. **Week 13:** Deprecate old API, celebrate success! 🎉
