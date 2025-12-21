<div align="center">
  <img src="Documentation/assets/logo.png" alt="Clubs Logo" width="200"/>
  
  # Clubs
  
  **A comprehensive club management application**
  
  [![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
  [![Backend](https://img.shields.io/badge/backend-Go-00ADD8.svg)](Backend/)
  [![Frontend](https://img.shields.io/badge/frontend-React-61DAFB.svg)](Frontend/)
</div>

---

## 📋 Overview

Clubs is a full-stack club management application designed to help organizations manage their clubs, members, events, fines, shifts, and more. The application provides a comprehensive platform for club administration with an intuitive user interface and powerful backend.

## ✨ Key Features

- 🏢 **Club Management**: Create and manage multiple clubs with customizable settings
- 👥 **Member Management**: Handle member registration, roles, and permissions
- 📅 **Event Scheduling**: Create and manage events with recurring event support
- 💰 **Fine Management**: Track and manage fines with customizable templates
- 📊 **Shift Scheduling**: Organize and assign shifts to members
- 👔 **Team Organization**: Create teams within clubs for better organization
- 📰 **News & Notifications**: Keep members informed with announcements
- 🔐 **Secure Authentication**: OAuth2/OIDC via Keycloak and Magic Link email authentication
- ☁️ **Azure Integration**: Seamless integration with Azure services

## 🚀 Getting Started

For detailed setup instructions, see [Local Development Guide](Documentation/LocalDev.md).

### Quick Start

1. Clone the repository
2. Open in VS Code with Dev Container support
3. Wait for the container to build and start
4. Access the application at `http://localhost:5173`

## 📚 Documentation

- [Local Development Setup](Documentation/LocalDev.md)
- [Backend API Documentation](Documentation/Backend/API.md)
- [Frontend Design System](Documentation/Frontend/README.md)
- [Adding New Tables](Documentation/Backend/AddNewTable.md)

## 🛠️ Technology Stack

### Backend
- **Language**: Go 1.25
- **Database**: PostgreSQL with GORM
- **Authentication**: OAuth2/OIDC (Keycloak), JWT, Magic Link
- **Cloud**: Azure (Blob Storage, Communication Services)

### Frontend
- **Framework**: React 19 with TypeScript
- **Build Tool**: Vite
- **Routing**: React Router v7
- **UI**: Custom design system with dark theme
- **i18n**: English and German support

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
