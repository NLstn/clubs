# Feature Settings and OData API Mapping

## Overview

This document provides the required analysis of toggleable settings and their corresponding OData APIs/entities, as requested in the feature specification.

## Settings Analysis

### Source: ClubSettings Model

Location: `Backend/models/club_settings.go`

```go
type ClubSettings struct {
    FinesEnabled             bool  // Toggles fines feature
    ShiftsEnabled            bool  // Toggles shifts feature
    TeamsEnabled             bool  // Toggles teams feature
    MembersListVisible       bool  // Controls member list visibility (not API-toggling)
    DiscoverableByNonMembers bool  // Controls club discoverability (not API-toggling)
}
```

### API-Controlling Settings

Three settings control backend API access:

1. **FinesEnabled** - Controls Fine and FineTemplate entities
2. **ShiftsEnabled** - Controls Shift and ShiftMember entities
3. **TeamsEnabled** - Controls Team and TeamMember entities

Note: `MembersListVisible` and `DiscoverableByNonMembers` control visibility/discoverability but do not disable APIs entirely.

## Settings-to-Entity-to-API Mapping

### 1. Fines Feature (`FinesEnabled`)

#### Entities
- **Fine** (`Backend/models/fines.go`)
  - Represents a fine assigned to a club member
  - OData hooks: BeforeReadCollection, BeforeReadEntity, BeforeCreate, BeforeUpdate, BeforeDelete

- **FineTemplate** (`Backend/models/fine_templates.go`)
  - Represents a reusable fine template for quick fine creation
  - OData hooks: BeforeReadCollection, BeforeReadEntity, BeforeCreate, BeforeUpdate, BeforeDelete

#### OData Endpoints
- `GET /api/v2/Fines` - List all fines
- `GET /api/v2/Fines('id')` - Get specific fine
- `POST /api/v2/Fines` - Create fine
- `PUT/PATCH /api/v2/Fines('id')` - Update fine
- `DELETE /api/v2/Fines('id')` - Delete fine
- `GET /api/v2/FineTemplates` - List all fine templates
- `GET /api/v2/FineTemplates('id')` - Get specific template
- `POST /api/v2/FineTemplates` - Create template
- `PUT/PATCH /api/v2/FineTemplates('id')` - Update template
- `DELETE /api/v2/FineTemplates('id')` - Delete template

#### Behavior When Disabled
- Entity requests with ID: **HTTP 400** (FeatureDisabled error)
- Collection queries: Return empty results (filtered by OData hooks)
- Create operations: Validated in OData hooks

---

### 2. Shifts Feature (`ShiftsEnabled`)

#### Entities
- **Shift** (`Backend/models/shift_schedules.go`)
  - Represents a time slot for an event that members can sign up for
  - OData hooks: BeforeReadCollection, BeforeReadEntity, BeforeCreate, BeforeUpdate, BeforeDelete

- **ShiftMember** (`Backend/models/shift_schedules.go`)
  - Represents a member assigned to a shift
  - OData hooks: BeforeReadCollection, BeforeReadEntity, BeforeCreate, BeforeUpdate, BeforeDelete

#### OData Endpoints
- `GET /api/v2/Shifts` - List all shifts
- `GET /api/v2/Shifts('id')` - Get specific shift
- `POST /api/v2/Shifts` - Create shift
- `PUT/PATCH /api/v2/Shifts('id')` - Update shift
- `DELETE /api/v2/Shifts('id')` - Delete shift
- `GET /api/v2/ShiftMembers` - List shift assignments
- `GET /api/v2/ShiftMembers('id')` - Get specific assignment
- `POST /api/v2/ShiftMembers` - Assign member to shift
- `PUT/PATCH /api/v2/ShiftMembers('id')` - Update assignment
- `DELETE /api/v2/ShiftMembers('id')` - Remove assignment

#### Behavior When Disabled
- Entity requests with ID: **HTTP 400** (FeatureDisabled error)
- Collection queries: Return empty results (filtered by OData hooks)
- Create operations: Validated in OData hooks

---

### 3. Teams Feature (`TeamsEnabled`)

#### Entities
- **Team** (`Backend/models/teams.go`)
  - Represents a sub-group within a club
  - OData hooks: BeforeReadCollection, BeforeReadEntity, BeforeCreate, BeforeUpdate, BeforeDelete

- **TeamMember** (`Backend/models/teams.go`)
  - Represents a member's membership in a team
  - OData hooks: BeforeReadCollection, BeforeReadEntity, BeforeCreate, BeforeUpdate, BeforeDelete

#### OData Endpoints
- `GET /api/v2/Teams` - List all teams
- `GET /api/v2/Teams('id')` - Get specific team
- `POST /api/v2/Teams` - Create team
- `PUT/PATCH /api/v2/Teams('id')` - Update team
- `DELETE /api/v2/Teams('id')` - Delete team
- `GET /api/v2/TeamMembers` - List team memberships
- `GET /api/v2/TeamMembers('id')` - Get specific membership
- `POST /api/v2/TeamMembers` - Add member to team
- `PUT/PATCH /api/v2/TeamMembers('id')` - Update membership
- `DELETE /api/v2/TeamMembers('id')` - Remove member from team

#### Behavior When Disabled
- Entity requests with ID: **HTTP 400** (FeatureDisabled error)
- Collection queries: Return empty results (filtered by OData hooks)
- Create operations: Validated in OData hooks

---

## Implementation Summary

### Settings Check Mechanism

**Location**: Middleware-based with database lookups

**Process**:
1. Request arrives at OData API endpoint
2. Authentication middleware validates user
3. Feature check middleware:
   - Extracts entity set name from URL
   - Checks if entity requires feature validation
   - For entity-specific requests (with ID):
     - Looks up entity's club ID
     - Queries ClubSettings for that club
     - Returns HTTP 400 if feature disabled
   - For collections/creates: passes to OData hooks

**Configuration Source**: Database table `club_settings`

### Error Response Format

When a feature is disabled, endpoints return:

```json
{
  "error": {
    "code": "FeatureDisabled",
    "message": "bad request: [feature name] feature is disabled for this club"
  }
}
```

HTTP Status: **400 Bad Request**

### Coverage Summary

| Request Type | Feature Check | HTTP 400 |
|--------------|---------------|----------|
| GET /Entity('id') | ✅ Middleware | ✅ Yes |
| PUT/PATCH /Entity('id') | ✅ Middleware | ✅ Yes |
| DELETE /Entity('id') | ✅ Middleware | ✅ Yes |
| GET /Entities | ⚠️ OData Hooks | ❌ No (empty results) |
| POST /Entities | ⚠️ OData Hooks | ❌ No (validation error) |

### OData Hooks Applied

All six entities implement these OData lifecycle hooks:

1. **ODataBeforeReadCollection** - Filters entities by club membership
2. **ODataBeforeReadEntity** - Validates access to specific entity
3. **ODataBeforeCreate** - Validates creation permissions
4. **ODataBeforeUpdate** - Validates update permissions
5. **ODataBeforeDelete** - Validates deletion permissions

These hooks provide an additional layer of security and filtering beyond the middleware.

## Testing

**Test File**: `Backend/odata/feature_check_middleware_test.go`

**Coverage**:
- 8 comprehensive test cases
- Tests for each feature (Fines, Shifts, Teams)
- Tests for both enabled and disabled states
- Tests for non-feature entities (ensure no false positives)
- Tests for collection queries and metadata

**Results**: All tests passing with race detector enabled

## Documentation

- **Implementation Guide**: `Documentation/Backend/FeatureSettingsEnforcement.md`
- **API Documentation**: `Documentation/Backend/API.md`
- **Code Locations**:
  - Settings Model: `Backend/models/club_settings.go`
  - Middleware: `Backend/odata/feature_check_middleware.go`
  - Error Types: `Backend/models/errors.go`
  - Tests: `Backend/odata/feature_check_middleware_test.go`

## Conclusion

The implementation provides comprehensive feature toggle control for three major club features (Fines, Shifts, Teams), affecting six OData entities and their corresponding API endpoints. The middleware-based approach ensures HTTP 400 responses for disabled features while maintaining backward compatibility and security through OData hooks.
