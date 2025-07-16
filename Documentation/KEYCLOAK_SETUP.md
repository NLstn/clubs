# Keycloak Authentication Setup

This document describes how to set up and use Keycloak authentication in the Clubs application.

## Overview

The application now supports two authentication methods:
1. **Magic Link Authentication** - Email-based passwordless authentication (existing)
2. **Keycloak OIDC Authentication** - Single Sign-On (SSO) authentication via Keycloak

## Keycloak Configuration

### Server Setup

Your Keycloak instance is configured at: `https://auth.clubsstaging.dev`
- **Realm**: `clubs`
- **Client ID**: `clubs-frontend`

### Environment Variables

#### Backend (.env)
```bash
KEYCLOAK_SERVER_URL=https://auth.clubsstaging.dev
KEYCLOAK_REALM=clubs
KEYCLOAK_CLIENT_ID=clubs-frontend
KEYCLOAK_CLIENT_SECRET=  # Optional for public clients
KEYCLOAK_REDIRECT_URL=http://localhost:5173/auth/callback
```

#### Frontend (.env)
```bash
VITE_KEYCLOAK_URL=https://auth.clubsstaging.dev/realms/clubs
VITE_KEYCLOAK_CLIENT_ID=clubs-frontend
```

## How It Works

### Authentication Flow

1. **User clicks "Login with Keycloak"** on the login page
2. **Redirect to Keycloak** - User is redirected to Keycloak login page
3. **User authenticates** with Keycloak (username/password, social login, etc.)
4. **Authorization code exchange** - Keycloak redirects back with authorization code
5. **Token validation** - Backend exchanges code for tokens and validates with Keycloak
6. **User creation/update** - User is created or updated in local database
7. **JWT tokens issued** - Application issues its own JWT tokens for API access

### Database Integration

- Users authenticated via Keycloak are stored in the same `users` table
- A new `keycloak_id` field stores the Keycloak subject ID
- Email-based matching links existing users with Keycloak accounts
- Full name is automatically populated from Keycloak profile

### API Endpoints

#### Keycloak Authentication
- `GET /api/v1/auth/keycloak/login` - Get Keycloak authorization URL
- `POST /api/v1/auth/keycloak/callback` - Handle Keycloak callback
- `POST /api/v1/auth/keycloak/refresh` - Refresh tokens

#### Existing Magic Link Authentication
- `POST /api/v1/auth/requestMagicLink` - Request magic link
- `GET /api/v1/auth/verifyMagicLink` - Verify magic link
- `POST /api/v1/auth/refreshToken` - Refresh tokens
- `POST /api/v1/auth/logout` - Logout

## Frontend Components

### New Components
- `KeycloakCallback` - Handles the OAuth callback from Keycloak
- `keycloakService` - Service class for Keycloak OIDC operations

### Updated Components
- `Login` - Now includes "Login with Keycloak" button
- `ProtectedRoute` - Enhanced to preserve redirect paths
- `AuthProvider` - Works seamlessly with both authentication methods

## Security Features

- **CSRF Protection** - State parameter validation
- **Token Validation** - OIDC ID token verification
- **Automatic Renewal** - Silent token refresh
- **Secure Storage** - Tokens stored in localStorage with HttpOnly options where available

## Development Setup

1. **Start Backend**:
   ```bash
   cd Backend
   go run main.go
   ```

2. **Start Frontend**:
   ```bash
   cd Frontend
   npm run dev
   ```

3. **Configure Environment**: Copy `.env.example` to `.env` in both directories

## Production Considerations

1. **HTTPS Required** - Keycloak requires HTTPS in production
2. **Client Secrets** - Use client secrets for confidential clients
3. **Token Security** - Consider using secure cookies instead of localStorage
4. **CORS Configuration** - Ensure proper CORS settings for your domains

## Troubleshooting

### Common Issues

1. **CORS Errors**
   - Ensure Keycloak client has correct redirect URIs configured
   - Check CORS settings in both Keycloak and backend

2. **Token Verification Failed**
   - Verify Keycloak server URL and realm are correct
   - Check network connectivity to Keycloak server

3. **User Not Created**
   - Check backend logs for database errors
   - Verify user email is available in Keycloak token

### Debug Information

Enable debug logging by setting environment variables:
```bash
# Backend
LOG_LEVEL=debug

# Frontend
VITE_DEBUG=true
```

## Migration from Magic Link Only

Existing users can seamlessly transition to Keycloak:

1. **Email Matching** - If a user with the same email exists, their account is linked to Keycloak
2. **Profile Completion** - User profile information is updated with Keycloak data
3. **Backwards Compatibility** - Magic link authentication continues to work alongside Keycloak

## Keycloak Client Configuration

In your Keycloak admin console, ensure the `clubs-frontend` client has:

- **Client Type**: Public (or Confidential if using client secret)
- **Valid Redirect URIs**: 
  - `http://localhost:5173/auth/callback` (development)
  - `https://yourdomain.com/auth/callback` (production)
- **Valid Post Logout Redirect URIs**:
  - `http://localhost:5173/login` (development)
  - `https://yourdomain.com/login` (production)
- **Web Origins**: `+` (to allow all valid redirect URIs)

## API Documentation

For detailed API documentation, see [API.md](Documentation/API.md).
