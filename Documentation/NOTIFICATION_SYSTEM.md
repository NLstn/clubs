# Notification System

The club management system includes a comprehensive notification system that allows users to receive both in-app and email notifications about various activities.

## Features

### In-App Notifications
- Notifications appear in the header with a bell icon
- Unread count badge shows the number of unread notifications
- Dropdown menu displays recent notifications with timestamps
- Click on notifications to mark them as read and navigate to relevant pages
- **Delete button** on each notification for individual removal
- **Smart navigation** - Notifications automatically navigate to the appropriate page:
  - **Invitation notifications** → Invites page (`/profile/invites`)
  - **Event notifications** → Club details page
  - **Fine notifications** → Club details page
  - **News notifications** → Club details page
- "Mark all read" option for bulk management

### Email Notifications
- Traditional email notifications sent to user's registered email
- Can be enabled/disabled independently from in-app notifications
- Uses existing email infrastructure (Azure Communication Services)

### Notification Types
1. **Member Added** - When you are added to a club (not sent for invite acceptances)
2. **Invite Received** - When you receive an invitation to join a club
3. **Event Created** - When new events are created in your clubs
4. **Fine Assigned** - When you are assigned a fine
5. **News Created** - When news is published in your clubs

### User Preferences
Users can control their notification preferences through:
- Profile → Notifications settings page
- Separate toggles for in-app and email notifications
- Settings are saved automatically when changed
- Default preferences are created for new users

## Technical Implementation

### Backend Components
- `models/notification.go` - Notification and preferences data models
  - `DeleteNotification()` - Function to delete individual notifications
- `handlers/notifications.go` - API endpoints for notification management
  - `DeleteNotification()` - Handler for DELETE requests to remove notifications
- `notifications/notifications.go` - Email notification functions
- Database tables: `notifications` and `user_notification_preferences`

### Frontend Components
- `components/layout/NotificationDropdown.tsx` - Header notification display with smart navigation
  - **Invitation notifications** (`type: 'invite_received'`) → Navigate to `/profile/invites`
  - **Event notifications** (`type: 'event_created'`) → Navigate to `/clubs/{clubId}`
  - **Fine notifications** (`type: 'fine_assigned'`) → Navigate to `/clubs/{clubId}`
  - **News notifications** (`type: 'news_created'`) → Navigate to `/clubs/{clubId}`
  - **Delete functionality** - × button on each notification for individual removal
  - Automatically closes dropdown when navigating
  - Displays appropriate icons for each notification type
- `pages/settings/NotificationSettings.tsx` - User preference management
- `hooks/useNotifications.ts` - React hook for notification state management
  - `deleteNotification()` - Function to delete individual notifications

### API Endpoints
- `GET /api/v1/notifications` - Retrieve user notifications
- `GET /api/v1/notifications/count` - Get unread notification count
- `PUT /api/v1/notifications/{id}` - Mark specific notification as read
- `DELETE /api/v1/notifications/{id}` - Delete specific notification
- `PUT /api/v1/notifications/mark-all-read` - Mark all notifications as read
- `GET /api/v1/notification-preferences` - Get user preferences
- `PUT /api/v1/notification-preferences` - Update user preferences

## Usage

### For Users
1. Click the bell icon in the header to view notifications
2. **Click on notifications** to navigate directly to relevant pages:
   - **Invitation notifications** → Navigate to your invites page (`/profile/invites`)
   - **Event, Fine, or News notifications** → Navigate to the related club page
3. **Delete individual notifications** by clicking the × button on each notification
4. Click on other notifications to mark them as read
5. Access notification settings via Profile → Notifications
6. Toggle in-app and email preferences for each notification type
7. View and manage notification history

### For Developers
The notification system is extensible for new notification types:

1. Add new notification type to the preferences model
2. Create notification when the event occurs:
   ```go
   models.CreateNotification(userID, "type", "title", "message", &clubID, nil, nil)
   ```
3. Add email notification function if needed
4. Update frontend preference settings
5. **Add navigation logic** in `NotificationDropdown.tsx` for new notification types:
   ```tsx
   // In handleNotificationClick function
   else if (notification.type === 'new_type' && notification.clubId) {
     setIsOpen(false);
     navigate(`/appropriate/path/${notification.clubId}`);
   }
   ```

## Migration from Legacy System
- Existing email notifications continue to work
- New in-app notifications are added alongside email notifications
- User preferences default to enabled for backward compatibility
- No breaking changes to existing notification code

## Invite Notification Behavior

The system implements a smart notification approach for invitations:

### When an Invite is Sent
- The invited user receives an "Invite Received" notification
- This notification includes the club name and invitation details
- The notification is linked to the specific invite for tracking

### When an Invite is Accepted
- The "Invite Received" notification is automatically removed
- **No "Member Added" notification is sent** (to avoid duplicate notifications)
- The member record is marked with `AcceptedViaInvite: true` for tracking

### When an Invite is Rejected
- The "Invite Received" notification is automatically removed
- No membership is created

### Direct Member Addition (Non-Invite)
- When admins add members directly (not via invite)
- Standard "Member Added" notification is sent
- Member record has `AcceptedViaInvite: false`

This approach ensures users receive relevant notifications without being overwhelmed by duplicate or redundant messages.