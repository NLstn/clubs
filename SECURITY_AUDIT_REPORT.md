# Security Audit Report - Clubs Application
**Date:** December 19, 2025 (Updated)  
**Auditor:** GitHub Copilot Security Scan  
**Repository:** NLstn/clubs

## Executive Summary

A comprehensive security audit was conducted on the Clubs application, focusing on authorization, authentication, and data isolation. **Three vulnerabilities were discovered and fixed** (two CRITICAL from previous audit, one HIGH from this audit), along with several security improvements. No SQL injection, XSS, or other common web vulnerabilities were found.

### Severity Classification
- **CRITICAL**: Vulnerabilities that allow unauthorized access to resources or privilege escalation
- **HIGH**: Significant security issues that could lead to data exposure
- **MEDIUM**: Security concerns that should be addressed but have limited impact
- **LOW**: Minor security improvements

## Critical Vulnerabilities Found and Fixed

### 1. Cross-Club Resource Manipulation (CRITICAL - FIXED)

**Severity:** CRITICAL  
**Status:** ✅ FIXED  
**Impact:** Club isolation breach, unauthorized data manipulation

#### Description
A critical vulnerability was discovered where club administrators could manipulate resources (Events, Fines, Shifts) across club boundaries by specifying a TeamID or EventID from a different club.

#### Attack Scenario
1. Attacker is an admin of Club A
2. Attacker creates an event with ClubID = Club A but TeamID = Team from Club B
3. This bypasses club isolation and allows cross-club data manipulation

#### Affected Entities
- **Events**: TeamID validation missing
- **Fines**: TeamID validation missing
- **Shifts**: EventID validation missing

#### Fix Implemented
Added validation in `ODataBeforeCreate` and `ODataBeforeUpdate` hooks to verify:
- TeamID belongs to the specified ClubID (Events, Fines)
- EventID belongs to the specified ClubID (Shifts)

**Code Changes:**
- `Backend/models/events.go`: Lines 441-445, 474-478
- `Backend/models/fines.go`: Lines 157-161, 182-186
- `Backend/models/shift_schedules.go`: Lines 190-195, 213-218

**Test Coverage:**
- `Backend/models/security_audit_test.go`: TestEventCreationTeamIDClubIDMismatch
- `Backend/models/security_audit_test.go`: TestFineCreationTeamIDClubIDMismatch

---

### 2. Privilege Escalation via Role Manipulation (CRITICAL - FIXED)

**Severity:** CRITICAL  
**Status:** ✅ FIXED  
**Impact:** Unauthorized privilege escalation to owner role

#### Description
Two privilege escalation vulnerabilities were discovered:

1. **Flawed canChangeRole() Logic**: The function allowed admins to promote any "member" to "owner" role due to incorrect logic in the permission check
2. **Direct OData API Manipulation**: The `Member.ODataBeforeUpdate()` hook did not validate role changes, allowing direct PATCH requests to bypass proper authorization

#### Attack Scenarios

**Scenario 1: Admin promotes member to owner**
```javascript
// Admin updates a regular member
PATCH /api/v2/Members/{member-id}
{
  "Role": "owner"  // Should be rejected but was allowed
}
```

**Scenario 2: Admin self-promotion**
```javascript
// Admin promotes themselves
PATCH /api/v2/Members/{admin-own-id}
{
  "Role": "owner"  // Should be rejected but was allowed
}
```

#### Root Cause
Line 256 in `canChangeRole()`:
```go
// VULNERABLE CODE (before fix)
if changingUserRole == "admin" && (targetMember.Role == "member" || newRole == "admin") {
    return true, nil
}
```

This logic allows admins to change roles when `targetMember.Role == "member"`, regardless of what the `newRole` is.

#### Fix Implemented

1. **Fixed canChangeRole() logic**:
   - Admins can only change between "member" and "admin" roles
   - Admins CANNOT promote to owner or demote from owner
   - Only owners can manage owner roles

2. **Added role change validation in ODataBeforeUpdate()**:
   - Detects role changes and applies proper authorization
   - Uses the same `canChangeRole()` logic as the UpdateMemberRole function

3. **Added role restrictions in ODataBeforeCreate()**:
   - Only owners can create new owners
   - Only admins/owners can create new admins

**Code Changes:**
- `Backend/models/members.go`: Lines 236-271 (canChangeRole fix)
- `Backend/models/members.go`: Lines 333-373 (ODataBeforeUpdate fix)
- `Backend/models/members.go`: Lines 310-343 (ODataBeforeCreate fix)

**Test Coverage:**
- `Backend/models/security_audit_test.go`: TestPrivilegeEscalationViaRoleUpdate
- `Backend/models/security_audit_test.go`: TestPrivilegeEscalationViaCreate

---

### 3. TeamMember Privilege Escalation via Role Update (HIGH - FIXED)

**Severity:** HIGH  
**Status:** ✅ FIXED  
**Impact:** Team-level privilege escalation, unauthorized team admin promotion

#### Description
Similar to the Member privilege escalation vulnerability (#2), the `TeamMember.ODataBeforeUpdate()` hook did not validate role changes, allowing users to manipulate team roles via direct OData API calls.

#### Attack Scenarios

**Scenario 1: Regular team member self-promotes to team admin**
```javascript
// Team member updates their own role
PATCH /api/v2/TeamMembers/{member-id}
{
  "Role": "admin"  // Should be rejected but was allowed
}
```

**Scenario 2: Bypass last admin protection**
```javascript
// Demote the last team admin
PATCH /api/v2/TeamMembers/{last-admin-id}
{
  "Role": "member"  // Should be rejected to maintain at least one admin
}
```

#### Root Cause
The `TeamMember.ODataBeforeUpdate()` function (lines 567-596 in `teams.go`) only checked if the user had permission to update team members generally, but did not:
- Detect role changes
- Validate role change permissions using `canChangeRole()`
- Prevent demoting the last team admin
- Enforce immutability of TeamID

#### Fix Implemented

1. **Added role change detection and validation**:
   - Detects when the Role field is being changed
   - Uses existing `canChangeRole()` logic for authorization
   - Prevents demoting the last team admin (ErrLastTeamAdminDemotion)
   - Enforces TeamID immutability

2. **Maintains existing permission checks for non-role updates**:
   - Club admins/owners can update team members
   - Team admins can update team members
   - Users cannot change TeamID (immutable)

**Code Changes:**
- `Backend/models/teams.go`: Lines 567-625 (ODataBeforeUpdate fix)

**Test Coverage:**
- `Backend/models/team_security_test.go`: TestTeamMemberPrivilegeEscalationViaRoleUpdate
- `Backend/models/team_security_test.go`: TestTeamMemberPrivilegeEscalationViaCreate

---

## Security Strengths

### ✅ Authentication & Authorization

1. **Multi-layered Authentication**
   - JWT Bearer tokens (15-min access, 30-day refresh)
   - API Key authentication with bcrypt hashing
   - Magic Link email authentication
   - OAuth/OIDC via Keycloak
   - Composite auth middleware supports multiple methods

2. **Rate Limiting**
   - Auth endpoints: 5 requests per minute per IP
   - API endpoints: 30 requests per 5 seconds per IP
   - Automatic cleanup of inactive rate limiters
   - Protection against brute force attacks

3. **OData Authorization Hooks**
   - All models implement ODataBefore* hooks
   - Read operations filter by club membership
   - Create/Update/Delete operations validate permissions
   - Proper separation between club admins, owners, and members

### ✅ Data Isolation

1. **Club Boundary Enforcement**
   - Users can only access data from clubs they belong to
   - All entity read operations filter by club membership
   - Cross-club queries are properly prevented

2. **Role-Based Access Control**
   - Owner: Full control over club
   - Admin: Management capabilities (events, fines, members)
   - Member: Read access and self-service operations

### ✅ SQL Injection Prevention

- **All database queries use parameterized queries** with `?` placeholders
- No string concatenation found in SQL statements
- GORM ORM provides additional protection

### ✅ Session Management

- Refresh tokens hashed with SHA-256
- Session tracking by User-Agent and IP
- Users can view and revoke active sessions
- Expired sessions automatically filtered

---

## Additional Security Observations

### API Key Security

**Strengths:**
- API keys use bcrypt hashing (cost factor 12)
- 256-bit cryptographic random generation
- Keys stored as hashes, never in plaintext
- Prefix-based filtering for efficient validation
- Expiration and active status tracking
- Last used timestamp for auditing

**Recommendations:**
- Consider implementing API key scopes/permissions (partially implemented)
- Add API key usage logging for security monitoring
- Implement key rotation policies

### Input Validation

**Current State:**
- GORM provides type safety for structured data
- OData hooks validate foreign key relationships
- No obvious XSS vulnerabilities (backend doesn't render HTML)

**Recommendations:**
- Add explicit input sanitization for user-provided strings
- Implement max length restrictions on text fields
- Add validation for email formats, URLs, etc.

### CORS Configuration

**Current Implementation:**
```go
Access-Control-Allow-Origin: *
```

**Recommendation:**
- Configure specific allowed origins in production
- Use environment variables for origin configuration
- Consider implementing CORS preflight caching

---

## Testing Summary

### Security Tests Created
- ✅ Cross-club boundary violation tests
- ✅ Privilege escalation prevention tests (Member)
- ✅ Privilege escalation prevention tests (TeamMember) - NEW
- ✅ Role permission enforcement tests
- ✅ API key validation tests
- ✅ Team role change validation tests - NEW

### Test Results
- All backend tests pass: ✅ (8 packages, 0.322s)
- Security audit tests pass: ✅ (10 tests total)
  - Member privilege escalation tests: 4 tests
  - TeamMember privilege escalation tests: 6 tests (NEW)
- CodeQL security scan: ✅ (0 alerts)

### Test Commands
```bash
# Run all tests
cd Backend && go test ./...

# Run security-specific tests
cd Backend && go test -v -run ".*[Ss]ecurity.*|.*[Pp]rivilege.*|.*[Aa]udit.*" ./models/

# Run CodeQL scan
# (Already integrated in CI/CD)
```

---

## Recommendations

### High Priority

1. ~~**Review Team Role Management** (Similar to member roles)~~ ✅ **COMPLETED**
   - ~~Verify team admins cannot escalate to club owner~~
   - ~~Ensure proper team hierarchy enforcement~~
   - **Status:** Fixed in this audit. TeamMember role changes now properly validated.

2. **Implement API Key Scopes**
   - Define granular permissions for API keys
   - Limit API key capabilities based on use case
   - Document API key best practices

3. **Add Security Logging**
   - Log all authorization failures
   - Track failed login attempts
   - Monitor unusual API patterns
   - Implement alert thresholds

### Medium Priority

4. **CORS Hardening**
   - Whitelist specific origins in production
   - Use environment-based configuration

5. **Input Validation Layer**
   - Add explicit max length validations
   - Implement email format validation
   - Sanitize user-provided text fields

6. **Session Security**
   - Consider adding session fingerprinting
   - Implement suspicious activity detection
   - Add geographic location tracking (optional)

### Low Priority

7. **API Rate Limiting Improvements**
   - Per-user rate limits in addition to per-IP
   - Configurable rate limit tiers
   - Rate limit exemptions for trusted clients

8. **Audit Logging**
   - Comprehensive audit trail for all changes
   - Queryable audit log interface
   - Retention policies for logs

---

## Compliance Notes

### GDPR Considerations
- User data access is properly restricted
- Users can delete their own accounts (leave clubs)
- Email addresses properly protected from unauthorized access
- Session management allows users to view/revoke access

### Data Privacy
- Club members cannot see members of other clubs
- Email addresses only visible to club admins/owners
- User information properly scoped to club context

---

## Conclusion

The security audit identified and fixed **three vulnerabilities** that could have allowed:
1. Cross-club data manipulation (CRITICAL - Previous audit)
2. Unauthorized privilege escalation at club level (CRITICAL - Previous audit)
3. Unauthorized privilege escalation at team level (HIGH - This audit)

All vulnerabilities have been **completely remediated** with comprehensive fixes and test coverage.

The application demonstrates strong security practices overall, including:
- ✅ Proper authentication mechanisms (JWT, API Keys, Magic Links, OAuth/OIDC)
- ✅ Role-based access control at both club and team levels
- ✅ SQL injection prevention (parameterized queries)
- ✅ Data isolation between clubs
- ✅ Rate limiting protection
- ✅ Comprehensive authorization hooks on all entities
- ✅ Proper validation of role changes and last admin protections

**Overall Security Posture: EXCELLENT** (after all fixes)

The remaining recommendations are preventive measures and enhancements rather than critical fixes.

---

## Appendix: Files Modified

### Security Fixes (This Audit - December 19, 2025)
1. `Backend/models/teams.go` - TeamMember privilege escalation fix
   - Added role change detection in ODataBeforeUpdate
   - Enforced proper authorization using canChangeRole()
   - Added TeamID immutability check

### Test Coverage (This Audit)
2. `Backend/models/team_security_test.go` - New team security test suite
   - TestTeamMemberPrivilegeEscalationViaRoleUpdate (4 test cases)
   - TestTeamMemberPrivilegeEscalationViaCreate (2 test cases)

### Documentation (This Audit)
3. `SECURITY_AUDIT_REPORT.md` - Updated with new findings

### Security Fixes (Previous Audit)
4. `Backend/models/events.go` - Cross-club validation
5. `Backend/models/fines.go` - Cross-club validation
6. `Backend/models/shift_schedules.go` - Cross-club validation
7. `Backend/models/members.go` - Member privilege escalation fixes

### Test Coverage (Previous Audit)
8. `Backend/models/security_audit_test.go` - Security test suite for Members

---

**Report Generated:** December 19, 2025  
**Next Review Recommended:** Q2 2026 or upon significant architecture changes
