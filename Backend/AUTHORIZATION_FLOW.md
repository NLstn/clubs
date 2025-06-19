# Authorization Flow Documentation

This document describes the HTTP cookie-based authentication system implemented in the Clubs application.

## Overview

The application uses secure HTTP cookies for authentication instead of localStorage or Authorization headers. This provides better security by preventing XSS attacks from accessing tokens.

## Authentication Flow

### 1. Magic Link Request
- **Endpoint**: `POST /api/v1/auth/requestMagicLink`
- **Input**: `{"email": "user@example.com"}`
- **Response**: `204 No Content`
- **Action**: Sends a magic link email to the user

### 2. Magic Link Verification
- **Endpoint**: `GET /api/v1/auth/verifyMagicLink?token=<magic_token>`
- **Response**: `204 No Content`
- **Cookies Set**:
  - `access_token`: JWT token valid for 15 minutes
  - `refresh_token`: JWT token valid for 30 days
- **Cookie Attributes**:
  - `HttpOnly`: Prevents JavaScript access
  - `SameSite=Strict`: CSRF protection
  - `Secure`: Set automatically for HTTPS environments
  - `Path=/`: Available to all routes

### 3. Token Refresh
- **Endpoint**: `POST /api/v1/auth/refreshToken`
- **Input**: Refresh token from cookie (automatic)
- **Response**: `204 No Content`
- **Cookies Updated**:
  - `access_token`: New JWT token valid for 15 minutes

### 4. Logout
- **Endpoint**: `POST /api/v1/auth/logout`
- **Input**: Refresh token from cookie (automatic)
- **Response**: `204 No Content`
- **Action**: 
  - Deletes all refresh tokens from database
  - Clears cookies by setting `MaxAge=-1`

## Authentication Middleware

### Token Validation
- **Endpoint Protection**: Applied to protected routes via `AuthMiddleware`
- **Token Source**: Only reads from `access_token` cookie
- **Validation**: JWT signature verification using HMAC-SHA256
- **Context**: Adds `userID` to request context for handlers

### Request Flow
1. Middleware extracts `access_token` from cookie
2. Validates JWT signature and expiration
3. Extracts `user_id` claim from token
4. Adds `userID` to request context
5. Continues to protected handler

## Security Features

### Cookie Security
- **HttpOnly**: Prevents XSS attacks by blocking JavaScript access
- **SameSite=Strict**: Prevents CSRF attacks by restricting cross-site requests
- **Secure Flag**: Automatically enabled for HTTPS environments
- **Path Restriction**: Cookies scoped to root path only

### Token Security
- **Short-lived Access Tokens**: 15-minute expiration reduces exposure window
- **Long-lived Refresh Tokens**: 30-day expiration for user convenience
- **Database Storage**: Refresh tokens stored in database for revocation
- **Automatic Refresh**: Frontend handles token refresh transparently

## Frontend Integration

### Automatic Cookie Handling
- **Axios Configuration**: `withCredentials: true` enables automatic cookie transmission
- **No Manual Headers**: No need to set Authorization headers
- **Token Refresh**: Interceptor handles expired tokens automatically

### Authentication State
- **Cookie Detection**: Frontend checks for `access_token` cookie presence
- **State Management**: React context tracks authentication status
- **Auto-logout**: Redirects to login when tokens are invalid or missing

## Error Handling

### Common Responses
- `401 Unauthorized`: Missing or invalid token
- `400 Bad Request`: Malformed request
- `404 Not Found`: User not found
- `500 Internal Server Error`: Server-side errors

### Token Expiration
- **Access Token**: Automatically refreshed by frontend interceptor
- **Refresh Token**: Requires new magic link authentication
- **Invalid Tokens**: Automatic redirect to login page

## Environment Configuration

### Required Variables
- `FRONTEND_URL`: Used to determine HTTPS and set Secure cookie flag
- Cookie security attributes automatically configured based on environment

### Development vs Production
- **Development**: Cookies work over HTTP (Secure=false)
- **Production**: Cookies require HTTPS (Secure=true)
- **Detection**: Based on `FRONTEND_URL` starting with "https://"

## Migration Notes

This system is **cookie-only** and does not support backwards compatibility with Authorization headers or localStorage. All authentication must use HTTP cookies.

## Testing

### Unit Tests
- Cookie setting and clearing verification
- Token validation from cookies only
- Security attribute verification (HttpOnly, SameSite, Secure)

### Integration Tests
- End-to-end authentication flow
- Token refresh with cookies
- Logout cookie clearing
- Protected endpoint access

## Troubleshooting

### Common Issues
1. **Cross-origin requests**: Ensure `withCredentials: true` or `credentials: 'include'`
2. **Missing cookies**: Check that cookies are being set by server
3. **HTTPS requirements**: Secure flag requires HTTPS in production
4. **SameSite restrictions**: Cookies only sent with same-site requests

### Debug Tips
- Check browser dev tools Application/Storage tab for cookies
- Verify cookie attributes match security requirements
- Ensure CORS is configured for credentials
- Check server logs for token validation errors