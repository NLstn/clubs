# AI Agent Instructions for Clubs Project

## Project Overview

**Clubs** is a full-stack club management application designed to help organizations manage their clubs, members, events, fines, shifts, and more. The application provides a comprehensive platform for club administration with features including:

- Club creation and management
- Member management with roles and permissions
- Event scheduling (including recurring events)
- Fine management with customizable templates
- Shift scheduling and management
- Team organization
- News and notifications
- Invitation and join request workflows
- OAuth/OIDC authentication via Keycloak
- Azure integration for storage and communication services

## Technology Stack

### Backend
- **Language:** Go 1.25
- **Web Framework:** Custom HTTP handlers
- **Database:** PostgreSQL with GORM ORM
- **Authentication:** OAuth2/OIDC (Keycloak), JWT tokens, Magic Link email authentication
- **Cloud Services:** Azure (Blob Storage, Communication Services)
- **Testing:** Go testing framework with testify

### Frontend
- **Framework:** React 19 with TypeScript
- **Build Tool:** Vite
- **Routing:** React Router v7
- **State Management:** React Context
- **Internationalization:** i18next
- **Authentication:** oidc-client-ts
- **HTTP Client:** Axios
- **Testing:** Vitest with Testing Library
- **Linting:** ESLint 9

## Project Structure

```
/workspace
├── Backend/           # Go backend application
│   ├── auth/         # Authentication logic (Keycloak, JWT)
│   ├── azure/        # Azure service integrations
│   ├── database/     # Database connection and migrations
│   ├── handlers/     # HTTP request handlers (API routes)
│   ├── models/       # Data models and business logic
│   ├── notifications/# Notification service
│   ├── tools/        # Utility functions
│   ├── main.go       # Application entry point
│   └── go.mod        # Go dependencies
├── Frontend/         # React frontend application
│   ├── src/
│   │   ├── components/ # Reusable UI components
│   │   ├── context/    # React context providers
│   │   ├── hooks/      # Custom React hooks
│   │   ├── i18n/       # Internationalization files
│   │   ├── pages/      # Page components
│   │   └── utils/      # Utility functions
│   ├── package.json  # Node dependencies
│   └── vite.config.ts# Vite configuration
└── Documentation/    # Project documentation
    ├── LocalDev.md   # Local development setup
    ├── Backend/      # Backend-specific docs
    └── Frontend/     # Frontend-specific docs
```

## Documentation Locations

- **Local Development Setup:** `/workspace/Documentation/LocalDev.md`
- **Backend API Documentation:** `/workspace/Documentation/Backend/API.md`
- **Backend - Adding New Tables:** `/workspace/Documentation/Backend/AddNewTable.md`
- **Frontend Design System:** `/workspace/Documentation/Frontend/README.md`
- **Frontend Component Documentation:** `/workspace/Documentation/Frontend/components/`
  - Input.md - Form input components
  - Modal.md - Modal dialog patterns
  - Table.md - Data table components
  - TypeAheadDropdown.md - Autocomplete dropdowns

## Development Environment

This project uses a **Dev Container** for a consistent development environment. The dev container automatically handles all setup including:
- Go toolchain and dependencies
- Node.js and npm
- PostgreSQL database
- All required development tools

Simply open the project in VS Code and select "Reopen in Container" when prompted.

### VS Code Tasks

The project provides VS Code tasks for development:
- **Start Backend:** Runs the Go backend with hot-reload using Air
- **Start Frontend:** Runs the Vite dev server
- **Start Development Environment:** Runs both backend and frontend in parallel

## Quality Checks

After making code changes, run the appropriate quality checks based on the changed files. These checks mirror the CI/CD pipeline and should pass before considering the work complete.

### Frontend Changes (Frontend/**)

After modifying files in the `Frontend/` directory, run the following commands in order:

1. **Lint Check**
   ```bash
   cd Frontend && npm run lint
   ```
   - Must pass with no errors
   - Ensures code style consistency and catches common issues

2. **Build Check**
   ```bash
   cd Frontend && npm run build
   ```
   - Must complete successfully
   - Verifies TypeScript compilation and build process

3. **Test Suite**
   ```bash
   cd Frontend && npm run test
   ```
   - All tests must pass
   - Verifies functionality is not broken

### Backend Changes (Backend/**)

After modifying files in the `Backend/` directory, run the following commands in order:

1. **Verify Dependencies**
   ```bash
   cd Backend && go mod verify
   ```
   - Ensures go.mod and go.sum are in sync

2. **Build Check**
   ```bash
   cd Backend && go build -v ./...
   ```
   - Must compile successfully
   - Catches compilation errors

3. **Test Suite with Race Detection**
   ```bash
   cd Backend && go test -v -race -coverprofile=coverage.out ./...
   ```
   - All tests must pass
   - Race detector must not find any data races
   - Generates coverage report

4. **Display Coverage** (optional but recommended)
   ```bash
   cd Backend && go tool cover -func=coverage.out | grep total
   ```
   - Shows total test coverage percentage

### Notes

- If checks fail, fix the issues before proceeding
- Quality checks mirror CI/CD requirements to ensure successful merges

## Key Development Guidelines

### Backend Development

1. **Adding New Database Tables:**
   - Define model in `Backend/models/` with proper GORM tags
   - Add model to `AutoMigrate` in `Backend/database/database.go`
   - Follow existing patterns for foreign keys and relationships

2. **Adding New API Endpoints:**
   - Create handler function in appropriate `Backend/handlers/` file
   - Register route in `Backend/handlers/api.go`
   - Add authentication middleware if required
   - Document endpoint in `Documentation/Backend/API.md`
   - Write comprehensive tests in corresponding `*_test.go` file

3. **Testing:**
   - Write tests for all handlers in `*_test.go` files
   - Use test helpers from `Backend/handlers/test_helpers.go`
   - Test both success and error cases
   - Use SQLite in-memory database for tests

### Frontend Development

1. **Component Creation:**
   - Follow the design system in `Documentation/Frontend/README.md`
   - Use existing CSS variables for colors and spacing
   - Ensure components are responsive and accessible
   - Add proper TypeScript types

2. **Styling:**
   - Use CSS variables defined in the design system
   - Follow dark theme patterns
   - Maintain consistent spacing and typography

3. **Internationalization:**
   - Add translations to `src/i18n/` files
   - Use `useTranslation` hook for all user-facing text
   - Support both English and German

4. **Testing:**
   - Write unit tests for components using Vitest
   - Use Testing Library for component testing
   - Mock API calls with axios-mock-adapter

## CI/CD Pipeline

The project has automated CI/CD workflows that run on every push and pull request:

- **Frontend Pipeline:** Runs lint, build, and tests with coverage
- **Backend Pipeline:** Runs dependency verification, build, and tests with race detection
- **Docker Build:** Builds and deploys backend Docker image to Azure on main branch

Your local quality checks must match these CI/CD requirements to ensure successful merges.

## Code Quality Standards

- **Go:** Follow standard Go conventions, use `gofmt`
- **TypeScript/React:** Follow ESLint rules, use proper typing
- **Testing:** Maintain high test coverage, test edge cases
- **Documentation:** Update relevant docs when adding features
- **Git:** Write clear commit messages, reference issues

## Authentication Flow

The application supports multiple authentication methods:
1. **Magic Link:** Email-based passwordless authentication
2. **OAuth/OIDC:** Integration with Keycloak for SSO
3. **JWT Tokens:** Access and refresh token mechanism

See `Documentation/Backend/API.md` for detailed authentication endpoint documentation.

## Azure Integration

The backend integrates with Azure services:
- **Blob Storage:** For file uploads and storage
- **Communication Services:** For sending emails (magic links, notifications)

Credentials and configuration are managed via environment variables.

---

## Quick Reference

**Run all frontend checks:**
```bash
cd Frontend && npm run lint && npm run build && npm run test
```

**Run all backend checks:**
```bash
cd Backend && go mod verify && go build -v ./... && go test -v -race ./...
```

**Start development environment:**
Use VS Code task "Start Development Environment" or run both tasks manually.

---

*Remember: Quality checks are mandatory. No exceptions.*
