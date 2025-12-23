# iOS Clubs App Implementation Guide

## Project Overview

This is a native iOS app built with SwiftUI that mirrors the web frontend of the Clubs Management Application. The app follows Apple's Human Interface Guidelines and uses modern iOS design patterns.

## Project Structure

```
ios/
├── Models/
│   └── APIModels.swift          # Data models for API responses
├── Services/
│   ├── AuthenticationManager.swift  # Auth state and JWT management
│   ├── APIService.swift           # Network layer for API calls
│   └── KeychainService.swift      # Secure token storage
├── Views/
│   ├── Authentication/
│   │   └── LoginView.swift        # Magic link authentication
│   ├── Clubs/
│   │   └── ClubsListView.swift    # Clubs list and detail views
│   ├── Events/
│   │   └── EventsListView.swift   # Events list
│   ├── Fines/
│   │   └── FinesListView.swift    # Fines list with filtering
│   └── Profile/
│       └── ProfileView.swift      # User profile and logout
├── Utilities/
│   └── DesignSystem.swift         # Reusable UI components
├── iosApp.swift                   # App entry point
└── ContentView.swift              # Main tab bar view
```

## Key Features

### Authentication
- **Magic Link Authentication**: Request and verify magic links via email
- **Secure Token Storage**: JWT tokens stored securely in Keychain
- **Token Refresh**: Automatic token refresh before expiration
- **Session Management**: Handles token rotation and logout

### Clubs Management
- **Club Listing**: View all clubs the user is a member of
- **Club Details**: View club information, members, and events
- **Member Directory**: See club members and their roles
- **Statistics**: View member count and event statistics

### Events
- **Event Listing**: View all upcoming events across clubs
- **Event Details**: See event title, date, location, and attendance count
- **Recurring Events**: Indication of recurring events
- **Sorting**: Events sorted by start time

### Fines Management
- **Fine Listing**: View all fines with status
- **Status Filtering**: Filter fines by Pending, Paid, or Overdue
- **Color Coding**: Visual indicators for fine status
- **Amount Display**: Show fine amount in euros

### User Profile
- **User Information**: Display user name and email
- **Account Details**: Show user role and account information
- **Logout**: Secure logout that invalidates tokens

## Architecture

### MVVM with State Management
- Uses `@EnvironmentObject` for dependency injection
- `AuthenticationManager` handles authentication state and logic
- View models manage data fetching and state

### Networking
- `APIService` provides all API communication
- Automatic JWT token injection in headers
- Error handling with localized error messages
- Async/await for modern concurrency

### Security
- `KeychainService` for secure token storage
- JWT token validation and expiration checking
- Refresh token rotation for improved security

## Color Scheme

The app uses a dark theme inspired by the web frontend:
- **Primary Green**: #4CAF50 (for actions and success states)
- **Secondary Blue**: #646CFF (for secondary actions)
- **Red**: #F44336 (for destructive actions and errors)
- **Dark Backgrounds**: #242424 and #333333

## Design Patterns

### Navigation
- **Tab Bar Navigation**: Main navigation with 4 tabs (Clubs, Events, Fines, Profile)
- **Navigation Stacks**: Detail views using NavigationStack
- **Modal Presentations**: Alerts for confirmations

### Loading States
- Progress indicators during data fetching
- Empty states for empty lists
- Error banners for API errors

### User Feedback
- Loading spinners during operations
- Error alerts with retry options
- Success confirmations

## iOS Version Requirements

- **Minimum iOS**: 15.0+
- **Target iOS**: 17.0+
- Uses latest SwiftUI features (iOS 17)

## API Integration

The app connects to the backend API at `http://localhost:8080/api/v1`

### Main Endpoints Used
- `POST /auth/requestMagicLink` - Request authentication link
- `GET /auth/verifyMagicLink?token=` - Verify magic link
- `POST /auth/refreshToken` - Refresh access token
- `POST /auth/logout` - Logout user
- `GET /user` - Get current user info
- `GET /clubs` - Get user's clubs
- `GET /clubs/{id}` - Get club details
- `GET /clubs/{id}/members` - Get club members
- `GET /clubs/{id}/events` - Get club events
- `GET /events` - Get all events
- `GET /fines` - Get user's fines
- `GET /clubs/{id}/fines` - Get club fines

## Setup & Running

### Prerequisites
- Xcode 15.0+
- iOS deployment target 15.0+
- Backend server running at http://localhost:8080

### Building
1. Open `ios/ios.xcodeproj` in Xcode
2. Select the iOS simulator or device
3. Press Cmd+R to build and run

### Configuration
- Update `baseURL` in `APIService` if backend is on different address
- Configure Keychain service identifier if needed

## Design System Components

### Buttons
- `PrimaryButton` - Main action buttons (green background)
- `SecondaryButton` - Alternative actions

### Containers
- `CardContainer` - Reusable card for content grouping
- `ErrorBanner` - Error message display

### States
- `EmptyState` - Display when no data available
- `LoadingSkeleton` - Placeholder during loading
- Shimmer animation for better UX

## Future Enhancements

- [ ] Push notifications for events and fines
- [ ] Offline mode with local caching
- [ ] Rich member profiles
- [ ] Event RSVP functionality
- [ ] Fine payment integration
- [ ] News and announcements
- [ ] Shift scheduling
- [ ] Deep linking for notifications
- [ ] Widget support
- [ ] Share extensions for events

## Contributing

When adding new features:
1. Follow existing code structure and patterns
2. Use SwiftUI's latest features
3. Maintain dark theme consistency
4. Add proper error handling
5. Include loading states
6. Test on multiple device sizes
