<div align="center">
  <img src="assets/logo.png" alt="Clubs Logo" width="150"/>
  
  # Quick Start Guide
  
  **Get up and running with Clubs in 5 minutes**
</div>

---

## 🎯 Prerequisites

Before you begin, ensure you have:
- ✅ Docker Desktop installed (for Dev Container)
- ✅ Visual Studio Code with Dev Containers extension
- ✅ Git installed on your machine

---

## 🚀 Step 1: Clone the Repository

```bash
git clone https://github.com/NLstn/clubs.git
cd clubs
```

---

## 🐳 Step 2: Open in Dev Container

1. Open the project in Visual Studio Code
2. When prompted, click **"Reopen in Container"**
3. Wait for the container to build (first time takes ~5 minutes)

**What happens automatically:**
- PostgreSQL database is started
- Keycloak authentication server is configured
- Go and Node.js dependencies are installed
- Test user is created

---

## 🎮 Step 3: Start the Application

### Option A: Use VS Code Tasks (Recommended)

1. Press `Ctrl+Shift+P` (or `Cmd+Shift+P` on Mac)
2. Type "Tasks: Run Task"
3. Select **"Start Development Environment"**

This starts both backend and frontend simultaneously!

### Option B: Manual Start

**Terminal 1 - Backend:**
```bash
cd Backend
go run main.go
```

**Terminal 2 - Frontend:**
```bash
cd Frontend
npm run dev
```

---

## 🌐 Step 4: Access the Application

Once started, open your browser and navigate to:

**Frontend Application:**
```
http://localhost:5173
```

**Keycloak Admin Console:**
```
http://localhost:8081/admin
```

---

## 🔐 Step 5: Login

Use the pre-configured test account:

**Username:** `testuser`  
**Password:** `testpass`

Or use Magic Link authentication with any email address (in dev mode, check backend logs for the magic link).

---

## 🎉 Step 6: Explore Features

Once logged in, try these actions:

### 1. Create Your First Club
- Click "Create Club"
- Fill in the name and description
- Upload a club logo (optional)
- Click "Create"

### 2. Invite Members
- Navigate to your club
- Go to "Members" → "Invite Member"
- Enter an email address
- Send invitation

### 3. Create an Event
- Go to "Events" → "Create Event"
- Set event details
- Choose date and time
- Click "Create Event"

### 4. Explore the Dashboard
- View recent activity
- Check upcoming events
- Browse news and announcements

---

## 🛠️ Development Tips

### Hot Reload

Both frontend and backend support hot-reload:
- **Frontend**: Vite provides instant HMR
- **Backend**: Air provides automatic restart on file changes

### Running Tests

**Backend tests:**
```bash
cd Backend
go test ./...
```

**Frontend tests:**
```bash
cd Frontend
npm run test
```

### Linting

**Frontend lint:**
```bash
cd Frontend
npm run lint
```

### Database Access

Connect to the PostgreSQL database:
```bash
# Inside dev container
psql -h db -U clubs_dev -d clubs_dev
# Password: clubs_dev_password
```

---

## 📁 Project Structure Quick Reference

```
clubs/
├── Backend/          # Go backend API
│   ├── handlers/    # API route handlers
│   ├── models/      # Database models
│   ├── auth/        # Authentication logic
│   └── main.go      # Entry point
│
├── Frontend/        # React frontend
│   ├── src/
│   │   ├── components/  # Reusable components
│   │   ├── pages/       # Page components
│   │   ├── context/     # React contexts
│   │   └── utils/       # Utility functions
│   └── package.json
│
└── Documentation/   # Project documentation
    ├── Backend/     # Backend docs
    ├── Frontend/    # Frontend docs
    └── LocalDev.md  # Development guide
```

---

## 🔧 Common Commands

### Backend
```bash
# Run backend
go run main.go

# Build backend
go build

# Run tests
go test ./...

# Install dependencies
go mod download
```

### Frontend
```bash
# Start dev server
npm run dev

# Build for production
npm run build

# Run tests
npm run test

# Lint code
npm run lint
```

---

## 🐛 Troubleshooting

### Port Already in Use

If you get "port already in use" errors:

```bash
# Check what's using the port
lsof -i :5173  # Frontend port
lsof -i :8080  # Backend port
lsof -i :8081  # Keycloak port

# Kill the process
kill -9 <PID>
```

### Database Connection Issues

```bash
# Check if database is running
docker ps | grep postgres

# Restart database container
docker restart <container-id>

# Check logs
docker logs <container-id>
```

### Frontend Build Errors

```bash
# Clear node_modules and reinstall
cd Frontend
rm -rf node_modules package-lock.json
npm install
```

### Backend Build Errors

```bash
# Clean go cache
go clean -cache -modcache

# Download dependencies
cd Backend
go mod download
go mod verify
```

---

## 📚 Next Steps

Now that you're set up, explore these resources:

1. 📖 [User Guide](USER_GUIDE.md) - Comprehensive user documentation
2. 🏗️ [Architecture Overview](ARCHITECTURE.md) - System design and architecture
3. 🎨 [Frontend Design System](Frontend/README.md) - UI components and styling
4. 🔌 [API Documentation](Backend/API.md) - Backend API reference
5. 💾 [Adding New Tables](Backend/AddNewTable.md) - Database schema guide

---

## 🆘 Getting Help

- **Documentation**: Check the `/Documentation` folder
- **Issues**: Report bugs on GitHub Issues
- **Community**: Join the discussion on GitHub Discussions

---

## ✨ Tips for New Developers

1. **Read the docs first**: Familiarize yourself with the architecture
2. **Follow the design system**: Use existing components and patterns
3. **Write tests**: Maintain test coverage for new features
4. **Use TypeScript strictly**: Avoid `any` types
5. **Follow Go conventions**: Use `gofmt` and follow idiomatic Go
6. **Check CI before pushing**: Run lints and tests locally

---

**Happy Coding! 🎉**

For more detailed information, see the [Local Development Guide](LocalDev.md).
