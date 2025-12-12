# Keycloak Setup Checklist

Use this checklist to ensure your Keycloak development environment is properly configured.

## Initial Setup

- [ ] **Open in Dev Container**
  - Open project in VS Code
  - Click "Reopen in Container" when prompted
  - Wait 2-3 minutes for all services to start

- [ ] **Copy Environment Files**
  ```bash
  cp Backend/.env.example Backend/.env
  cp Frontend/.env.example Frontend/.env
  ```

- [ ] **Verify Keycloak is Running**
  - Open `http://localhost:8081` in your browser
  - You should see the Keycloak welcome page
  - Or run: `.devcontainer/test-keycloak.sh`

## Testing Authentication

- [ ] **Start Backend**
  ```bash
  cd Backend
  go run main.go
  ```
  - Should see: "Server starting on port 8080"
  - Should see: "Keycloak initialized successfully" (or warning if not configured)

- [ ] **Start Frontend**
  ```bash
  cd Frontend
  npm install  # First time only
  npm run dev
  ```
  - Should see: "Local: http://localhost:5173"

- [ ] **Test Login Flow**
  - Open `http://localhost:5173` in browser
  - Navigate to login page
  - Click "Sign in with Keycloak" (or similar button)
  - Should redirect to Keycloak login at `http://localhost:8081`

- [ ] **Login with Test User**
  - Username: `testuser`
  - Password: `testpass`
  - Should redirect back to application
  - Should be logged in successfully

## Verification

- [ ] **Check Logs** - No errors related to Keycloak or authentication

- [ ] **Check Token** - In browser console:
  ```javascript
  localStorage.getItem('auth_token')
  // Should return a JWT token
  ```

- [ ] **Access Protected Routes** - Try accessing protected pages/features

- [ ] **Logout** - Test logout functionality works correctly

## Admin Console Access

- [ ] **Access Admin Console**
  - Open `http://localhost:8081/admin`
  - Username: `admin`
  - Password: `admin`
  - Should see Keycloak admin dashboard

- [ ] **Verify Realm**
  - Realm dropdown (top left) should show "clubs-dev"
  - Click on "clubs-dev" to view realm settings

- [ ] **Verify Client**
  - Go to "Clients" in left menu
  - Should see "clubs-frontend" client
  - Client should be enabled

- [ ] **Verify Users**
  - Go to "Users" in left menu
  - Should see at least 2 users: "testuser" and "admin"

## Troubleshooting

If any step fails, check:

- [ ] **Keycloak container is running**
  ```bash
  docker ps | grep keycloak
  ```

- [ ] **Keycloak logs**
  ```bash
  docker compose -f .devcontainer/docker-compose.yml logs keycloak
  ```

- [ ] **PostgreSQL databases exist**
  ```bash
  PGPASSWORD=clubs_dev_password psql -U clubs_dev -h localhost -l | grep -E "clubs_dev|keycloak_dev"
  ```

- [ ] **Environment variables are correct**
  - Check `Backend/.env` matches `Backend/.env.example`
  - Check `Frontend/.env` matches `Frontend/.env.example`

- [ ] **Clear browser cache and local storage**
  - Open DevTools (F12) → Application → Storage → Clear site data

## Common Issues

### "Keycloak not accessible"
- Wait 2 minutes after starting devcontainer
- Check `docker compose logs keycloak`
- Restart: `docker compose restart keycloak`

### "Realm not found"
- Check admin console at `http://localhost:8081/admin`
- Manually import `.devcontainer/keycloak-init/clubs-realm.json`

### "Invalid redirect URI"
- Verify frontend URL in `Backend/.env`: `FRONTEND_URL=http://localhost:5173`
- Check client redirect URIs in Keycloak admin console

### "Token validation failed"
- Check `KEYCLOAK_SERVER_URL` in `Backend/.env`
- Check `VITE_KEYCLOAK_URL` in `Frontend/.env`
- Ensure both point to `http://localhost:8081`

## Need More Help?

Refer to detailed documentation:
- [KEYCLOAK_QUICKSTART.md](KEYCLOAK_QUICKSTART.md) - Complete guide with examples
- [ARCHITECTURE.md](ARCHITECTURE.md) - System architecture and data flows
- [README.md](README.md) - DevContainer overview
- [keycloak-init/README.md](keycloak-init/README.md) - Configuration details

---

✅ **All checks passed?** You're ready to develop with Keycloak authentication!
