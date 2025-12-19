# Authentication

Clubs supports multiple authentication methods to ensure secure access to the platform.

## Authentication Methods

### 1. Keycloak (OAuth/OIDC)

In the development environment, Clubs uses Keycloak for authentication. Keycloak provides secure Single Sign-On (SSO) capabilities.

**Development Test User Credentials:**

- **Username:** `testuser`
- **Password:** `testpass`

**Keycloak Admin Console Access:**

- URL: [http://localhost:8081/admin](http://localhost:8081/admin)
- **Username:** `admin`
- **Password:** `admin`

### 2. Magic Link Email Authentication

Clubs also supports passwordless authentication via Magic Link emails. Users receive a secure link via email that logs them in automatically.

## Login Steps

1. **Start Services**: Ensure the backend and frontend services are running
2. **Open Application**: Navigate to the frontend in your browser
3. **Choose Authentication Method**:
   - Click **Login with Keycloak** for OAuth/OIDC authentication
   - OR enter your email for Magic Link authentication
4. **Authenticate**:
   - For Keycloak: Use the test user credentials above
   - For Magic Link: Check your email for the login link

## First Time Setup

When logging in for the first time:

1. Complete your profile information
2. Create a new club or request to join an existing one
3. Configure your notification preferences

## Security Features

- **JWT Tokens**: Secure access and refresh token mechanism
- **Role-Based Access Control**: Different permissions for admins and members
- **Session Management**: View and manage active sessions from your profile
