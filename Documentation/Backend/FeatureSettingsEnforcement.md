# Feature Settings Enforcement - HTTP 400 Response Implementation

## Overview

This document describes the implementation of HTTP 400 (Bad Request) responses when club features are disabled via settings.

## Feature Settings

The `ClubSettings` model controls which features are enabled for each club:

```go
type ClubSettings struct {
    FinesEnabled   bool  // Controls Fine and FineTemplate APIs
    ShiftsEnabled  bool  // Controls Shift and ShiftMember APIs
    TeamsEnabled   bool  // Controls Team and TeamMember APIs
    // ... other settings
}
```

## Affected OData Entities

### Fines Feature (`FinesEnabled`)
- **Entities**: `Fine`, `FineTemplate`
- **Endpoints**: `/api/v2/Fines`, `/api/v2/FineTemplates`

### Shifts Feature (`ShiftsEnabled`)
- **Entities**: `Shift`, `ShiftMember`
- **Endpoints**: `/api/v2/Shifts`, `/api/v2/ShiftMembers`

### Teams Feature (`TeamsEnabled`)
- **Entities**: `Team`, `TeamMember`
- **Endpoints**: `/api/v2/Teams`, `/api/v2/TeamMembers`

## Implementation

### Architecture

The feature enforcement is implemented as middleware that sits between authentication and the OData service:

```
Request → Auth Middleware → Feature Check Middleware → OData Service
```

### Middleware Behavior

The `FeatureCheckMiddleware` intercepts requests and:

1. Extracts the entity set name from the request path
2. Determines if the entity requires a feature check
3. For requests with entity IDs (e.g., `/Fines('uuid')`):
   - Looks up the entity in the database
   - Retrieves the associated club ID
   - Checks if the feature is enabled for that club
   - Returns HTTP 400 if disabled
4. For other requests (collections, creates):
   - Passes through to OData hooks
   - OData hooks handle authorization and filtering

### HTTP 400 Response Format

When a feature is disabled, the middleware returns an OData v4 compliant error:

```json
{
  "error": {
    "code": "FeatureDisabled",
    "message": "bad request: fines feature is disabled for this club"
  }
}
```

### Code Components

#### 1. Error Types (`Backend/models/errors.go`)

```go
type FeatureDisabledError struct {
    Feature string
    Message string
}

func NewFeatureDisabledError(featureName string) *FeatureDisabledError
func CheckFeatureEnabled(clubID, featureName string) error
func IsFeatureEnabled(clubID, featureName string) bool
```

#### 2. Middleware (`Backend/odata/feature_check_middleware.go`)

```go
func FeatureCheckMiddleware() func(http.Handler) http.Handler
```

Key functions:
- `getClubIDFromRequest()` - Extracts club ID from request
- `getClubIDFromEntity()` - Looks up entity's club ID in database
- `writeODataError()` - Writes OData-compliant error response

#### 3. Integration (`Backend/main.go`)

```go
odataWithFeatureCheck := odata.FeatureCheckMiddleware()(odataV2Mux)
odataWithAuth := http.StripPrefix("/api/v2", odata.AuthMiddleware(jwtSecret)(odataWithFeatureCheck))
```

## Request Handling Details

### Fully Handled (Returns HTTP 400)
- **GET** with entity ID: `/Fines('abc-123')`
- **PUT/PATCH** with entity ID: `/Fines('abc-123')`
- **DELETE** with entity ID: `/Fines('abc-123')`

For these requests, the middleware can determine the club ID by looking up the entity in the database before passing to OData.

### Pass-Through to OData Hooks
- **GET** collections: `/Fines`
- **POST** creates: `/Fines`

For these requests:
- Collections: Cannot determine a single club ID. OData hooks filter results by club membership.
- Creates: Club ID is in request body. OData hooks validate and check settings.

This is acceptable because:
1. Collection queries already filter by club membership in OData hooks
2. Disabled features return empty collections (correct behavior)
3. Create operations validate club ID and settings in OData hooks

## Testing

### Test Coverage

Comprehensive tests in `Backend/odata/feature_check_middleware_test.go`:

1. **TestFeatureCheckMiddleware_FinesDisabled** - Verifies HTTP 400 for disabled fines
2. **TestFeatureCheckMiddleware_FinesEnabled** - Verifies pass-through when enabled
3. **TestFeatureCheckMiddleware_TeamsDisabled** - Verifies HTTP 400 for disabled teams
4. **TestFeatureCheckMiddleware_ShiftsDisabled** - Verifies HTTP 400 for disabled shifts
5. **TestFeatureCheckMiddleware_FineTemplateDisabled** - Verifies template handling
6. **TestFeatureCheckMiddleware_NonFeatureEntity** - Verifies non-feature entities pass through
7. **TestFeatureCheckMiddleware_CollectionQuery** - Verifies collections pass through
8. **TestFeatureCheckMiddleware_Metadata** - Verifies metadata requests unaffected

### Running Tests

```bash
cd Backend
go test ./odata -run TestFeatureCheckMiddleware -v
```

All tests pass with race detector enabled.

## Usage Examples

### Scenario 1: Feature Disabled

**Request:**
```http
GET /api/v2/Fines('abc-123') HTTP/1.1
Authorization: Bearer <token>
```

**Response** (when `FinesEnabled = false`):
```http
HTTP/1.1 400 Bad Request
Content-Type: application/json
OData-Version: 4.01

{
  "error": {
    "code": "FeatureDisabled",
    "message": "bad request: fines feature is disabled for this club"
  }
}
```

### Scenario 2: Feature Enabled

**Request:**
```http
GET /api/v2/Fines('abc-123') HTTP/1.1
Authorization: Bearer <token>
```

**Response** (when `FinesEnabled = true`):
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "@odata.context": "$metadata#Fines/$entity",
  "ID": "abc-123",
  "ClubID": "club-456",
  "Reason": "Late to practice",
  "Amount": 10.00,
  ...
}
```

### Scenario 3: Collection Query

**Request:**
```http
GET /api/v2/Fines HTTP/1.1
Authorization: Bearer <token>
```

**Response** (when `FinesEnabled = false`):
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "@odata.context": "$metadata#Fines",
  "value": []
}
```

Note: Collection queries return empty results when features are disabled, as OData hooks filter by club membership and features.

## Known Limitations

1. **Collection Queries**: Cannot return HTTP 400 for collection queries because:
   - Results may span multiple clubs
   - Cannot determine a single club context
   - OData hooks handle filtering appropriately

2. **Create Operations**: Cannot check settings in middleware because:
   - Club ID is in request body (not URL)
   - Reading body would interfere with OData processing
   - OData hooks validate settings during creation

These limitations are by design and do not impact security or functionality.

## Future Enhancements

Possible improvements:

1. **Request Body Parsing**: Parse POST request bodies to check settings before OData processing
2. **Bulk Operations**: Handle bulk create/update operations with feature checks
3. **Caching**: Cache club settings to reduce database queries
4. **Metrics**: Add prometheus metrics for feature-disabled rejections

## Security Considerations

- The middleware runs after authentication, ensuring only authenticated users are checked
- Club membership is verified by OData hooks (separate layer of security)
- Settings are stored per-club, not globally, allowing fine-grained control
- The middleware does not expose which clubs exist (returns 400, not 404)

## Performance Impact

- Minimal: Only affects requests with entity IDs
- Database lookup: One additional query per request (club settings lookup)
- Can be optimized with caching if needed
- Collection queries and metadata requests unaffected

## Maintenance

When adding new feature-controlled entities:

1. Add the feature flag to `ClubSettings` model
2. Add entity-to-feature mapping in `entitySetToFeature` map in middleware
3. Create model with proper OData hooks
4. Add tests for the new feature

## References

- ClubSettings Model: `Backend/models/club_settings.go`
- Middleware Implementation: `Backend/odata/feature_check_middleware.go`
- Error Types: `Backend/models/errors.go`
- Tests: `Backend/odata/feature_check_middleware_test.go`
