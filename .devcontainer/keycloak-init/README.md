# Keycloak Initialization Files

This directory contains the initialization files for the Keycloak authentication server used in the development container.

## Files

### `init-keycloak-db.sql`
SQL script that runs during PostgreSQL initialization to create:
- A separate database `keycloak_dev` for Keycloak's data
- A dedicated user `keycloak_dev` with access to that database

This ensures Keycloak and the Clubs application use separate databases within the same PostgreSQL instance.

### `clubs-realm.json`
Keycloak realm configuration that is automatically imported on first startup. This includes:

#### Realm Settings
- **Realm Name**: `clubs-dev`
- **SSL**: Not required (dev mode)
- **Registration**: Enabled
- **Login with email**: Enabled
- **Remember me**: Enabled

#### OAuth/OIDC Client
- **Client ID**: `clubs-frontend`
- **Type**: Public client (no client secret required)
- **Protocol**: OpenID Connect
- **Redirect URIs**: 
  - `http://localhost:5173/*`
  - `http://localhost:5173/auth/callback`
  - `http://localhost:5173/auth/silent-callback`
- **Web Origins**: `http://localhost:5173`
- **PKCE**: S256 (required for public clients)

#### Protocol Mappers
The client includes standard OIDC mappers for:
- Email address
- Given name (first name)
- Family name (last name)
- Full name

#### Pre-configured Users
Two test users are created automatically:

1. **Standard User**
   - Username: `testuser`
   - Email: `testuser@example.com`
   - Password: `testpass`
   - Roles: `user`

2. **Admin User**
   - Username: `admin`
   - Email: `admin@example.com`
   - Password: `admin`
   - Roles: `user`, `admin`

## Modifying the Configuration

To modify the Keycloak configuration:

1. Make changes to `clubs-realm.json`
2. Rebuild the dev container or restart the Keycloak service
3. The realm will be re-imported on startup

Alternatively, you can:
1. Access the Keycloak Admin Console at `http://localhost:8081/admin`
2. Login with username `admin` and password `admin`
3. Make changes through the UI
4. Export the realm to replace `clubs-realm.json` if you want to persist the changes

## Testing Authentication

To test the authentication flow:

1. Start the dev container
2. Wait for Keycloak to fully initialize (check logs with `docker-compose logs keycloak`)
3. Start the backend (`cd Backend && go run main.go`)
4. Start the frontend (`cd Frontend && npm run dev`)
5. Navigate to `http://localhost:5173/login`
6. Click "Sign in with Keycloak" (or similar)
7. Login with `testuser` / `testpass` or `admin` / `admin`
8. You should be redirected back to the application with an authenticated session

## Security Note

⚠️ **These credentials and configurations are for development only!**

Never use these settings in production:
- Change all passwords
- Use proper SSL/TLS
- Configure proper redirect URIs
- Review and adjust all security settings
- Consider using a separate Keycloak instance
