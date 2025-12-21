<div align="center">
  <img src="../assets/logo.png" alt="Clubs Logo" width="150"/>
  
  # CSRF Protection
  
  **Cross-Site Request Forgery protection mechanisms**
</div>

---

# CSRF Protection

This document describes the Cross-Site Request Forgery (CSRF) protection mechanisms implemented in the Clubs backend API.

## Overview

The application implements multiple layers of CSRF protection:

1. **OAuth State Token Validation** - HMAC-signed tokens with IP validation for OAuth/Keycloak flows
2. **JSON-Based API** - Natural CSRF protection (browsers won't submit JSON via forms)
3. **JWT Tokens in Headers** - Authentication tokens in Authorization header (not cookies)

## OAuth State Token Protection

### How It Works

The OAuth/Keycloak authentication flow uses cryptographically signed state tokens to prevent CSRF attacks and session fixation vulnerabilities.

#### State Token Format

State tokens are HMAC-SHA256 signed and contain:
- **Nonce**: 32-byte random value (base64url encoded)
- **Timestamp**: Unix timestamp when token was created
- **Client IP Hash**: SHA-256 hash of client IP address
- **Signature**: HMAC-SHA256 signature of nonce + timestamp + IP hash

Format: `nonce.timestamp.signature`

Example: `dGVzdG5vbmNl.1703074800.a1b2c3d4e5f6...`

#### Security Features

1. **Signature Verification**: Tamper-proof using HMAC-SHA256 with server secret
2. **Timestamp Validation**: Tokens expire after 10 minutes
3. **IP Binding**: Token is bound to client IP address hash
4. **One-Time Use**: Nonce is stored in database and consumed on use
5. **Replay Attack Prevention**: Each state token can only be used once

### Implementation

#### Login Endpoint (`GET /api/v1/auth/keycloak/login`)

```json
{
  "authURL": "https://keycloak.example.com/...",
  "state": "nonce.timestamp.signature",
  "codeVerifier": "pkce-verifier"
}
```

The state token is:
1. Generated with HMAC signature including client IP hash
2. Stored in database (`oauth_states` table) for one-time use validation
3. Returned to client for OAuth redirect

#### Callback Endpoint (`POST /api/v1/auth/keycloak/callback`)

```json
{
  "code": "oauth-authorization-code",
  "state": "nonce.timestamp.signature",
  "codeVerifier": "pkce-verifier"
}
```

The callback validates:
1. **State presence**: Required parameter
2. **Signature validation**: HMAC signature must be valid
3. **Timestamp validation**: Token must not be expired (10 min)
4. **IP validation**: Client IP must match the IP used during login
5. **Nonce consumption**: Nonce must exist in database and not be used
6. **One-time use**: Nonce is deleted after successful validation

**Validation Flow:**
```
1. Extract nonce from state token
2. Verify HMAC signature (prevents tampering)
3. Check timestamp (prevents replay after expiry)
4. Validate IP hash matches (prevents CSRF from different network)
5. Check nonce exists in database (prevents fabricated states)
6. Delete nonce (one-time use, prevents replay)
```

### Database Model

The `oauth_states` table tracks state nonces:

```sql
CREATE TABLE oauth_states (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    state TEXT UNIQUE NOT NULL,        -- The nonce value
    ip_hash TEXT NOT NULL,              -- SHA-256 hash of client IP
    expires_at TIMESTAMP NOT NULL,      -- Expiration time (10 minutes)
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_oauth_states_expires_at ON oauth_states(expires_at);
```

### Cleanup

Expired state tokens are automatically cleaned up by calling:
```go
models.CleanupExpiredOAuthStates()
```

This should be called periodically (e.g., via cron job or scheduled task).

## CSRF Protection for Other Endpoints

### Current Protection

All API endpoints use:

1. **JSON Content-Type**: Endpoints require `application/json`, which browsers won't send from forms
2. **JWT in Authorization Header**: Tokens are in headers, not cookies
3. **No Cookie-Based Authentication**: Application doesn't use cookies for auth

These mechanisms provide natural CSRF protection because:
- HTML forms cannot set custom `Content-Type: application/json`
- HTML forms cannot set custom `Authorization` headers
- Browsers enforce Same-Origin Policy on AJAX requests

### Additional Protection (Future)

For endpoints that might need additional CSRF protection:

1. **X-CSRF-Token Header**: Validate CSRF token on state-changing operations
2. **Double Submit Cookie**: If cookies are added in the future
3. **SameSite Cookie Attribute**: If cookies are added in the future

## Security Secrets

### Environment Variables

The following secrets must be configured:

```bash
# Required for JWT tokens and CSRF
JWT_SECRET=your-strong-secret-here

# Optional: Separate CSRF secret (falls back to JWT_SECRET)
CSRF_SECRET=your-csrf-secret-here
```

**Important:**
- Use cryptographically random secrets (at least 256 bits)
- Keep secrets secure and never commit to version control
- Rotate secrets periodically
- Use different secrets for production and development

## API Reference

### CSRF Package (`github.com/NLstn/clubs/csrf`)

#### `Init() error`
Initializes CSRF protection with secrets from environment variables.

#### `GenerateStateToken(ipHash string) (string, error)`
Generates a signed OAuth state token.

Parameters:
- `ipHash`: SHA-256 hash of client IP address

Returns: Signed state token in format `nonce.timestamp.signature`

#### `ValidateStateToken(stateToken string, ipHash string) (string, bool)`
Validates a signed OAuth state token.

Parameters:
- `stateToken`: Token to validate
- `ipHash`: SHA-256 hash of client IP address

Returns:
- `nonce`: The nonce value if valid
- `valid`: Boolean indicating if token is valid

#### `HashIP(ip string) string`
Creates a SHA-256 hash of an IP address.

Parameters:
- `ip`: IP address string

Returns: Hex-encoded SHA-256 hash

### Models Package

#### `CreateOAuthState(state string, ipHash string) error`
Stores an OAuth state nonce for validation.

#### `ValidateAndConsumeOAuthState(state string, ipHash string) (bool, error)`
Validates and consumes (deletes) an OAuth state nonce.

#### `CleanupExpiredOAuthStates() error`
Removes expired OAuth state nonces from the database.

## Testing

The CSRF protection is thoroughly tested:

- `Backend/csrf/csrf_test.go` - Token generation, validation, expiration
- `Backend/models/oauth_state_test.go` - State storage, consumption, cleanup

Run tests:
```bash
cd Backend
go test ./csrf/... -v
go test ./models/... -run OAuth -v
```

## Security Considerations

### Threats Mitigated

✅ **CSRF on OAuth Login** - State token prevents attackers from forcing users to log in with attacker's account
✅ **Session Fixation** - One-time use nonces prevent session fixation attacks  
✅ **Replay Attacks** - Timestamp and nonce consumption prevent replay  
✅ **State Tampering** - HMAC signature prevents tampering with state parameters  
✅ **Cross-Network CSRF** - IP validation prevents CSRF from different networks

### Remaining Considerations

⚠️ **IP Changes**: Users changing networks during OAuth flow will fail validation
  - Mitigation: 10-minute token expiry provides reasonable window
  - Consider: Relaxing IP validation for mobile users if needed

⚠️ **Shared IPs**: Users behind NAT/proxy with same IP
  - Mitigation: Nonce uniqueness still provides protection
  - Impact: Multiple users from same network can each have their own states

⚠️ **X-Forwarded-For Header Trust**: IP validation relies on X-Forwarded-For header
  - **Security Risk**: This header can be spoofed by clients if backend is directly accessible
  - **Mitigation**: Ensure backend is only accessible through a trusted reverse proxy
  - **Best Practice**: Configure proxy to strip/override client-provided X-Forwarded-For headers
  - **Alternative**: Use RemoteAddr only (requires removing X-Forwarded-For logic in getClientIP)
  - **Note**: The `getClientIP()` function includes security warnings in code comments

✅ **No Cookie-Based CSRF**: Application uses JWT in headers, not cookies
  - No SameSite attribute needed
  - No double-submit cookie pattern needed

## References

- [OWASP CSRF Prevention Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html)
- [OAuth 2.0 Security Best Current Practice](https://datatracker.ietf.org/doc/html/draft-ietf-oauth-security-topics)
- [PKCE (RFC 7636)](https://datatracker.ietf.org/doc/html/rfc7636)

## Changelog

### 2025-12-20
- Implemented HMAC-signed OAuth state tokens
- Added server-side state nonce validation
- Added IP address binding for state tokens
- Created comprehensive test coverage
- Documented CSRF protection mechanisms
