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
- **Host**: `localhost`
- **Port**: `5432`
- **Database**: `clubs_dev`
- **User**: `clubs_dev`
- **Password**: `clubs_dev_password`

### Keycloak Database
- **Host**: `localhost`
- **Port**: `5432`
- **Database**: `keycloak_dev`
- **User**: `keycloak_dev`
- **Password**: `keycloak_dev_password`

These credentials are already configured in the `Backend/.env.example` file.

## Keycloak Configuration

The devcontainer includes a pre-configured Keycloak instance for authentication:

- **URL**: `http://localhost:8081`
- **Admin Console**: `http://localhost:8081/admin`
- **Admin Username**: `admin`
- **Admin Password**: `admin`
- **Realm**: `clubs-dev`
- **Frontend Client ID**: `clubs-frontend`

### Pre-configured Test Users

Two test users are automatically created:

1. **Standard User**
   - Username: `testuser`
   - Password: `testpass`
   - Email: `testuser@example.com`

2. **Admin User**
   - Username: `admin`
   - Password: `admin`
   - Email: `admin@example.com`

## Getting Started

1. Install the [Dev Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) in VS Code
2. Open the project folder in VS Code
3. When prompted, click "Reopen in Container" or use the command palette (F1) and select "Dev Containers: Reopen in Container"
4. Wait for the container to build and start (Keycloak may take 1-2 minutes to fully initialize)
5. Copy `Backend/.env.example` to `Backend/.env` to use the devcontainer credentials
6. Copy `Frontend/.env.example` to `Frontend/.env` to use the devcontainer Keycloak configuration
7. The services will be available at:
   - Backend API: `http://localhost:8080`
   - Frontend: `http://localhost:5173`
   - Keycloak: `http://localhost:8081`
   - PostgreSQL: `localhost:5432`

### Testing the Setup

To verify that Keycloak is properly configured, run the test script:

```bash
.devcontainer/test-keycloak.sh
```

This will check:
- Keycloak accessibility
- Realm configuration
- OIDC endpoints
- Database connections
- Test user authentication

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

1. Verify environment variables in `Backend/.env` and `Frontend/.env` match the values in `.env.example`
2. Ensure Keycloak is accessible at `http://localhost:8081`
3. Check that the frontend redirect URIs are configured correctly in Keycloak
4. Clear browser cache and local storage if experiencing issues
