# Keycloak Quick Start Guide

This guide will help you get started with the integrated Keycloak authentication in your development environment.

## Automatic Setup

When you open the project in the Dev Container, Keycloak is automatically:
1. ✅ Started with PostgreSQL backend
2. ✅ Configured with the `clubs-dev` realm
3. ✅ Set up with the `clubs-frontend` client
4. ✅ Populated with test users

No manual configuration needed! 🎉

## First Time Setup

### 1. Start the Dev Container

Open the project in VS Code and reopen in the container when prompted.

### 2. Copy Environment Files

```bash
# Backend configuration
cp Backend/.env.example Backend/.env

# Frontend configuration
cp Frontend/.env.example Frontend/.env
```

These files already contain the correct Keycloak configuration for local development.

### 3. Verify Keycloak is Running

Wait 1-2 minutes for Keycloak to fully initialize, then run:

```bash
.devcontainer/test-keycloak.sh
```

You should see all checks pass with green checkmarks.

### 4. Start the Application

**Terminal 1 - Backend:**
```bash
cd Backend
go run main.go
```

**Terminal 2 - Frontend:**
```bash
cd Frontend
npm install  # First time only
npm run dev
```

## Using Keycloak Authentication

### Access the Application

1. Open your browser to `http://localhost:5173`
2. Click on the login/sign-in button
3. You'll be redirected to Keycloak
4. Login with one of the test users:
   - **Standard user**: `testuser` / `testpass`
   - **Admin user**: `admin` / `admin`
5. You'll be redirected back to the application, now authenticated

### Authentication Flow

```
Frontend (localhost:5173)
    ↓
    1. User clicks "Login with Keycloak"
    ↓
Keycloak (localhost:8081)
    ↓
    2. User enters credentials
    ↓
    3. Keycloak issues authorization code
    ↓
Backend (localhost:8080)
    ↓
    4. Backend exchanges code for tokens
    ↓
    5. Backend validates tokens and creates session
    ↓
Frontend
    ↓
    6. User is logged in
```

## Admin Console Access

To manage users, clients, and realm settings:

1. Open `http://localhost:8081/admin`
2. Login with:
   - Username: `admin`
   - Password: `admin`
3. Select the `clubs-dev` realm from the dropdown

### Common Admin Tasks

#### Add a New User

1. Go to Admin Console → clubs-dev realm
2. Click "Users" in the left menu
3. Click "Add user"
4. Fill in the details:
   - Username (required)
   - Email
   - First name
   - Last name
5. Click "Create"
6. Go to the "Credentials" tab
7. Set a password and disable "Temporary"
8. Click "Set password"

#### Modify Redirect URIs

If you need to add more redirect URIs (e.g., for different ports):

1. Go to Admin Console → clubs-dev realm
2. Click "Clients" in the left menu
3. Click "clubs-frontend"
4. Scroll to "Valid redirect URIs"
5. Add your URI (e.g., `http://localhost:3000/*`)
6. Click "Save"

## Troubleshooting

### Keycloak Not Starting

**Check container status:**
```bash
docker ps
```

Look for the `keycloak` container. If it's not running:

```bash
docker compose -f .devcontainer/docker-compose.yml logs keycloak
```

### "Realm not found" Error

The realm import may have failed. Check logs:

```bash
docker compose -f .devcontainer/docker-compose.yml logs keycloak | grep -i import
```

If needed, manually import the realm:
1. Go to `http://localhost:8081/admin`
2. Click the realm dropdown (top left)
3. Click "Create Realm"
4. Click "Browse" and select `.devcontainer/keycloak-init/clubs-realm.json`
5. Click "Create"

### Authentication Fails

**Check environment variables:**

Backend `.env` should have:
```
KEYCLOAK_SERVER_URL=http://localhost:8081
KEYCLOAK_REALM=clubs-dev
KEYCLOAK_CLIENT_ID=clubs-frontend
FRONTEND_URL=http://localhost:5173
```

Frontend `.env` should have:
```
VITE_KEYCLOAK_URL=http://localhost:8081/realms/clubs-dev
VITE_KEYCLOAK_CLIENT_ID=clubs-frontend
```

**Clear browser data:**
1. Open browser DevTools (F12)
2. Go to Application → Storage
3. Click "Clear site data"
4. Try logging in again

### Database Connection Issues

**Verify databases exist:**
```bash
# Check clubs database
psql -U clubs_dev -h localhost -d clubs_dev -c '\l'

# Check keycloak database
PGPASSWORD=keycloak_dev_password psql -U keycloak_dev -h localhost -d keycloak_dev -c '\l'
```

If databases don't exist, you may need to recreate the containers:
```bash
docker compose -f .devcontainer/docker-compose.yml down -v
docker compose -f .devcontainer/docker-compose.yml up -d
```

## Development Tips

### Testing Different Users

You can quickly switch between test users to test different permission levels:
- Use `testuser` for standard user flows
- Use `admin` for administrative features

### Viewing Tokens

In the browser DevTools console:
```javascript
// Get the stored access token
localStorage.getItem('auth_token')

// Get the Keycloak ID token
localStorage.getItem('keycloak_id_token')
```

### Backend Token Validation

The backend validates tokens in the Keycloak middleware. Check logs for validation errors:
```bash
cd Backend
go run main.go
# Look for lines containing "Invalid token" or "Keycloak"
```

### Testing Token Expiration

Tokens expire after 1 hour by default. To test expiration:
1. Login and note the time
2. Wait for the token to expire (or modify the realm settings for shorter expiration)
3. Try to access a protected resource
4. You should be prompted to login again

## Production Considerations

⚠️ **Important**: This configuration is for development only!

Before deploying to production:

1. **Change all passwords** (admin, database, test users)
2. **Enable SSL/TLS** (set `KC_HOSTNAME_STRICT_HTTPS=true`)
3. **Use a proper Keycloak instance** (not dev mode)
4. **Configure proper redirect URIs** (your production domain)
5. **Review security settings** (password policies, session timeouts)
6. **Enable logging and monitoring**
7. **Use a separate Keycloak database server**
8. **Consider Keycloak clustering** for high availability

## Additional Resources

- [Keycloak Documentation](https://www.keycloak.org/documentation)
- [OIDC Specification](https://openid.net/connect/)
- [Keycloak Admin CLI](https://www.keycloak.org/docs/latest/server_admin/#admin-cli)
- [Project Documentation](../Documentation/)

## Need Help?

If you encounter issues not covered here:
1. Check the [troubleshooting section](#troubleshooting)
2. Review the logs (backend, frontend, Keycloak)
3. Check the [.devcontainer/README.md](.devcontainer/README.md)
4. Check the [keycloak-init/README.md](keycloak-init/README.md)
