# Security Documentation

This document outlines the security measures implemented in the Clubs application and provides guidance for maintaining security.

## Table of Contents
1. [Authentication](#authentication)
2. [Authorization](#authorization)
3. [Security Fixes Applied](#security-fixes-applied)
4. [Security Best Practices](#security-best-practices)
5. [Deployment Security](#deployment-security)

## Authentication

The application supports multiple authentication methods:

### 1. JWT Bearer Tokens
- **Access Tokens**: 15-minute expiration
- **Refresh Tokens**: 30-day expiration
- **Algorithm**: HS256 (HMAC with SHA-256)
- **Storage**: Refresh tokens hashed with SHA-256 in database

### 2. API Keys
- **Format**: `prefix_base64RandomString`
- **Storage**: Bcrypt hashed (cost 12)
- **Validation**: Constant-time comparison via bcrypt
- **Features**: Expiration dates, active/inactive status, per-key permissions

### 3. Magic Link Email Authentication
- **Tokens**: 32-byte cryptographically secure random tokens
- **Expiration**: Configurable per magic link
- **Delivery**: Azure Communication Services

### 4. OAuth/OIDC via Keycloak
- **Provider**: Keycloak
- **Flows**: Authorization Code Flow with PKCE
- **Token Exchange**: Backend validates tokens with Keycloak

## Authorization

Authorization is implemented at multiple levels:

### OData v2 API Authorization

All models implement OData authorization hooks:

#### Read Authorization Hooks
- `ODataBeforeReadCollection`: Filters collections based on user permissions
- `ODataBeforeReadEntity`: Validates access to individual entities

#### Write Authorization Hooks
- `ODataBeforeCreate`: Validates creation permissions and sets audit fields
- `ODataBeforeUpdate`: Validates update permissions and sets audit fields
- `ODataBeforeDelete`: Validates deletion permissions

### Role-Based Access Control (RBAC)

#### Club Roles
- **Owner**: Full control, cannot be demoted if last owner
- **Admin**: Manage members, teams, events, fines, news
- **Member**: View club content, RSVP to events, view own fines

#### Team Roles
- **Admin**: Manage team members, team settings
- **Member**: View team content, participate in team activities

### Authorization Rules by Entity

| Entity | Read | Create | Update | Delete |
|--------|------|--------|--------|--------|
| **Club** | Members + Discoverable | Any user | Admin/Owner | Owner only |
| **Member** | Club members | Admin/Owner | Admin/Owner (role restrictions apply) | Admin/Owner or self |
| **Team** | Club members | Club Admin/Owner | Club Admin/Owner or Team Admin | Club Admin/Owner |
| **Event** | Club members | Admin/Owner | Admin/Owner | Admin/Owner |
| **Fine** | Club members | Admin/Owner | Admin/Owner | Admin/Owner |
| **News** | Club members | Admin/Owner | Admin/Owner | Admin/Owner |
| **Shift** | Club members | Admin/Owner | Admin/Owner | Admin/Owner |
| **EventRSVP** | Club members | Self only | Self only | Self or Admin/Owner |
| **APIKey** | Self only | Self only | Self only | Self only |
| **Invite** | Self + Admins | Admin/Owner | N/A | Self or Admin/Owner |
| **JoinRequest** | Self + Admins | Self only | N/A | Self or Admin/Owner |

## Security Fixes Applied

### 1. CORS Wildcard Vulnerability (CRITICAL)
**Date Fixed**: December 19, 2024

**Issue**: The CORS middleware used `Access-Control-Allow-Origin: *` which allows any website to make requests to the API, creating a significant security vulnerability.

**Impact**: 
- Any malicious website could make authenticated requests on behalf of users
- Credentials and sensitive data could be exposed to untrusted origins
- Session hijacking possible

**Fix**: 
```go
// Before (INSECURE)
w.Header().Set("Access-Control-Allow-Origin", "*")

// After (SECURE)
allowedOrigin := os.Getenv("FRONTEND_URL")
origin := r.Header.Get("Origin")
if origin == allowedOrigin {
    w.Header().Set("Access-Control-Allow-Origin", origin)
    w.Header().Set("Access-Control-Allow-Credentials", "true")
}
```

**Configuration Required**: Set `FRONTEND_URL` environment variable to your frontend domain (e.g., `https://app.example.com`)

### 2. X-Forwarded-For Trust Issue (HIGH)
**Date Fixed**: December 19, 2024

**Issue**: Rate limiting blindly trusted the `X-Forwarded-For` header, allowing attackers to bypass rate limits by spoofing IP addresses.

**Impact**:
- Attackers could bypass rate limits on authentication endpoints
- Brute force attacks on passwords and API keys would not be rate limited
- Denial of service attacks possible

**Fix**:
```go
// Priority order for IP extraction:
// 1. X-Real-IP (set by trusted reverse proxy)
// 2. X-Forwarded-For (first IP only - client IP)
// 3. RemoteAddr (direct connection)

ip := r.Header.Get("X-Real-IP")
if ip == "" {
    xff := r.Header.Get("X-Forwarded-For")
    if xff != "" {
        ips := strings.Split(xff, ",")
        ip = strings.TrimSpace(ips[0])
    }
}
if ip == "" {
    // Extract IP from RemoteAddr (format: "IP:port")
    addr := r.RemoteAddr
    if idx := strings.LastIndex(addr, ":"); idx != -1 {
        ip = addr[:idx]
    } else {
        ip = addr
    }
}
```

**Deployment Note**: Ensure your reverse proxy (nginx, load balancer, etc.) sets `X-Real-IP` header.

### 3. Team Authorization Enhancement (MEDIUM)
**Date Fixed**: December 19, 2024

**Issue**: `Team.ODataBeforeUpdate` only checked for club admin/owner permissions, not team admin permissions.

**Impact**:
- Team admins could not update their own teams
- Inconsistency with documented authorization model

**Fix**: Updated to use existing `CanUserEditTeam()` method which properly checks both club admins/owners AND team admins.

### 4. ClubSettings Audit Trail (MEDIUM)
**Date Fixed**: December 19, 2024

**Issue**: `ClubSettings.ODataBeforeUpdate` didn't set the `UpdatedBy` and `UpdatedAt` audit fields.

**Impact**:
- Loss of audit trail for club settings changes
- Cannot determine who made sensitive configuration changes

**Fix**: Added proper audit field setting in the update hook.

## Security Best Practices

### For Developers

1. **Never Trust Client Input**
   - Always validate and sanitize user input
   - Use parameterized queries (GORM handles this)
   - Validate foreign key relationships

2. **Follow Authorization Patterns**
   - Implement all OData authorization hooks
   - Use database-level filtering in read hooks
   - Validate ownership/permissions in write hooks

3. **Audit Trail**
   - Always set `CreatedBy`, `UpdatedBy`, `CreatedAt`, `UpdatedAt` fields
   - Log sensitive operations

4. **Rate Limiting**
   - Apply appropriate rate limits to all public endpoints
   - Use stricter limits for authentication endpoints
   - Monitor rate limit violations

5. **Error Messages**
   - Don't leak sensitive information in error messages
   - Use generic "Unauthorized" or "Forbidden" messages
   - Log detailed errors server-side only

### For DevOps/Deployment

1. **Environment Variables**
   - **Required**: `JWT_SECRET` - Use a strong, randomly generated secret (min 32 bytes)
   - **Required**: `FRONTEND_URL` - Set to your frontend domain (e.g., `https://app.example.com`)
   - Never commit secrets to version control

2. **HTTPS Only**
   - Always use HTTPS in production
   - Set secure cookie flags
   - Enable HSTS (HTTP Strict Transport Security)

3. **Reverse Proxy Configuration**
   ```nginx
   # nginx example
   location /api/ {
       proxy_pass http://backend:8080;
       proxy_set_header X-Real-IP $remote_addr;
       proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
       proxy_set_header X-Forwarded-Proto $scheme;
       proxy_set_header Host $host;
   }
   ```

4. **Database Security**
   - Use strong database passwords
   - Enable SSL/TLS for database connections in production
   - Restrict database access to application servers only

5. **API Key Security**
   - Rotate API keys regularly
   - Implement key expiration
   - Monitor API key usage
   - Revoke compromised keys immediately

## Deployment Security

### Production Checklist

- [ ] `JWT_SECRET` set to strong random value
- [ ] `FRONTEND_URL` set to production frontend domain
- [ ] HTTPS enabled with valid SSL certificate
- [ ] Database SSL/TLS enabled
- [ ] Reverse proxy configured with `X-Real-IP`
- [ ] Rate limiting tested and tuned
- [ ] Security headers enabled (CSP, X-Frame-Options, etc.)
- [ ] Logs configured for security monitoring
- [ ] Regular security updates scheduled
- [ ] Backup strategy implemented

### Security Monitoring

Monitor for:
- Failed authentication attempts
- Rate limit violations
- Unusual API usage patterns
- Unauthorized access attempts
- Database query errors

### Incident Response

If a security issue is discovered:

1. **Assess Impact**: Determine what data/systems are affected
2. **Contain**: Revoke compromised credentials, block IPs if needed
3. **Fix**: Apply patches, update configurations
4. **Notify**: Inform affected users if necessary
5. **Review**: Conduct post-mortem and update security measures

## Contact

For security issues, please contact the security team immediately. Do not open public issues for security vulnerabilities.

## Updates

This document should be updated whenever:
- Security fixes are applied
- Authorization rules change
- New authentication methods are added
- Security best practices evolve
