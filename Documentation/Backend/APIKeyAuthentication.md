# API Key Authentication Implementation Plan

## Overview

API Key authentication will provide a secure way for service-to-service communication and programmatic API access. Unlike JWT tokens (which expire quickly and require refresh), API Keys will be long-lived credentials tied to specific users or services with configurable permissions.

---

## Implementation Plan

### **Phase 1: Database Schema & Models**

#### 1.1 Create API Key Model
**File:** `/workspace/Backend/models/api_key.go`

```go
type APIKey struct {
    ID          string    `json:"ID" gorm:"type:uuid;default:gen_random_uuid();primaryKey" odata:"key"`
    UserID      string    `json:"UserID" gorm:"type:uuid;not null" odata:"required"`
    Name        string    `json:"Name" gorm:"not null" odata:"required"`
    KeyHash     string    `json:"-" gorm:"uniqueIndex;not null"` // Never exposed via API
    KeyPrefix   string    `json:"KeyPrefix" gorm:"not null" odata:"immutable"`
    Permissions []string  `json:"Permissions" gorm:"type:text[]" odata:"nullable"`
    LastUsedAt  *time.Time `json:"LastUsedAt,omitempty" gorm:"type:timestamp" odata:"nullable"`
    ExpiresAt   *time.Time `json:"ExpiresAt,omitempty" gorm:"type:timestamp" odata:"nullable"`
    IsActive    bool      `json:"IsActive" gorm:"default:true" odata:"required"`
    CreatedAt   time.Time `json:"CreatedAt" odata:"immutable"`
    UpdatedAt   time.Time `json:"UpdatedAt"`
}
```

**Key Design Decisions:**
- Store only **hashed** keys (SHA-256), never plaintext
- Use `KeyPrefix` for easy identification without revealing the full key (e.g., "sk_live_abc12345")
- Support **optional expiration** for enhanced security
- Track `LastUsedAt` for auditing and monitoring
- Support granular **permissions** (future-proof for scope-based access control)
- Use PascalCase for JSON fields (OData v4 convention)
- `KeyHash` is never exposed via API (`json:"-"` tag)

#### 1.2 Database Migration
- Add `APIKey` model to `Backend/database/database.go` AutoMigrate
- Ensure proper indexes:
  - Unique index on `key_hash`
  - Index on `user_id` for faster queries
  - Index on `expires_at` for cleanup queries

---

### **Phase 2: Backend Authentication Logic**

#### 2.1 API Key Generation Function
**File:** `/workspace/Backend/auth/auth.go`

```go
func GenerateAPIKey(prefix string) (plainKey string, keyHash string, keyPrefix string, error)
```

**Implementation Details:**
- Generate cryptographically secure random key (32+ bytes using `crypto/rand`)
- Format: `sk_live_<random_string>` or `sk_test_<random_string>`
- Return three values:
  - `plainKey`: Full key shown **once** to user
  - `keyHash`: SHA-256 hash for database storage
  - `keyPrefix`: First 8-12 chars for identification (e.g., "sk_live_abc12345")
- Example: `sk_live_abc12345def67890ghij12345klmno67890`

#### 2.2 API Key Validation Function
**File:** `/workspace/Backend/auth/auth.go`

```go
func ValidateAPIKey(keyStr string) (userID string, permissions []string, error)
```

**Implementation Details:**
- Hash the provided key using SHA-256
- Look up hash in database with single query
- Validate:
  - Key exists
  - `IsActive` is true
  - `ExpiresAt` is null or in the future
- Update `LastUsedAt` timestamp asynchronously
- Return user ID and permissions for authorization checks
- Return descriptive errors for logging/debugging

#### 2.3 API Key Middleware
**File:** `/workspace/Backend/handlers/middlewares.go`

```go
func APIKeyAuthMiddleware(next http.Handler) http.Handler
```

**Implementation Details:**
- Check for API key in two possible locations:
  1. `X-API-Key` header (preferred)
  2. `Authorization: ApiKey <key>` header (alternative)
- Validate key using `ValidateAPIKey()`
- Set user context using same pattern as JWT middleware:
  ```go
  ctx := context.WithValue(r.Context(), auth.UserIDKey, userID)
  ```
- Log authentication failures for security monitoring
- Continue to next handler on success

#### 2.4 Composite Auth Middleware
**File:** `/workspace/Backend/handlers/middlewares.go`

```go
func CompositeAuthMiddleware(next http.Handler) http.Handler
```

**Implementation Details:**
- Try JWT Bearer token authentication first (existing flow)
- If no Bearer token found, try API Key authentication
- If neither succeeds, return 401 Unauthorized
- **This allows endpoints to accept both auth methods seamlessly**
- Maintains backward compatibility with existing clients

---

### **Phase 3: API Endpoints**

#### 3.1 OData v4 Implementation
**File:** `/workspace/Backend/odata/entities.go`
- Register `APIKey` entity in OData service
- Configure entity set name: `APIKeys`
- Set up navigation properties to `User` entity

**File:** `/workspace/Backend/odata/custom_handlers.go`
- Implement custom create action that:
  - Generates API key using `GenerateAPIKey()`
  - Returns **plaintext key only once** (critical UX!)
  - Stores hashed version in database
  - Returns response with both key and metadata

**OData v4 Endpoints:**
- `POST /api/v4/APIKeys` - Create new API key
- `GET /api/v4/APIKeys` - List user's API keys (filtered by current user)
- `GET /api/v4/APIKeys('{id}')` - Get specific API key details
- `DELETE /api/v4/APIKeys('{id}')` - Revoke/delete API key
- `PATCH /api/v4/APIKeys('{id}')` - Update name, permissions, or active status

**Query Options Support:**
- `$filter` - Filter by Name, IsActive, etc.
- `$orderby` - Sort by CreatedAt, LastUsedAt
- `$select` - Select specific fields
- `$expand` - Expand User navigation property

#### 3.2 Custom Action for API Key Creation
**OData v4 Action:** `CreateAPIKey`

**Request:**
```json
{
  "Name": "Production Server",
  "ExpiresAt": "2026-12-31T23:59:59Z",
  "Permissions": ["read:events", "write:members"]
}
```

**Response (one-time only):**
```json
{
  "APIKey": "sk_live_abc12345def67890ghij12345klmno67890",
  "ID": "uuid-here",
  "KeyPrefix": "sk_live_abc12345",
  "Name": "Production Server",
  "ExpiresAt": "2026-12-31T23:59:59Z",
  "CreatedAt": "2025-12-13T10:30:00Z"
}
```

**Security Considerations:**
- Only show plaintext `APIKey` field at creation (never again)
- List/Get endpoints never return the `KeyHash` field
- Users can only manage their own keys (enforced by OData hooks)
- Admin users can manage any keys (future enhancement)
- Rate limit API key creation (max 10 keys per user, configurable)

---

### **Phase 4: Frontend Integration**

#### 4.1 API Key Management Page
**File:** `/workspace/Frontend/src/pages/APIKeysPage.tsx`

**Features:**
- ODataTable component showing existing API keys
- Columns: Key Prefix, Name, Last Used, Expires, Status
- "Create New Key" button
- Inline actions: Edit name, Revoke, Delete
- Empty state with helpful guidance

**Table Configuration:**
```typescript
const columns = [
  { field: 'KeyPrefix', label: 'Key Prefix', sortable: true },
  { field: 'Name', label: 'Name', sortable: true },
  { field: 'LastUsedAt', label: 'Last Used', type: 'datetime' },
  { field: 'ExpiresAt', label: 'Expires', type: 'datetime' },
  { field: 'IsActive', label: 'Status', type: 'boolean' }
];
```

#### 4.2 API Key Creation Modal
**File:** `/workspace/Frontend/src/components/dashboard/APIKeyModal.tsx`

**Critical UX Flow:**
1. User clicks "Create API Key" button
2. Modal opens with form:
   - **Name** (required): Text input for descriptive name
   - **Expiration** (optional): Date picker or preset options
   - **Permissions** (optional): Multi-select checkboxes (future)
3. On submit, backend returns plaintext key
4. **Modal transitions to "Key Created" view:**
   - Large, prominent display of full API key
   - Copy-to-clipboard button with success feedback
   - ⚠️ Warning message:
     > **Important:** This is the only time you'll see this key. Copy it now and store it securely. You won't be able to retrieve it again.
   - Checkbox: "I have saved this API key"
   - Close button (disabled until checkbox checked)
5. After modal closes, key cannot be retrieved again

**Component Structure:**
```typescript
<Modal>
  {!keyCreated ? (
    <APIKeyForm onSubmit={handleCreate} />
  ) : (
    <APIKeyDisplay
      apiKey={createdKey}
      onClose={handleClose}
    />
  )}
</Modal>
```

#### 4.3 TypeScript Types
**File:** `/workspace/Frontend/src/types/apikey.ts`

```typescript
interface APIKey {
  ID: string;
  UserID: string;
  Name: string;
  KeyPrefix: string;
  Permissions?: string[];
  LastUsedAt?: string;
  ExpiresAt?: string;
  IsActive: boolean;
  CreatedAt: string;
  UpdatedAt: string;
}

interface APIKeyCreationRequest {
  Name: string;
  ExpiresAt?: string;
  Permissions?: string[];
}

interface APIKeyCreationResponse extends APIKey {
  APIKey: string;  // Plaintext, shown once
}
```

#### 4.4 API Service
**File:** `/workspace/Frontend/src/utils/apiKeyService.ts`

```typescript
export const apiKeyService = {
  create: (data: APIKeyCreationRequest) => 
    api.post<APIKeyCreationResponse>('/api/v4/APIKeys', data),
  
  list: () => 
    api.get<{ value: APIKey[] }>('/api/v4/APIKeys'),
  
  get: (id: string) => 
    api.get<APIKey>(`/api/v4/APIKeys('${id}')`),
  
  update: (id: string, data: Partial<APIKey>) => 
    api.patch(`/api/v4/APIKeys('${id}')`, data),
  
  delete: (id: string) => 
    api.delete(`/api/v4/APIKeys('${id}')`)
};
```

---

### **Phase 5: Documentation**

#### 5.1 Update API Documentation
**File:** `/workspace/Documentation/Backend/API.md`

Add new section after existing authentication:

```markdown
## API Key Authentication

### Overview
API Keys provide long-lived credentials for programmatic access to the Clubs API. Unlike JWT tokens, API keys don't expire automatically (unless configured) and are ideal for:
- Server-to-server integrations
- CI/CD pipelines
- Third-party service integrations
- Long-running scripts or automation

### Authentication Methods
The API accepts API keys in two formats:

**Preferred (Custom Header):**
```
X-API-Key: sk_live_abc12345def67890ghij12345klmno67890
```

**Alternative (Authorization Header):**
```
Authorization: ApiKey sk_live_abc12345def67890ghij12345klmno67890
```

### Security Best Practices
1. **Store securely:** Treat API keys like passwords. Use environment variables or secret managers.
2. **Never commit:** Don't commit API keys to version control.
3. **Use descriptive names:** Name keys by their purpose (e.g., "Production CI/CD").
4. **Rotate regularly:** Create new keys and revoke old ones periodically.
5. **Set expiration dates:** Use expiration dates for temporary integrations.
6. **Monitor usage:** Check "Last Used" dates to identify unused keys.
7. **Revoke immediately:** If a key is compromised, revoke it immediately.

### API Key Prefixes
- `sk_live_*` - Production/live environment keys
- `sk_test_*` - Development/test environment keys (future)

### Rate Limits
API key endpoints have the following rate limits:
- Key creation: 5 per hour per user
- Key listing: Standard API rate limits
- API calls using keys: Standard API rate limits
```

#### 5.2 Create Comprehensive API Key Guide
**File:** `/workspace/Documentation/Backend/APIKeys.md`

**Contents:**
- Introduction and use cases
- When to use API keys vs JWT tokens
- Step-by-step guide to create and use API keys
- Permission scopes (extensible for future)
- Key rotation strategies and best practices
- Security considerations and threat model
- Troubleshooting common issues
- FAQ section

#### 5.3 Update Frontend Documentation
**File:** `/workspace/Documentation/Frontend/README.md`

Add section about API Key management UI components and patterns.

---

### **Phase 6: Testing**

#### 6.1 Backend Unit Tests

**File:** `/workspace/Backend/auth/auth_test.go`

Test cases:
- `TestGenerateAPIKey`: Verify key format, uniqueness, hash generation
- `TestValidateAPIKey_Valid`: Valid key returns correct user ID
- `TestValidateAPIKey_Invalid`: Invalid key returns error
- `TestValidateAPIKey_Expired`: Expired key returns specific error
- `TestValidateAPIKey_Inactive`: Inactive key returns error
- `TestValidateAPIKey_UpdatesLastUsed`: Verify LastUsedAt is updated

**File:** `/workspace/Backend/models/api_key_test.go`

Test cases:
- Model creation and validation
- GORM relationships with User model
- Permission array serialization/deserialization

**File:** `/workspace/Backend/handlers/middlewares_test.go`

Test cases:
- `TestAPIKeyAuthMiddleware_ValidKey`: Accept valid key
- `TestAPIKeyAuthMiddleware_InvalidKey`: Reject invalid key
- `TestAPIKeyAuthMiddleware_MissingKey`: Return 401 when no key provided
- `TestAPIKeyAuthMiddleware_HeaderFormats`: Test both X-API-Key and Authorization formats
- `TestCompositeAuthMiddleware_JWT`: JWT auth still works
- `TestCompositeAuthMiddleware_APIKey`: API key auth works
- `TestCompositeAuthMiddleware_Neither`: Reject when neither provided

#### 6.2 OData Integration Tests

**File:** `/workspace/Backend/odata/api_key_test.go`

Test cases:
- `TestCreateAPIKey`: Create key and verify response includes plaintext
- `TestCreateAPIKey_DuplicateName`: Allow duplicate names (per user)
- `TestListAPIKeys`: List returns only user's keys
- `TestGetAPIKey`: Get specific key (no plaintext in response)
- `TestUpdateAPIKey`: Update name and permissions
- `TestDeleteAPIKey`: Delete/revoke key
- `TestAPIKeyFiltering`: Test $filter queries
- `TestAPIKeyOrdering`: Test $orderby queries
- `TestAPIKeyWithExpiredToken`: Cannot access with expired key

#### 6.3 Frontend Component Tests

**File:** `/workspace/Frontend/src/components/dashboard/__tests__/APIKeyModal.test.tsx`

Test cases:
- Render create form correctly
- Submit form creates API key
- Display created key with copy button
- Warning message is visible
- Close button disabled until confirmed
- Copy to clipboard functionality works

**File:** `/workspace/Frontend/src/pages/__tests__/APIKeysPage.test.tsx`

Test cases:
- Render table with API keys
- Create button opens modal
- Delete button revokes key
- Edit functionality works
- Empty state displays correctly

#### 6.4 Integration Tests

**File:** `/workspace/Backend/handlers/api_integration_test.go`

End-to-end scenarios:
1. Create user → Create API key → Use key to call protected endpoint
2. Create key → Revoke key → Verify key no longer works
3. Create key with expiration → Wait/mock time → Verify key expired
4. Create multiple keys → Verify all work independently
5. Test key with insufficient permissions (future)

---

### **Phase 7: Security & Production Readiness**

#### 7.1 Rate Limiting

**Implementation:**
- Apply stricter rate limits to API key creation endpoint:
  - 5 keys per hour per user
  - Maximum 10 active keys per user
- Track failed authentication attempts:
  - Lock key after 10 failed attempts in 5 minutes
  - Send notification to user about suspicious activity
- Implement exponential backoff for repeated failures

**Configuration:**
```go
var apiKeyCreationLimiter = NewIPRateLimiter(rate.Limit(5.0/3600.0), 5)
const maxKeysPerUser = 10
```

#### 7.2 Audit Logging

**Events to Log:**
- API key created (user ID, key prefix, name)
- API key deleted/revoked (user ID, key prefix)
- API key updated (user ID, key prefix, changed fields)
- Failed authentication with API key (key prefix, IP, endpoint)
- Key used after long inactivity (potential compromised key)

**Log Format:**
```json
{
  "timestamp": "2025-12-13T10:30:00Z",
  "event": "api_key_created",
  "user_id": "uuid",
  "key_prefix": "sk_live_abc12345",
  "key_name": "Production Server",
  "ip_address": "192.168.1.1"
}
```

**Implementation:**
- Create `Backend/models/audit_log.go`
- Store in separate `audit_logs` table
- Implement automatic cleanup (retain 90 days)
- Add admin dashboard to view audit logs (future)

#### 7.3 Key Rotation

**Features:**
- Email notification 7 days before expiration
- Warning banner in UI when keys are expiring soon
- "Rotate Key" button that:
  1. Creates new key
  2. Shows both old and new keys
  3. Allows grace period (both keys work)
  4. Automatically revokes old key after grace period

**Implementation:**
- Background job checks for keys expiring in 7 days
- Sends email via Azure Communication Services
- Frontend component shows warning badges

#### 7.4 Permissions/Scopes (Future Enhancement)

**Scope System Design:**

Define permission scopes:
```go
const (
    ScopeReadEvents    = "read:events"
    ScopeWriteEvents   = "write:events"
    ScopeReadMembers   = "read:members"
    ScopeWriteMembers  = "write:members"
    ScopeReadFines     = "read:fines"
    ScopeWriteFines    = "write:fines"
    ScopeAdmin         = "admin"  // Full access
)
```

**Enforcement:**
- Middleware checks required scopes for each endpoint
- OData hooks verify permissions before CRUD operations
- Frontend shows scope selector during key creation
- Default to read-only scopes for safety

**Migration Path:**
- Phase 1: Add Permissions field (unused)
- Phase 2: Define scope constants and middleware
- Phase 3: Implement UI for scope selection
- Phase 4: Enforce scopes on all endpoints

#### 7.5 Additional Security Measures

**Key Lifecycle Management:**
- Automatic expiration of unused keys (no activity in 90 days)
- Maximum key lifetime (e.g., 1 year)
- Force rotation policy for security-critical applications

**Monitoring:**
- Dashboard showing:
  - Total active keys per user
  - Last used dates
  - Keys never used (potential forgotten keys)
  - Keys with suspicious activity patterns
- Alerts for unusual API key usage patterns

**Compliance:**
- GDPR: Include API keys in user data export/deletion
- SOC 2: Audit logging and key rotation policies
- Documentation for security audits

---

## Key Design Principles

### 1. Security First
- **Never store plaintext keys** - Only SHA-256 hashes
- **Show keys only once** - No retrieval after creation
- **Support expiration dates** - Time-bound credentials
- **Audit all operations** - Complete audit trail
- **Rate limit aggressively** - Prevent abuse
- **Revoke immediately** - Instant revocation without delay

### 2. OData v4 Compatibility
- **Use PascalCase** for all JSON fields
- **Follow entity registration patterns** from existing code
- **Integrate with existing middleware** in OData service
- **Support standard query options** ($filter, $orderby, $select, $expand)
- **Implement custom actions** for special operations like key creation

### 3. User Experience
- **Clear security warnings** - Impossible to miss important info
- **Easy copy-to-clipboard** - One-click copying
- **Descriptive names** - Help users organize multiple keys
- **Visual indicators** - Show expiration, last usage, status
- **Helpful empty states** - Guide users through first key creation
- **Confirmation dialogs** - Prevent accidental deletion

### 4. Backward Compatibility
- **Existing JWT auth continues** - No breaking changes
- **Composite middleware** - Both auth methods work simultaneously
- **No changes to existing endpoints** - Drop-in replacement
- **Optional migration path** - Gradual adoption possible

### 5. Developer Experience
- **Clear documentation** - Examples for every use case
- **Consistent API patterns** - Follows existing conventions
- **TypeScript types** - Full type safety in frontend
- **Error messages** - Descriptive errors for debugging
- **Testing utilities** - Easy to test API key auth in unit tests

---

## Implementation Order

### Phase 1: Foundation (Days 1-2)
1. Create `APIKey` model with proper GORM tags
2. Add database migration
3. Write model unit tests
4. Implement key generation function
5. Implement key validation function
6. Write auth function unit tests

### Phase 2: Authentication (Days 2-3)
1. Create API key middleware
2. Create composite auth middleware
3. Write middleware unit tests
4. Update existing routes to use composite auth
5. Test backward compatibility

### Phase 3: API Endpoints (Days 3-4)
1. Register APIKey entity in OData
2. Implement custom create action
3. Implement standard CRUD operations
4. Add OData hooks for security
5. Write OData integration tests

### Phase 4: Testing (Day 4-5)
1. Write comprehensive backend tests
2. End-to-end integration tests
3. Security vulnerability testing
4. Load testing for rate limits

### Phase 5: Frontend (Days 5-7)
1. Create TypeScript types
2. Implement API service
3. Build APIKeysPage component
4. Build APIKeyModal component
5. Add routing and navigation
6. Write frontend component tests

### Phase 6: Documentation (Day 7)
1. Update API.md with authentication methods
2. Create comprehensive APIKeys.md guide
3. Add inline code comments
4. Update README if needed

### Phase 7: Production Features (Days 8-10)
1. Implement audit logging
2. Add rate limiting
3. Create expiration notifications
4. Build monitoring dashboard
5. Security review and testing

---

## Estimated Effort

| Phase | Backend | Frontend | Testing | Documentation | Total |
|-------|---------|----------|---------|---------------|-------|
| Phase 1 | 4h | - | 2h | 0.5h | 6.5h |
| Phase 2 | 3h | - | 2h | 0.5h | 5.5h |
| Phase 3 | 4h | - | 3h | 1h | 8h |
| Phase 4 | - | - | 4h | - | 4h |
| Phase 5 | - | 6h | 2h | 1h | 9h |
| Phase 6 | - | - | - | 3h | 3h |
| Phase 7 | 4h | 2h | 2h | 1h | 9h |
| **Total** | **15h** | **8h** | **15h** | **7h** | **45h** |

**Note:** Times are estimates. Actual implementation may vary based on:
- Complexity of permission system
- Number of edge cases discovered during testing
- UI/UX refinements based on user feedback
- Security review findings

---

## Success Criteria

### Must Have (MVP)
- ✅ API keys can be created and stored securely (hashed)
- ✅ API keys work for authentication (alongside JWT)
- ✅ Users can list and delete their API keys
- ✅ Keys are shown only once at creation
- ✅ Basic rate limiting on creation
- ✅ Comprehensive test coverage (>80%)
- ✅ Documentation for developers

### Should Have (v1.1)
- ✅ Expiration date support
- ✅ Last used timestamp tracking
- ✅ Audit logging for security events
- ✅ Email notifications for expiration
- ✅ Admin view for all keys (security team)

### Could Have (v1.2+)
- ✅ Permission/scope system
- ✅ Key rotation workflow
- ✅ Usage analytics dashboard
- ✅ IP whitelisting
- ✅ Key templates for common scenarios
- ✅ Bulk operations (revoke multiple keys)

---

## Migration Strategy

### For Existing Users
1. API keys are an **additive feature** - no migration needed
2. JWT authentication continues to work unchanged
3. Users opt-in by creating their first API key
4. Documentation guides users when to use each auth method

### For Administrators
1. Monitor adoption via audit logs
2. Identify users creating many keys (potential integrations)
3. Reach out to power users for feedback
4. Create internal documentation for support team

### Rollback Plan
If issues arise:
1. Disable API key creation (keep existing keys working)
2. Revert composite middleware to JWT-only
3. Mark API key endpoints as maintenance mode
4. Fix issues and re-enable gradually

---

## Open Questions & Decisions

### 1. Key Format
**Decision:** Use `sk_live_` prefix for consistency with industry standards (Stripe, GitHub)
- Alternatives considered: `clubs_key_`, `ck_`
- Rationale: Familiar pattern, clear indication of key type

### 2. Maximum Keys Per User
**Decision:** 10 active keys per user initially
- Can be increased based on user feedback
- Prevents abuse while allowing reasonable integrations
- Admin users may have higher limits

### 3. Default Expiration
**Decision:** No default expiration (optional field)
- Users choose expiration when creating key
- Future: Add organization policies to enforce expiration
- Rationale: Flexibility for different use cases

### 4. Permission System Timeline
**Decision:** Implement field in Phase 1, enforce in future phase
- Allows data model to be future-proof
- Gives time to design comprehensive permission system
- Can gather user feedback on needed permissions

### 5. Test vs Live Keys
**Decision:** Only `sk_live_` prefix in initial release
- `sk_test_` prefix reserved for future multi-environment support
- Keeps initial implementation simpler
- Easy to add later without breaking changes

---

## Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Keys leaked to public repos | Medium | High | Clear warnings, documentation, git scanning tools |
| Keys never rotated | High | Medium | Expiration reminders, usage monitoring |
| Rate limit bypass | Low | Medium | Multiple rate limiting layers, monitoring |
| Hash collision | Very Low | High | Use SHA-256 (cryptographically secure) |
| Database breach | Low | High | Encryption at rest, access controls, audit logs |
| Key reuse across systems | Medium | Medium | Clear naming conventions, user education |
| Performance impact | Low | Low | Indexed queries, caching, async operations |

---

## Dependencies

### External Libraries
- None required (all functionality uses existing dependencies)

### Infrastructure
- PostgreSQL (existing)
- Azure Communication Services (existing, for email notifications)
- Redis (optional, for rate limiting cache in future)

### Internal Dependencies
- Existing auth package
- Existing OData service
- Existing middleware infrastructure
- Existing frontend components (Modal, Table, etc.)

---

## Post-Launch Activities

### Week 1
- Monitor error rates and authentication failures
- Gather user feedback on UX
- Fix critical bugs
- Update documentation based on questions

### Month 1
- Analyze usage patterns
- Identify most common use cases
- Plan permission system based on real needs
- Create video tutorials if needed

### Quarter 1
- Implement permission/scope system (Phase 7.4)
- Add usage analytics dashboard
- Enhance monitoring and alerting
- Consider IP whitelisting feature

---

## Appendix

### A. Example API Key Workflow

**1. User Creates Key:**
```http
POST /api/v4/APIKeys
Content-Type: application/json

{
  "Name": "Production CI/CD",
  "ExpiresAt": "2026-12-31T23:59:59Z"
}
```

**Response:**
```json
{
  "APIKey": "sk_live_abc12345def67890ghij12345klmno67890",
  "ID": "550e8400-e29b-41d4-a716-446655440000",
  "KeyPrefix": "sk_live_abc12345",
  "Name": "Production CI/CD",
  "ExpiresAt": "2026-12-31T23:59:59Z",
  "IsActive": true,
  "CreatedAt": "2025-12-13T10:30:00Z"
}
```

**2. User Makes API Call:**
```http
GET /api/v4/Events
X-API-Key: sk_live_abc12345def67890ghij12345klmno67890
```

**3. Backend Validates:**
- Hashes provided key
- Looks up hash in database
- Checks IsActive and ExpiresAt
- Updates LastUsedAt
- Sets user context
- Proceeds to handler

**4. User Lists Keys:**
```http
GET /api/v4/APIKeys
Authorization: Bearer jwt_token
```

**Response:**
```json
{
  "value": [
    {
      "ID": "550e8400-e29b-41d4-a716-446655440000",
      "KeyPrefix": "sk_live_abc12345",
      "Name": "Production CI/CD",
      "LastUsedAt": "2025-12-13T11:45:00Z",
      "ExpiresAt": "2026-12-31T23:59:59Z",
      "IsActive": true,
      "CreatedAt": "2025-12-13T10:30:00Z"
    }
  ]
}
```

Note: Full key is never returned after creation.

### B. Error Codes

| Code | Message | Scenario |
|------|---------|----------|
| 401 | API key is invalid | Key hash not found in database |
| 401 | API key is inactive | IsActive = false |
| 401 | API key has expired | ExpiresAt < now |
| 403 | Insufficient permissions | Permission check failed (future) |
| 429 | Too many API keys | User exceeded max keys limit |
| 429 | Rate limit exceeded | Too many creation requests |

### C. Database Schema

```sql
CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    key_hash VARCHAR(64) NOT NULL UNIQUE,
    key_prefix VARCHAR(20) NOT NULL,
    permissions TEXT[],
    last_used_at TIMESTAMP,
    expires_at TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_api_keys_user_id ON api_keys(user_id);
CREATE INDEX idx_api_keys_key_hash ON api_keys(key_hash);
CREATE INDEX idx_api_keys_expires_at ON api_keys(expires_at);
```

---

**Document Version:** 1.0  
**Last Updated:** December 13, 2025  
**Author:** Development Team  
**Status:** Approved for Implementation
