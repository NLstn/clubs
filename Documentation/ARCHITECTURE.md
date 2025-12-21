<div align="center">
  <img src="assets/logo.png" alt="Clubs Logo" width="150"/>
  
  # Architecture Overview
  
  **System architecture and design of the Clubs Management Application**
</div>

---

## 📐 System Architecture

The Clubs application follows a modern three-tier architecture with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────────┐
│                        User Interface                            │
│  React 19 + TypeScript + Vite + React Router                   │
│  - Responsive dark-themed UI                                     │
│  - Component-based architecture                                  │
│  - i18n support (EN/DE)                                         │
└────────────────────────┬────────────────────────────────────────┘
                         │ HTTP/REST API (JSON)
                         │ OAuth2/OIDC, JWT, Magic Link Auth
┌────────────────────────┴────────────────────────────────────────┐
│                      Backend API Layer                           │
│  Go 1.25 + GORM + Custom HTTP Handlers                         │
│  - RESTful API (v1) & OData API (v2)                           │
│  - JWT authentication & authorization                            │
│  - Business logic & validation                                   │
│  - Job scheduler for background tasks                            │
└────────────────────────┬────────────────────────────────────────┘
                         │ SQL Queries (GORM ORM)
┌────────────────────────┴────────────────────────────────────────┐
│                      Data Storage Layer                          │
│  PostgreSQL 16                                                   │
│  - Relational database                                           │
│  - ACID transactions                                             │
│  - Auto-migration via GORM                                       │
└─────────────────────────────────────────────────────────────────┘

External Services:
┌──────────────────┐  ┌──────────────────┐  ┌────────────────────┐
│    Keycloak      │  │  Azure Blob      │  │  Azure Comm.       │
│  Identity & SSO  │  │    Storage       │  │    Services        │
└──────────────────┘  └──────────────────┘  └────────────────────┘
```

## 🔄 Request Flow

### 1. User Authentication Flow

```
User → Frontend → Keycloak (OAuth2/OIDC)
                     ↓
          Backend ← Token Validation
                     ↓
          JWT Generation & Session
                     ↓
          Frontend ← Access & Refresh Tokens
```

### 2. API Request Flow

```
Frontend → API Request (with JWT in Authorization header)
              ↓
        Rate Limiter Middleware
              ↓
        Authentication Middleware
              ↓
        Route Handler
              ↓
        Business Logic / Model Layer
              ↓
        Database (GORM)
              ↓
        Response (JSON) → Frontend
```

## 🗄️ Database Schema Overview

### Core Entities

**Clubs**: The primary organizational unit
- Manages members, events, fines, shifts
- Custom settings and configuration
- Image storage via Azure Blob

**Members**: Users associated with clubs
- Role-based permissions (Admin, Member, etc.)
- Invitation and join request workflows
- Activity tracking

**Events**: Club activities and meetings
- One-time or recurring events
- RSVP functionality
- Location and time management

**Fines**: Financial penalties management
- Template-based fine creation
- Payment tracking
- Custom fine types per club

**Shifts**: Work/duty scheduling
- Shift schedule templates
- Member assignment
- Date and time tracking

**Teams**: Sub-groups within clubs
- Team membership
- Hierarchical organization

**News**: Announcements and updates
- Club-wide or team-specific
- Timestamp tracking

### Key Relationships

```
Club ─┬─< Members
      ├─< Events ─< EventRSVPs
      ├─< Fines
      ├─< ShiftSchedules ─< Shifts
      ├─< Teams ─< TeamMembers
      └─< News
      
User ─< Members ─< Activity
```

## 🔐 Security Architecture

### Authentication Methods

1. **OAuth2/OIDC via Keycloak**
   - Single Sign-On (SSO) support
   - Enterprise-grade authentication
   - Token-based access

2. **Magic Link Authentication**
   - Passwordless email authentication
   - Time-limited tokens
   - Secure token validation

3. **JWT Tokens**
   - Access tokens (short-lived)
   - Refresh tokens (longer-lived, rotated)
   - Stateless authentication

### Authorization Model

- **Role-Based Access Control (RBAC)**
  - Club-level roles (Admin, Member, etc.)
  - Action-based permissions
  - Resource ownership validation

### Security Measures

- CSRF protection with state tokens
- Rate limiting on all endpoints
- SQL injection prevention via ORM
- XSS protection in frontend
- Secure password-less authentication
- Token rotation for refresh tokens

## 🚀 Deployment Architecture

### Frontend Deployment (Azure Static Web Apps)

```
Developer → GitHub Actions → Azure Static Web Apps
              ↓
         Build (Vite)
              ↓
         Deploy Static Files
              ↓
         CDN Distribution
```

### Backend Deployment (Azure Container Apps)

```
Developer → GitHub Actions → Docker Build → Azure Container Registry
                                               ↓
                                     Azure Container Apps
                                               ↓
                                     PostgreSQL Database
```

## 📊 Performance Considerations

### Frontend Optimizations
- Code splitting with React.lazy
- Vite's fast HMR for development
- Tree-shaking for minimal bundle size
- Asset optimization and caching

### Backend Optimizations
- GORM connection pooling
- Efficient query patterns with eager loading
- Indexed database columns
- Rate limiting to prevent abuse

### Caching Strategy
- Browser caching for static assets
- CDN caching for frontend
- Database query optimization
- JWT token caching (in-memory)

## 🔧 Development Workflow

1. **Local Development**
   - Dev Container with all dependencies
   - Hot-reload for both frontend and backend
   - Local PostgreSQL and Keycloak instances

2. **Version Control**
   - Git with feature branch workflow
   - Pull request reviews
   - Automated CI/CD checks

3. **Testing**
   - Backend: Go testing with testify
   - Frontend: Vitest + Testing Library
   - Integration tests for critical flows

4. **Deployment**
   - Automated deployment on main branch merge
   - Separate environments (dev/prod)
   - Docker containerization for backend

## 📚 Technology Choices Rationale

### Why Go for Backend?
- **Performance**: Compiled language with excellent concurrency
- **Simplicity**: Clean syntax, easy to maintain
- **Ecosystem**: Great libraries (GORM, JWT, etc.)
- **Deployment**: Single binary, easy containerization

### Why React 19 for Frontend?
- **Modern**: Latest React features and optimizations
- **TypeScript**: Type safety and better DX
- **Ecosystem**: Rich component libraries and tools
- **Performance**: Excellent rendering performance

### Why PostgreSQL?
- **Reliability**: ACID compliance, data integrity
- **Features**: Rich SQL support, JSON fields
- **Scalability**: Handles growth well
- **Open Source**: No licensing costs

### Why Azure?
- **Integration**: Seamless service integration
- **Scalability**: Auto-scaling capabilities
- **Security**: Enterprise-grade security
- **Global**: CDN and regional deployments

## 🔮 Future Enhancements

- Real-time notifications via WebSockets
- Mobile application (React Native)
- Advanced analytics and reporting
- Multi-language support expansion
- API versioning and GraphQL support
- Microservices architecture for scaling

---

For more detailed documentation, see:
- [Backend API Documentation](Backend/API.md)
- [Frontend Design System](Frontend/README.md)
- [Local Development Guide](LocalDev.md)
