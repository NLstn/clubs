# Development Container

This directory contains the configuration for the VS Code Development Container for the Clubs project.

## What's Included

- **Go 1.25**: Backend development environment
- **Node.js 22**: Frontend development environment
- **PostgreSQL 16**: Database server for local development
- **Keycloak 23.0**: Authentication and identity management server

## Database Configuration

The devcontainer includes a PostgreSQL database with the following default credentials:

### Clubs Application Database
- **Host**: `db` (Docker service name - use this in devcontainer)
- **Port**: `5432` (internal only, not exposed to host for security)
- **Database**: `clubs_dev`
- **User**: `clubs_dev`
- **Password**: `clubs_dev_password`

### Keycloak Database
- **Host**: `db` (Docker service name)
- **Port**: `5432` (internal only)
- **Database**: `keycloak_dev`
- **User**: `keycloak_dev`
- **Password**: `keycloak_dev_password`

**Important**: The database is only accessible within the Docker network (not from the host machine). This provides better security in GitHub Codespaces environments. Use the service name `db` instead of `localhost` when connecting from containers.

## Keycloak Configuration

The devcontainer includes a pre-configured Keycloak instance for authentication:

- **URL**: `http://localhost:8081` (accessible from both host and containers)
- **Admin Console**: `http://localhost:8081/admin`
- **Admin Username**: `admin`
- **Admin Password**: `admin`
- **Realm**: `clubs-dev`
- **Frontend Client ID**: `clubs-frontend`

### Port Forwarding Setup

To ensure Keycloak works correctly with OAuth2/OIDC:
- Keycloak runs on port `8080` inside its container
- Port `8081` on the host is mapped to Keycloak's port `8080`
- The app container uses `socat` to forward `localhost:8081` to `keycloak:8080`

This setup ensures both the host browser and the backend application can access Keycloak at `http://localhost:8081`, which is required for proper OAuth2/OIDC validation.

### PKCE Authentication Flow

The Keycloak integration uses PKCE (Proof Key for Code Exchange) for secure authentication:
- The backend generates a cryptographic code verifier and challenge
- The code challenge is sent to Keycloak during authorization
- The code verifier is stored in the browser's sessionStorage
- During callback, the verifier is sent to the backend for token exchange

This provides additional security for the OAuth2 authorization code flow.

### Pre-configured Test Users

The realm is configured to allow user registration. You can:
- **Register a new user** via the Keycloak login page
- **Use the admin account** to manage users and realm settings
  - Username: `admin`
  - Password: `admin`
  - Access: `http://localhost:8081/admin`

## Getting Started

1. Install the [Dev Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) in VS Code
2. Open the project folder in VS Code
3. When prompted, click "Reopen in Container" or use the command palette (F1) and select "Dev Containers: Reopen in Container"
4. Wait for the container to build and start (first startup may take 2-3 minutes as Keycloak initializes)
5. Copy `Backend/.env.example` to `Backend/.env` to use the devcontainer credentials
6. Copy `Frontend/.env.example` to `Frontend/.env` to use the devcontainer Keycloak configuration
7. The services will be available at:
   - Backend API: `http://localhost:8080`
   - Frontend: `http://localhost:5173`
   - Keycloak: `http://localhost:8081`

### First Time Setup

After the devcontainer starts:

1. **Verify Keycloak is running**: Navigate to `http://localhost:8081` - you should see the Keycloak welcome page
2. **Test authentication**: 
   - Go to `http://localhost:5173`
   - Click "Login with Keycloak"
   - Register a new user account
   - You should be redirected back to the application after successful registration

## Included VS Code Extensions

- Go language support
- ESLint for JavaScript/TypeScript linting
- Prettier for code formatting
- Docker tools
- PostgreSQL client for database management

## Data Persistence

The PostgreSQL data (including both the clubs and keycloak databases) is stored in a Docker volume (`postgres-data`), so your databases will persist between container rebuilds.

## Troubleshooting

### Keycloak not starting

If Keycloak fails to start or shows connection errors:
1. Check that PostgreSQL is healthy: `docker compose ps`
2. View Keycloak logs: `docker compose logs keycloak`
3. Ensure the Keycloak database was created: Connect to PostgreSQL and run `\l` to list databases

### Keycloak realm not imported

If the clubs-dev realm is not available:
1. The realm should import automatically on first startup
2. You can manually import it via the Admin Console at `http://localhost:8081/admin`
3. Navigate to the realm dropdown (top left) and select "Create Realm"
4. Import the file from `.devcontainer/keycloak-init/clubs-realm.json`

### Authentication not working

1. Verify environment variables in `Backend/.env` match:
   - `KEYCLOAK_SERVER_URL=http://localhost:8081`
   - `KEYCLOAK_REALM=clubs-dev`
   - `KEYCLOAK_CLIENT_ID=clubs-frontend`
2. Ensure Keycloak is accessible at `http://localhost:8081`
3. Check that the clubs-dev realm is properly imported
4. Clear browser cache and sessionStorage if experiencing issues
5. Check browser console for detailed error messages

### Port forwarding issues

If you get "Connection refused" errors:
1. Verify the app container's socat process is running: `ps aux | grep socat`
2. Check that Keycloak is healthy: `docker compose ps keycloak`
3. Restart the devcontainer if port forwarding is not working
