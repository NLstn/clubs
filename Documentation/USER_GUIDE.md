<div align="center">
  <img src="assets/logo.png" alt="Clubs Logo" width="150"/>
  
  # User Guide
  
  **Complete guide to using the Clubs Management Application**
</div>

---

## 📋 Table of Contents

1. [Getting Started](#getting-started)
2. [Authentication](#authentication)
3. [Dashboard](#dashboard)
4. [Managing Clubs](#managing-clubs)
5. [Members Management](#members-management)
6. [Events Management](#events-management)
7. [Fines Management](#fines-management)
8. [Shift Scheduling](#shift-scheduling)
9. [Teams](#teams)
10. [News & Announcements](#news--announcements)

---

## 🚀 Getting Started

### First Time Login

When you first access the Clubs application, you'll be presented with the login screen.

**Login Options:**
1. **Magic Link Authentication**: Enter your email to receive a login link
2. **Single Sign-On (SSO)**: Login via your organization's Keycloak account

### User Roles

The application supports multiple user roles:
- **Admin**: Full access to club management features
- **Member**: Access to view club information and participate in events
- **Guest**: Limited read-only access

---

## 🔐 Authentication

### Magic Link Login

1. Navigate to the login page
2. Enter your email address
3. Click "Send Magic Link"
4. Check your email for the login link
5. Click the link to automatically log in

**Security Features:**
- Links expire after 15 minutes
- One-time use only
- IP validation for security

### SSO Login via Keycloak

1. Click "Login with SSO"
2. Enter your organization credentials
3. Authorize the Clubs application
4. You'll be automatically redirected

---

## 🏠 Dashboard

The dashboard is your central hub for all club activities.

**Key Features:**
- **Recent Activity**: View latest updates across all your clubs
- **Upcoming Events**: See your next events and meetings
- **Quick Actions**: Fast access to common tasks
- **Club Switcher**: Easily switch between your clubs
- **Notifications**: Stay informed about important updates

**Dashboard Widgets:**
- Upcoming events calendar
- Recent news and announcements
- Pending tasks (RSVPs, fine payments, etc.)
- Club statistics overview

---

## 🏢 Managing Clubs

### Creating a Club

**Admin Only**

1. Click "Create Club" from the dashboard
2. Fill in club details:
   - Club name
   - Description
   - Club logo/image
   - Settings and preferences
3. Click "Create Club"

### Club Settings

Access club settings via the admin panel:

**General Settings:**
- Club name and description
- Logo and branding
- Visibility settings (public/private)
- Member join settings (open/invite-only)

**Advanced Settings:**
- Custom member roles
- Fine categories
- Event types
- Shift templates

### Joining a Club

**For Members:**

1. Browse available clubs
2. Click "Join Club" or "Request to Join"
3. For private clubs, wait for admin approval
4. Once approved, access club features

---

## 👥 Members Management

### Viewing Members

Navigate to the Members section to see:
- Complete member list
- Member roles and status
- Activity history
- Contact information

**Filtering & Sorting:**
- Filter by role (Admin, Member, etc.)
- Sort by name, join date, or activity
- Search by name or email

### Inviting Members

**Admin Only**

1. Navigate to Members → Invite Member
2. Enter email address
3. Select member role
4. Add optional welcome message
5. Click "Send Invitation"

**Invitation Process:**
- Member receives email invitation
- They can accept or decline
- Admins can track invitation status
- Invitations expire after 7 days

### Managing Join Requests

**Admin Only**

1. Navigate to Members → Join Requests
2. Review pending requests
3. View applicant information
4. Approve or decline request
5. Optional: Add a message

### Member Roles

**Changing Member Roles:**
1. Find member in members list
2. Click "Edit Member"
3. Select new role
4. Save changes

**Permission Levels:**
- **Admin**: Full club management access
- **Moderator**: Event and member management
- **Member**: Standard club participation
- **Guest**: Read-only access

---

## 📅 Events Management

### Viewing Events

The Events page displays:
- Upcoming events calendar
- Past events archive
- Event details and descriptions
- RSVP status

**Calendar Views:**
- Month view
- Week view
- List view
- Filter by event type

### Creating Events

**Admin/Moderator Only**

1. Click "Create Event"
2. Fill in event details:
   - Event name
   - Description
   - Date and time
   - Location
   - Event type
   - RSVP settings
3. For recurring events:
   - Select recurrence pattern (daily, weekly, monthly)
   - Set end date or number of occurrences
4. Click "Create Event"

**Event Types:**
- Meetings
- Training sessions
- Social gatherings
- Competitions
- Custom types

### RSVP Management

**For Members:**
1. View event details
2. Click RSVP button
3. Select: Attending, Maybe, or Not Attending
4. Add optional comment
5. Receive confirmation

**For Admins:**
- View RSVP list
- See attendance statistics
- Send reminders to non-responders
- Export attendee list

### Recurring Events

**How Recurring Events Work:**
- Master event template
- Individual instances can be modified
- Cancel single occurrence without affecting series
- Update all future events or just one

---

## 💰 Fines Management

### Viewing Fines

Navigate to Fines section to see:
- Outstanding fines
- Paid fines history
- Fine amount and reason
- Payment status

**Filter Options:**
- By status (pending, paid, waived)
- By member
- By date range
- By fine type

### Creating Fine Templates

**Admin Only**

1. Navigate to Fines → Templates
2. Click "Create Template"
3. Define template:
   - Fine name
   - Default amount
   - Description
   - Category
4. Save template

**Using Templates:**
- Quick fine creation with predefined values
- Maintain consistency
- Easy tracking by category

### Issuing Fines

**Admin Only**

1. Navigate to Fines → Issue Fine
2. Select member(s)
3. Choose fine template or create custom
4. Set amount and reason
5. Set due date
6. Click "Issue Fine"

**Bulk Fine Issuing:**
- Select multiple members
- Apply same fine to all
- Individual amounts can vary

### Payment Tracking

**Recording Payments:**
1. Find fine in fines list
2. Click "Record Payment"
3. Enter payment details:
   - Amount paid
   - Payment date
   - Payment method
   - Optional notes
4. Save payment

**Payment Methods:**
- Cash
- Bank transfer
- Digital payment
- Other (specify)

---

## 🕐 Shift Scheduling

### Viewing Shifts

The Shifts section shows:
- Shift calendar
- Assigned shifts
- Upcoming shifts
- Shift history

**Views:**
- Calendar view
- List view
- My shifts view
- Team shifts view

### Creating Shift Schedules

**Admin Only**

1. Navigate to Shifts → Create Schedule
2. Define schedule:
   - Schedule name
   - Description
   - Shift type
   - Duration
   - Required members per shift
3. Generate shift slots
4. Assign members

### Assigning Members to Shifts

**Methods:**
1. **Manual Assignment**:
   - Select shift
   - Choose member
   - Confirm assignment

2. **Auto-Assignment**:
   - Set member availability
   - Let system distribute fairly
   - Review and adjust

3. **Member Self-Selection**:
   - Members pick available slots
   - First-come, first-served
   - Admin approval if required

### Shift Reminders

- Automatic email reminders
- Configurable reminder timing
- SMS notifications (if configured)
- In-app notifications

---

## 👔 Teams

### Creating Teams

**Admin Only**

1. Navigate to Teams → Create Team
2. Enter team details:
   - Team name
   - Description
   - Team leader
3. Add team members
4. Save team

**Team Benefits:**
- Organize large clubs
- Team-specific events
- Targeted news/announcements
- Simplified member management

### Managing Team Members

1. Navigate to specific team
2. View current members
3. Add or remove members
4. Assign team roles
5. Track team activity

### Team Activities

Teams can have:
- Team-specific events
- Team news feed
- Team statistics
- Team competitions

---

## 📰 News & Announcements

### Viewing News

The News section displays:
- Latest announcements
- Club updates
- Important notifications
- Team news (if applicable)

**Sorting Options:**
- Most recent first
- By priority
- By category
- By team

### Creating News Posts

**Admin/Moderator Only**

1. Click "Create News Post"
2. Fill in details:
   - Title
   - Content (supports markdown)
   - Priority level
   - Target audience (all members, specific team)
   - Publish immediately or schedule
3. Attach images (optional)
4. Click "Publish"

**Post Categories:**
- General announcements
- Event updates
- Rule changes
- Achievements
- Custom categories

### News Notifications

Members are notified via:
- In-app notifications
- Email (based on preferences)
- Dashboard widget
- Mobile push (if app installed)

---

## ⚙️ User Settings

### Profile Management

1. Click your profile icon
2. Select "Settings"
3. Update:
   - Name and contact info
   - Profile picture
   - Bio/description
   - Privacy settings

### Notification Preferences

Configure how you receive notifications:
- Email notifications (on/off)
- In-app notifications (on/off)
- Notification frequency
- Specific notification types

### Language Settings

The application supports multiple languages:
- English
- German (Deutsch)

Change language in settings menu.

---

## 🔍 Search & Filters

### Global Search

Use the search bar to find:
- Clubs
- Members
- Events
- News posts
- Teams

**Search Tips:**
- Use quotation marks for exact phrases
- Search by name, email, or keyword
- Filter results by type

### Advanced Filters

Most list views support advanced filtering:
- Multiple criteria
- Date ranges
- Status filters
- Custom fields
- Save filter presets

---

## 📱 Mobile Experience

The application is fully responsive and works on:
- Desktop browsers
- Tablets
- Mobile phones

**Mobile Features:**
- Touch-optimized interface
- Responsive layouts
- Mobile navigation menu
- Optimized performance

---

## ❓ FAQ

### How do I reset my password?
The application uses passwordless authentication. Simply request a new magic link to log in.

### Can I be a member of multiple clubs?
Yes, you can join and participate in multiple clubs simultaneously.

### How do I leave a club?
Navigate to Club Settings and click "Leave Club". Admins cannot leave if they're the last admin.

### What happens to my data if I leave a club?
Your activity history remains for club records, but you lose access to club information.

### Can I export club data?
Admins can export member lists, event histories, and financial reports in CSV format.

### How do I report a bug or request a feature?
Contact your system administrator or use the feedback form in the application.

---

## 🆘 Support

For additional help:
- Check the [API Documentation](Backend/API.md)
- Review [Local Development Guide](LocalDev.md)
- Contact your club administrator
- Reach out to technical support

---

**Last Updated**: December 2024  
**Version**: 1.0
