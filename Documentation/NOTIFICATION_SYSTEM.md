# Notification System

The club management system includes a comprehensive notification system that allows users to receive both in-app and email notifications about various activities.

## Features

### In-App Notifications
- Notifications appear in the header with a bell icon
- Unread count badge shows the number of unread notifications
- Dropdown menu displays recent notifications with timestamps
- Click on notifications to mark them as read
- "Mark all read" option for bulk management

### Email Notifications
- Traditional email notifications sent to user's registered email
- Can be enabled/disabled independently from in-app notifications
- Uses existing email infrastructure (Azure Communication Services)

### Notification Types
1. **Member Added** - When you are added to a club
2. **Event Created** - When new events are created in your clubs
3. **Fine Assigned** - When you are assigned a fine
4. **News Created** - When news is published in your clubs

### User Preferences
Users can control their notification preferences through:
- Profile → Notifications settings page
- Separate toggles for in-app and email notifications
- Settings are saved automatically when changed
- Default preferences are created for new users

## Technical Implementation

### Backend Components
- `models/notification.go` - Notification and preferences data models
- `handlers/notifications.go` - API endpoints for notification management
- `notifications/notifications.go` - Email notification functions
- Database tables: `notifications` and `user_notification_preferences`

### Frontend Components
- `components/layout/NotificationDropdown.tsx` - Header notification display
- `pages/settings/NotificationSettings.tsx` - User preference management
- `hooks/useNotifications.ts` - React hook for notification state management

### API Endpoints
- `GET /api/v1/notifications` - Retrieve user notifications
- `GET /api/v1/notifications/count` - Get unread notification count
- `PUT /api/v1/notifications/{id}` - Mark specific notification as read
- `PUT /api/v1/notifications/mark-all-read` - Mark all notifications as read
- `GET /api/v1/notification-preferences` - Get user preferences
- `PUT /api/v1/notification-preferences` - Update user preferences

## Usage

### For Users
1. Click the bell icon in the header to view notifications
2. Access notification settings via Profile → Notifications
3. Toggle in-app and email preferences for each notification type
4. View and manage notification history

### For Developers
The notification system is extensible for new notification types:

1. Add new notification type to the preferences model
2. Create notification when the event occurs:
   ```go
   models.CreateNotification(userID, "type", "title", "message", &clubID, nil, nil)
   ```
3. Add email notification function if needed
4. Update frontend preference settings

## Migration from Legacy System
- Existing email notifications continue to work
- New in-app notifications are added alongside email notifications
- User preferences default to enabled for backward compatibility
- No breaking changes to existing notification code