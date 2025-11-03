# Development Container

This directory contains the configuration for the VS Code Development Container for the Clubs project.

## What's Included

- **Go 1.25**: Backend development environment
- **Node.js 20**: Frontend development environment
- **PostgreSQL 16**: Database server for local development

## Database Configuration

The devcontainer includes a PostgreSQL database with the following default credentials:

- **Host**: `localhost`
- **Port**: `5432`
- **Database**: `clubs_dev`
- **User**: `clubs_dev`
- **Password**: `clubs_dev_password`

These credentials are already configured in the `Backend/.env.example` file.

## Getting Started

1. Install the [Dev Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) in VS Code
2. Open the project folder in VS Code
3. When prompted, click "Reopen in Container" or use the command palette (F1) and select "Dev Containers: Reopen in Container"
4. Wait for the container to build and start
5. Copy `Backend/.env.example` to `Backend/.env` to use the devcontainer database credentials
6. The backend will be available on port 8080, frontend on port 5173, and PostgreSQL on port 5432

## Included VS Code Extensions

- Go language support
- ESLint for JavaScript/TypeScript linting
- Prettier for code formatting
- Docker tools
- PostgreSQL client for database management

## Data Persistence

The PostgreSQL data is stored in a Docker volume (`postgres-data`), so your database will persist between container rebuilds.
