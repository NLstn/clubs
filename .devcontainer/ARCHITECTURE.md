# DevContainer Architecture

This document describes the architecture of the development container setup for the Clubs project.

## Service Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                      Development Container                      │
│                                                                 │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │                    App Container                         │  │
│  │  - Ubuntu base image                                     │  │
│  │  - Go 1.25                                              │  │
│  │  - Node.js 22                                           │  │
│  │  - VS Code workspace at /workspace                      │  │
│  │                                                          │  │
│  │  Network Mode: Shares network with DB container         │  │
│  └──────────────────────────────────────────────────────────┘  │
│                                                                 │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │              PostgreSQL Container (db)                   │  │
│  │  - PostgreSQL 16 Alpine                                  │  │
│  │  - Port: 5432 (exposed)                                  │  │
│  │                                                          │  │
│  │  Databases:                                              │  │
│  │  ├─ clubs_dev (user: clubs_dev)                         │  │
│  │  │  └─ Application database                             │  │
│  │  └─ keycloak_dev (user: keycloak_dev)                   │  │
│  │     └─ Keycloak identity database                       │  │
│  │                                                          │  │
│  │  Volume: postgres-data (persistent)                      │  │
│  └──────────────────────────────────────────────────────────┘  │
│                                                                 │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │           Keycloak Container (keycloak)                  │  │
│  │  - Keycloak 23.0                                         │  │
│  │  - Port: 8081:8080 (mapped)                             │  │
│  │  - Mode: Development (start-dev)                        │  │
│  │                                                          │  │
│  │  Configuration:                                          │  │
│  │  ├─ Realm: clubs-dev (auto-imported)                    │  │
│  │  ├─ Client: clubs-frontend                              │  │
│  │  └─ Test Users: testuser, admin                         │  │
│  │                                                          │  │
│  │  Connected to: PostgreSQL (keycloak_dev)                │  │
│  └──────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

## Port Mapping

| Service    | Internal Port | Exposed Port | Description              |
|------------|--------------|--------------|--------------------------|
| Backend    | 8080         | 8080         | Clubs API Server         |
| Frontend   | 5173         | 5173         | Vite Dev Server          |
| Keycloak   | 8080         | 8081         | Authentication Server    |
| PostgreSQL | 5432         | 5432         | Database Server          |

## Network Configuration

The `app` container uses `network_mode: service:db`, which means:
- The app container shares the network namespace with the db container
- Services in the app container can access PostgreSQL at `localhost:5432`
- Services in the app container can access Keycloak via the Docker network

## Data Flow

### Authentication Flow

```
┌──────────┐         ┌──────────┐         ┌──────────┐         ┌──────────┐
│ Browser  │ ───1───>│ Frontend │ ───2───>│ Keycloak │         │   DB     │
│ :5173    │         │ React    │         │ :8081    │<───────>│ :5432    │
└──────────┘         └──────────┘         └──────────┘         └──────────┘
     ^                    │                     │
     │                    │                     │
     └────────────4───────┘                     │
                          │                     │
                          └────────3───────────>│
                                               │
                          ┌────────────────────v──────┐
                          │  Backend Go Server :8080   │
                          │  - Validates tokens        │
                          │  - API endpoints           │
                          └────────────┬───────────────┘
                                       │
                                       v
                          ┌────────────────────────────┐
                          │  PostgreSQL (clubs_dev)    │
                          │  - User data               │
                          │  - Application data        │
                          └────────────────────────────┘

1. User clicks "Login with Keycloak"
2. Frontend redirects to Keycloak
3. After login, Keycloak sends auth code to Backend
4. Backend validates and returns JWT to Frontend
```

### Database Separation

```
┌─────────────────────────────────────┐
│     PostgreSQL Server (:5432)       │
│                                     │
│  ┌──────────────────────────────┐  │
│  │  clubs_dev database          │  │
│  │  Owner: clubs_dev            │  │
│  │                              │  │
│  │  Used by:                    │  │
│  │  - Backend API               │  │
│  │                              │  │
│  │  Contains:                   │  │
│  │  - Users                     │  │
│  │  - Clubs                     │  │
│  │  - Events                    │  │
│  │  - Members                   │  │
│  │  - etc.                      │  │
│  └──────────────────────────────┘  │
│                                     │
│  ┌──────────────────────────────┐  │
│  │  keycloak_dev database       │  │
│  │  Owner: keycloak_dev         │  │
│  │                              │  │
│  │  Used by:                    │  │
│  │  - Keycloak server           │  │
│  │                              │  │
│  │  Contains:                   │  │
│  │  - Realm configuration       │  │
│  │  - Clients                   │  │
│  │  - Users (auth)              │  │
│  │  - Sessions                  │  │
│  │  - Tokens                    │  │
│  └──────────────────────────────┘  │
└─────────────────────────────────────┘
```

## Initialization Sequence

When the devcontainer starts:

1. **PostgreSQL Container** starts first
   - Initializes the `clubs_dev` database
   - Runs `init-keycloak-db.sql` to create `keycloak_dev` database and user
   - Healthcheck ensures database is ready

2. **Keycloak Container** starts after PostgreSQL is healthy
   - Connects to PostgreSQL (`keycloak_dev` database)
   - Imports the `clubs-realm.json` configuration
   - Creates admin user (admin/admin)
   - Creates test users (testuser, admin)
   - Starts in development mode (HTTP, no SSL required)

3. **App Container** starts after both services
   - Shares network with db container
   - Can access both PostgreSQL and Keycloak
   - Go and Node.js environments ready
   - Workspace mounted at `/workspace`

## Environment Variables

### Backend (.env)
```env
# Database
DATABASE_URL=localhost
DATABASE_PORT=5432
DATABASE_USER=clubs_dev
DATABASE_USER_PASSWORD=clubs_dev_password
DATABASE_NAME=clubs_dev

# Keycloak
KEYCLOAK_SERVER_URL=http://localhost:8081
KEYCLOAK_REALM=clubs-dev
KEYCLOAK_CLIENT_ID=clubs-frontend

# Application
FRONTEND_URL=http://localhost:5173
JWT_SECRET=<your-secret>
```

### Frontend (.env)
```env
# Backend API
VITE_API_HOST=http://localhost:8080

# Keycloak
VITE_KEYCLOAK_URL=http://localhost:8081/realms/clubs-dev
VITE_KEYCLOAK_CLIENT_ID=clubs-frontend
```

## Security Considerations

### Development Mode
- **HTTP only** (no SSL/TLS)
- **Simple passwords** (admin/admin, testpass)
- **Permissive CORS** settings
- **All services exposed** on localhost

### Why This is OK for Development
- All services run locally
- No external network exposure
- Easy to set up and use
- Fast iteration

### What Changes for Production
1. ✅ Enable HTTPS/TLS everywhere
2. ✅ Strong passwords and secrets
3. ✅ Restricted CORS policies
4. ✅ Proper SSL certificates
5. ✅ Network segmentation
6. ✅ Production-grade Keycloak (not dev mode)
7. ✅ Managed PostgreSQL service
8. ✅ Environment-specific configurations
9. ✅ Proper logging and monitoring
10. ✅ Regular security updates

## Volume Persistence

```
postgres-data volume
├─ clubs_dev database files
└─ keycloak_dev database files
```

**Persistence behavior:**
- Data persists across container restarts
- Data persists across devcontainer rebuilds
- Data is lost if volume is deleted (`docker compose down -v`)

## Troubleshooting

### Service Dependencies

If services don't start in the correct order:
```bash
# Check service health
docker compose ps

# View logs for a specific service
docker compose logs keycloak
docker compose logs db

# Restart services in order
docker compose restart db
docker compose restart keycloak
```

### Network Issues

If services can't communicate:
```bash
# Check network configuration
docker network ls
docker network inspect devcontainer_default

# Test connectivity from app container
docker exec -it devcontainer-app-1 curl http://localhost:8081
```

### Database Connection Issues

```bash
# Connect to PostgreSQL
docker exec -it devcontainer-db-1 psql -U clubs_dev -d clubs_dev

# List all databases
docker exec -it devcontainer-db-1 psql -U postgres -c '\l'
```

## Extending the Setup

### Adding a New Service

1. Add service to `docker-compose.yml`
2. Configure dependencies
3. Add port forwarding to `devcontainer.json`
4. Update environment files
5. Document in README

### Modifying Keycloak Configuration

1. Edit `clubs-realm.json`
2. Restart the devcontainer
3. Or manually import via Admin Console

### Database Migrations

The backend automatically runs migrations on startup using GORM AutoMigrate.

## References

- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Dev Containers Specification](https://containers.dev/)
- [Keycloak Documentation](https://www.keycloak.org/documentation)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
