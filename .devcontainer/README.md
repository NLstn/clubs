# Development Container (DevContainer) Setup

This repository includes a DevContainer configuration that provides a complete development environment for the Clubs application, including both frontend and backend development tools, as well as a PostgreSQL database.

## What's Included

### Development Environment
- **Go 1.24**: Backend development with Go modules
- **Node.js LTS**: Frontend development with npm
- **PostgreSQL 15**: Database with pre-configured clubs user and database

### VS Code Extensions
- **Go development**: Official Go extension with language server
- **Frontend development**: TypeScript, ESLint, Prettier, React snippets
- **Database tools**: PostgreSQL extension for database management
- **General development**: Git tools, Docker support

### Pre-configured Services
- **Backend API**: Available on port 8080
- **Frontend Dev Server**: Available on port 5173 (Vite)
- **PostgreSQL Database**: Available on port 5432
  - Username: `clubs`
  - Password: `yourpassword`
  - Database: `clubs`

## Getting Started

### Prerequisites
- [Docker](https://www.docker.com/get-started) installed and running
- [VS Code](https://code.visualstudio.com/) with the [Remote - Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)

### Opening the DevContainer

1. Clone this repository
2. Open the repository folder in VS Code
3. When prompted, click "Reopen in Container" or:
   - Press `F1` to open the command palette
   - Type "Remote-Containers: Reopen in Container"
   - Select it from the list

VS Code will build the development environment and install all dependencies automatically.

### First Time Setup

After the container starts:

1. The database will be automatically initialized with the `clubs` user and database
2. Backend dependencies will be downloaded (`go mod download`)
3. Frontend dependencies will be installed (`npm install`)

### Development Commands

#### Backend Development
```bash
# Navigate to backend directory
cd Backend

# Run the backend server
go run main.go

# Or use Air for hot reloading (if installed)
air

# Run tests
go test ./...

# Build the application
go build -o main .
```

#### Frontend Development
```bash
# Navigate to frontend directory
cd Frontend

# Start development server
npm run dev

# Run tests
npm test

# Build for production
npm run build

# Run linting
npm run lint
```

#### Database Access
The PostgreSQL database is automatically configured and accessible:
- **Host**: `postgres` (within container) or `localhost` (from host)
- **Port**: `5432`
- **Database**: `clubs`
- **Username**: `clubs`
- **Password**: `yourpassword`

You can connect using the PostgreSQL extension in VS Code or any database client.

## Environment Variables

The devcontainer automatically sets up the following environment variables:
- `DATABASE_URL=postgres`
- `DATABASE_PORT=5432`
- `DATABASE_USER=clubs`
- `DATABASE_USER_PASSWORD=yourpassword`
- `FRONTEND_URL=http://localhost:5173`

Additional environment variables can be added to `.devcontainer/.env`.

## Customization

### Adding VS Code Extensions
Edit `.devcontainer/devcontainer.json` and add extension IDs to the `extensions` array.

### Modifying the Development Environment
Edit `.devcontainer/Dockerfile` to add additional tools or modify the base environment.

### Database Configuration
Modify `.devcontainer/init-db.sh` to change the database initialization script.

## Troubleshooting

### Container Won't Start
- Ensure Docker is running
- Try rebuilding the container: Command Palette → "Remote-Containers: Rebuild Container"

### Database Connection Issues
- Verify the PostgreSQL service is running: `docker-compose ps`
- Check the database logs: `docker-compose logs postgres`

### Port Conflicts
If you have local services using ports 5173, 8080, or 5432, either:
- Stop the local services
- Modify the port mappings in `.devcontainer/docker-compose.yml`