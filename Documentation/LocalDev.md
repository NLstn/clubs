<div align="center">
  <img src="assets/logo.png" alt="Clubs Logo" width="150"/>
  
  # Local Development Guide
  
  **Setting up your local development environment**
</div>

---

## 📋 Overview

For local development you need to replicate the 3 components of clubs in your local environment:

1. **Database** (PostgreSQL)
2. **Backend** (Go application)
3. **Frontend** (React application)

## 🐳 Recommended: Dev Container Setup

The easiest way to get started is using the provided Dev Container configuration:

1. **Prerequisites**: Docker and VS Code with Dev Containers extension
2. **Open the project** in VS Code
3. **Select "Reopen in Container"** when prompted
4. **Wait for the container to build** - all dependencies will be automatically set up
5. **Start developing** - database, Keycloak, and all tools are ready to use

### Default Credentials (Dev Container)

**Test User** (for application login):
- Username: `testuser`
- Password: `testpass`

**Keycloak Admin Console** (http://localhost:8081/admin):
- Username: `admin`
- Password: `admin`

**PostgreSQL Database**:
- Host: `db` (use this in devcontainer)
- Port: `5432`
- Database: `clubs_dev`
- User: `clubs_dev`
- Password: `clubs_dev_password`

---

## 🔧 Manual Setup

If you prefer to set up the environment manually, follow these steps:

## Database

For the database, you have to do the following steps:

1. Create postgresql docker container
```bash
docker run --name postgres -e POSTGRES_PASSWORD=mysecretpassword -p 5432:5432 -d postgres:latest
```

2. Exec into the created container and run psql as user postgres
```bash
docker exec -it -u postgres postgres psql
```

3. Create a new user named clubs
```sql
CREATE USER clubs WITH PASSWORD 'yourpassword';
```

4. Create a new database also named clubs and set the user clubs as its owner, so it can create tables and stuff while automigrating the schema
```sql
CREATE DATABASE clubs OWNER clubs;
```

## Backend

The backend can be run standalone as well if the database is up at the place defined in .env. This file is only relevant in development.

Simply run the backend like this.
```bash
go run main.go
```

## Frontend

The frontend will also work out of the box if the backend is up and running, but you can choose which backend to use. Change the VITE_API_HOST in .env.development to the host of the backend and run

```bash
npm run dev
```

---

## 🎯 VS Code Tasks

The project provides convenient VS Code tasks for development:

- **Start Backend**: Runs the Go backend with hot-reload using Air
- **Start Frontend**: Runs the Vite dev server
- **Start Development Environment**: Runs both backend and frontend in parallel

To use these tasks:
1. Press `Ctrl+Shift+P` (or `Cmd+Shift+P` on Mac)
2. Type "Tasks: Run Task"
3. Select the task you want to run

## 🧪 Running Tests

### Backend Tests
```bash
cd Backend
go test -v ./...
```

### Frontend Tests
```bash
cd Frontend
npm run test
```

## 📦 Building for Production

### Backend
```bash
cd Backend
go build -o clubs-backend
```

### Frontend
```bash
cd Frontend
npm run build
```

## 🆘 Troubleshooting

### Port Already in Use
If you get errors about ports being in use:
- Backend (8080): Stop any other applications using port 8080
- Frontend (5173): Stop any other Vite dev servers
- Keycloak (8081): Stop any other services on port 8081

### Database Connection Issues
Make sure PostgreSQL is running and the credentials in your `.env` file match the database setup.

### Hot Reload Not Working
- Backend: Check that Air is properly installed
- Frontend: Try clearing the Vite cache: `rm -rf Frontend/node_modules/.vite`

---

For more information, see:
- [Backend API Documentation](Backend/API.md)
- [Frontend Design System](Frontend/README.md)