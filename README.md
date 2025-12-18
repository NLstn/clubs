# Clubs

A full-stack club management application for organizations to manage their clubs, members, events, fines, shifts, and more. Built with Go, PostgreSQL, React, and TypeScript.

## Features

- **Club Management**: Create and manage clubs with customizable settings
- **Member Management**: Handle members with role-based permissions (owner, admin, member)
- **Event Scheduling**: Create events including support for recurring events
- **RSVP System**: Track event attendance with going/maybe/not going responses
- **Fine Management**: Customizable fine templates and tracking
- **Shift Scheduling**: Organize shifts with member assignments
- **Team Organization**: Create teams within clubs
- **News & Notifications**: Keep members informed with in-app and email notifications
- **Privacy Controls**: Granular privacy settings for users and members
- **Activity Timeline**: Track club activities and changes
- **Invitation System**: Invite users via email or join request workflow
- **Authentication**: Multiple auth methods (Magic Link, Keycloak OAuth/OIDC)
- **API**: RESTful v1 API and OData v2 API with full CRUD operations

## Technology Stack

- **Backend**: Go 1.25, PostgreSQL, GORM ORM, Azure services
- **Frontend**: React 19, TypeScript, Vite, React Router v7
- **Authentication**: OAuth2/OIDC via Keycloak, JWT tokens, Magic Link
- **Cloud**: Azure Blob Storage, Azure Communication Services

## Quick Start

The easiest way to get started is using the included Dev Container:

1. Open the project in VS Code
2. Select "Reopen in Container" when prompted
3. Wait for the container to build (includes Go, Node.js, PostgreSQL, and Keycloak)
4. Use VS Code tasks to start the development environment

See [Documentation/LocalDev.md](Documentation/LocalDev.md) for manual setup instructions.

## Documentation

- [Local Development Setup](Documentation/LocalDev.md)
- [Backend API Documentation](Documentation/Backend/API.md)
- [Backend - Adding New Tables](Documentation/Backend/AddNewTable.md)
- [Frontend Design System](Documentation/Frontend/README.md)
- [Component Documentation](Documentation/Frontend/components/)

## Default Credentials

**Keycloak Test User**:
- Username: `testuser`
- Password: `testpass`

**Keycloak Admin Console** (http://localhost:8081/admin):
- Username: `admin`
- Password: `admin`

**PostgreSQL Database**:
- Host: `db` (in devcontainer) or `localhost` (host machine)
- Port: `5432`
- Database: `clubs_dev`
- User: `clubs_dev`
- Password: `clubs_dev_password`

## Development

Start the development environment using VS Code tasks:
- **Start Backend**: Runs Go backend with hot-reload
- **Start Frontend**: Runs Vite dev server
- **Start Development Environment**: Runs both in parallel

## Quality Checks

**Frontend**:
```bash
cd Frontend
npm run lint
npm run build
npm run test
```

**Backend**:
```bash
cd Backend
go mod verify
go build -v ./...
go test -v -race ./...
```

## License

[Add your license here]

