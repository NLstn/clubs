# Club Management API Documentation

## Base URL
All API endpoints are prefixed with `/api/v1/`

## Authentication
The API uses JWT-based authentication with magic link email authentication. Most endpoints require authentication via Bearer token.

### Rate Limiting
- Authentication endpoints: 5 requests per minute per IP
- API endpoints: 30 requests per 5 seconds per IP

---

## Authentication Endpoints

### Request Magic Link
**Endpoint:** `POST /api/v1/auth/requestMagicLink`  
**Authentication:** None required  
**Rate Limit:** Auth limiter (5/min)

**Description:** Request a magic link to be sent to the provided email address for authentication.

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Responses:**
- `204 No Content` - Magic link sent successfully
- `400 Bad Request` - Email required
- `500 Internal Server Error` - Database error or failed to send email

---

### Verify Magic Link
**Endpoint:** `GET /api/v1/auth/verifyMagicLink?token={token}`  
**Authentication:** None required  
**Rate Limit:** Auth limiter (5/min)

**Description:** Verify a magic link token and return access/refresh tokens.

**Query Parameters:**
- `token` (required) - The magic link token

**Response:**
```json
{
  "access": "jwt_access_token",
  "refresh": "jwt_refresh_token"
}
```

**Responses:**
- `200 OK` - Token verified, returns access and refresh tokens
- `400 Bad Request` - Token required
- `401 Unauthorized` - Invalid or expired token
- `500 Internal Server Error` - User error, JWT error, or database error

---

### Refresh Access Token
**Endpoint:** `POST /api/v1/auth/refreshToken`  
**Authentication:** Refresh token in Authorization header  
**Rate Limit:** Auth limiter (5/min)

**Description:** Refresh an expired access token using a valid refresh token.

**Headers:**
- `Authorization: {refresh_token}`

**Response:**
```json
{
  "access": "new_jwt_access_token"
}
```

**Responses:**
- `200 OK` - New access token generated
- `401 Unauthorized` - Refresh token required or invalid
- `404 Not Found` - User not found
- `500 Internal Server Error` - Failed to generate access token

---

### Logout
**Endpoint:** `POST /api/v1/auth/logout`  
**Authentication:** Refresh token in Authorization header  
**Rate Limit:** Auth limiter (5/min)

**Description:** Logout user by invalidating the refresh token.

**Headers:**
- `Authorization: {refresh_token}`

**Responses:**
- `204 No Content` - Successfully logged out
- `401 Unauthorized` - Refresh token required or invalid
- `404 Not Found` - User not found
- `500 Internal Server Error` - Failed to delete refresh token

---

## Club Management Endpoints

### Get All Clubs
**Endpoint:** `GET /api/v1/clubs`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Get all clubs that the authenticated user is a member of.

**Response:**
```json
[
  {
    "id": "club-uuid",
    "name": "Club Name",
    "description": "Club Description"
  }
]
```

**Responses:**
- `200 OK` - List of clubs
- `401 Unauthorized` - Invalid or missing token
- `500 Internal Server Error` - Database error

---

### Get Club by ID
**Endpoint:** `GET /api/v1/clubs/{clubid}`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Get details of a specific club. User must be a member.

**Path Parameters:**
- `clubid` (UUID) - Club identifier

**Response:**
```json
{
  "id": "club-uuid",
  "name": "Club Name",
  "description": "Club Description"
}
```

**Responses:**
- `200 OK` - Club details
- `401 Unauthorized` - Invalid or missing token
- `403 Forbidden` - User is not a member of the club
- `404 Not Found` - Club not found
- `500 Internal Server Error` - Database error

---

### Create Club
**Endpoint:** `POST /api/v1/clubs`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Create a new club. The authenticated user becomes the owner.

**Request Body:**
```json
{
  "name": "Club Name",
  "description": "Club Description"
}
```

**Response:**
```json
{
  "id": "new-club-uuid",
  "name": "Club Name",
  "description": "Club Description"
}
```

**Responses:**
- `201 Created` - Club created successfully
- `400 Bad Request` - Invalid request body
- `401 Unauthorized` - Invalid or missing token
- `500 Internal Server Error` - Database error

---

### Update Club
**Endpoint:** `PATCH /api/v1/clubs/{clubid}`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Update club details. Only the club owner can update.

**Path Parameters:**
- `clubid` (UUID) - Club identifier

**Request Body:**
```json
{
  "name": "Updated Club Name",
  "description": "Updated Club Description"
}
```

**Response:**
```json
{
  "id": "club-uuid",
  "name": "Updated Club Name",
  "description": "Updated Club Description"
}
```

**Responses:**
- `200 OK` - Club updated successfully
- `400 Bad Request` - Invalid request body
- `401 Unauthorized` - Invalid or missing token
- `403 Forbidden` - User is not the club owner
- `404 Not Found` - Club not found
- `500 Internal Server Error` - Database error

---

## Member Management Endpoints

### Get Club Members
**Endpoint:** `GET /api/v1/clubs/{clubid}/members`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Get all members of a club. User must be a member.

**Path Parameters:**
- `clubid` (UUID) - Club identifier

**Response:**
```json
[
  {
    "id": "member-uuid",
    "userId": "user-uuid",
    "name": "Member Name",
    "role": "owner|admin|member"
  }
]
```

**Responses:**
- `200 OK` - List of club members
- `400 Bad Request` - Invalid club ID format
- `401 Unauthorized` - Invalid or missing token
- `404 Not Found` - Club not found
- `500 Internal Server Error` - Database error

---

### Update Member Role
**Endpoint:** `PATCH /api/v1/clubs/{clubid}/members/{memberid}`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Update a member's role. Only admins/owners can update roles.

**Path Parameters:**
- `clubid` (UUID) - Club identifier
- `memberid` (UUID) - Member identifier

**Request Body:**
```json
{
  "role": "owner|admin|member"
}
```

**Responses:**
- `204 No Content` - Role updated successfully
- `400 Bad Request` - Invalid club/member ID format or invalid role
- `401 Unauthorized` - Invalid or missing token
- `403 Forbidden` - User is not an admin/owner
- `404 Not Found` - Club not found
- `500 Internal Server Error` - Database error

---

### Remove Club Member
**Endpoint:** `DELETE /api/v1/clubs/{clubid}/members/{memberid}`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Remove a member from the club. Only the club owner can remove members.

**Path Parameters:**
- `clubid` (UUID) - Club identifier
- `memberid` (UUID) - Member identifier

**Responses:**
- `204 No Content` - Member removed successfully
- `400 Bad Request` - Invalid club/member ID format
- `401 Unauthorized` - Invalid or missing token
- `403 Forbidden` - User is not the club owner
- `404 Not Found` - Club or member not found
- `500 Internal Server Error` - Database error

---

### Check Admin Rights
**Endpoint:** `GET /api/v1/clubs/{clubid}/isAdmin`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Check if the authenticated user has admin rights (owner or admin) in the club.

**Path Parameters:**
- `clubid` (UUID) - Club identifier

**Response:**
```json
{
  "isAdmin": true
}
```

**Responses:**
- `200 OK` - Admin status returned
- `401 Unauthorized` - Invalid or missing token
- `404 Not Found` - Club not found
- `500 Internal Server Error` - Database error

---

## Shift Management Endpoints

### Get Club Shifts
**Endpoint:** `GET /api/v1/clubs/{clubid}/shifts`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Get all shifts for a club.

**Path Parameters:**
- `clubid` (UUID) - Club identifier

**Response:**
```json
[
  {
    "id": "shift-uuid",
    "startTime": "2024-01-01T09:00:00Z",
    "endTime": "2024-01-01T17:00:00Z"
  }
]
```

**Responses:**
- `200 OK` - List of shifts
- `400 Bad Request` - Invalid club ID format
- `401 Unauthorized` - Invalid or missing token
- `500 Internal Server Error` - Database error

---

### Create Shift
**Endpoint:** `POST /api/v1/clubs/{clubid}/shifts`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Create a new shift for a club.

**Path Parameters:**
- `clubid` (UUID) - Club identifier

**Request Body:**
```json
{
  "startTime": "2024-01-01T09:00:00Z",
  "endTime": "2024-01-01T17:00:00Z"
}
```

**Response:**
```json
{
  "id": "new-shift-uuid"
}
```

**Responses:**
- `201 Created` - Shift created successfully
- `400 Bad Request` - Invalid club ID format, invalid request body, or missing times
- `401 Unauthorized` - Invalid or missing token
- `500 Internal Server Error` - Database error

---

### Get Shift Members
**Endpoint:** `GET /api/v1/clubs/{clubid}/shifts/{shiftid}/members`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Get all members assigned to a specific shift.

**Path Parameters:**
- `clubid` (UUID) - Club identifier
- `shiftid` (UUID) - Shift identifier

**Response:**
```json
[
  {
    "id": "user-uuid",
    "name": "Member Name"
  }
]
```

**Responses:**
- `200 OK` - List of shift members
- `400 Bad Request` - Invalid club/shift ID format
- `401 Unauthorized` - Invalid or missing token
- `500 Internal Server Error` - Database error

---

### Add Member to Shift
**Endpoint:** `POST /api/v1/clubs/{clubid}/shifts/{shiftid}/members`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Add a member to a shift.

**Path Parameters:**
- `clubid` (UUID) - Club identifier
- `shiftid` (UUID) - Shift identifier

**Request Body:**
```json
{
  "userId": "user-uuid"
}
```

**Response:**
```json
{
  "message": "Member added to shift successfully"
}
```

**Responses:**
- `201 Created` - Member added to shift
- `400 Bad Request` - Invalid club/shift/user ID format or invalid request body
- `401 Unauthorized` - Invalid or missing token
- `500 Internal Server Error` - Database error

---

### Remove Member from Shift
**Endpoint:** `DELETE /api/v1/clubs/{clubid}/shifts/{shiftid}/members/{memberid}`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Remove a member from a shift.

**Path Parameters:**
- `clubid` (UUID) - Club identifier
- `shiftid` (UUID) - Shift identifier
- `memberid` (UUID) - Member identifier

**Response:**
```json
{
  "message": "Member removed from shift successfully"
}
```

**Responses:**
- `200 OK` - Member removed from shift
- `400 Bad Request` - Invalid club/shift/member ID format
- `401 Unauthorized` - Invalid or missing token
- `500 Internal Server Error` - Database error

---

## Fine Management Endpoints

### Get Club Fines
**Endpoint:** `GET /api/v1/clubs/{clubid}/fines`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Get all fines for a club. User must be a member.

**Path Parameters:**
- `clubid` (UUID) - Club identifier

**Response:**
```json
[
  {
    "id": "fine-uuid",
    "userId": "user-uuid",
    "userName": "User Name",
    "reason": "Fine reason",
    "amount": 25.50,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z",
    "paid": false
  }
]
```

**Responses:**
- `200 OK` - List of fines
- `401 Unauthorized` - Invalid or missing token
- `403 Forbidden` - User is not a member of the club
- `404 Not Found` - Club not found
- `500 Internal Server Error` - Database error

---

### Create Fine
**Endpoint:** `POST /api/v1/clubs/{clubid}/fines`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Create a new fine for a user. Only admins/owners can create fines.

**Path Parameters:**
- `clubid` (UUID) - Club identifier

**Request Body:**
```json
{
  "userId": "user-uuid",
  "reason": "Fine reason",
  "amount": 25.50
}
```

**Response:**
```json
{
  "id": "new-fine-uuid",
  "userId": "user-uuid",
  "reason": "Fine reason",
  "amount": 25.50,
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:00:00Z",
  "paid": false
}
```

**Responses:**
- `201 Created` - Fine created successfully
- `400 Bad Request` - Invalid request payload or missing required fields
- `401 Unauthorized` - Invalid or missing token
- `403 Forbidden` - User is not an admin/owner
- `404 Not Found` - Club not found
- `500 Internal Server Error` - Database error

---

## Fine Template Management Endpoints

### Get Club Fine Templates
**Endpoint:** `GET /api/v1/clubs/{clubid}/fine-templates`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Get all fine templates for a club. User must be a member.

**Path Parameters:**
- `clubid` (UUID) - Club identifier

**Response:**
```json
[
  {
    "id": "template-uuid",
    "club_id": "club-uuid",
    "description": "Late arrival",
    "amount": 25.50,
    "created_at": "2024-01-01T10:00:00Z",
    "created_by": "user-uuid",
    "updated_at": "2024-01-01T10:00:00Z",
    "updated_by": "user-uuid"
  }
]
```

**Responses:**
- `200 OK` - List of fine templates
- `401 Unauthorized` - Invalid or missing token
- `403 Forbidden` - User is not a member of the club
- `404 Not Found` - Club not found
- `500 Internal Server Error` - Database error

---

### Create Fine Template
**Endpoint:** `POST /api/v1/clubs/{clubid}/fine-templates`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Create a new fine template. Only admins/owners can create templates.

**Path Parameters:**
- `clubid` (UUID) - Club identifier

**Request Body:**
```json
{
  "description": "Late arrival",
  "amount": 25.50
}
```

**Response:**
```json
{
  "id": "new-template-uuid",
  "club_id": "club-uuid",
  "description": "Late arrival",
  "amount": 25.50,
  "created_at": "2024-01-01T10:00:00Z",
  "created_by": "user-uuid",
  "updated_at": "2024-01-01T10:00:00Z",
  "updated_by": "user-uuid"
}
```

**Responses:**
- `201 Created` - Fine template created successfully
- `400 Bad Request` - Invalid request payload or missing required fields
- `401 Unauthorized` - Invalid or missing token
- `403 Forbidden` - User is not an admin/owner
- `404 Not Found` - Club not found
- `500 Internal Server Error` - Database error

---

### Update Fine Template
**Endpoint:** `PUT /api/v1/clubs/{clubid}/fine-templates/{templateid}`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Update an existing fine template. Only admins/owners can update templates.

**Path Parameters:**
- `clubid` (UUID) - Club identifier
- `templateid` (UUID) - Fine template identifier

**Request Body:**
```json
{
  "description": "Very late arrival",
  "amount": 30.00
}
```

**Response:**
```json
{
  "id": "template-uuid",
  "club_id": "club-uuid",
  "description": "Very late arrival",
  "amount": 30.00,
  "created_at": "2024-01-01T10:00:00Z",
  "created_by": "user-uuid",
  "updated_at": "2024-01-01T12:00:00Z",
  "updated_by": "user-uuid"
}
```

**Responses:**
- `200 OK` - Fine template updated successfully
- `400 Bad Request` - Invalid request payload or missing required fields
- `401 Unauthorized` - Invalid or missing token
- `403 Forbidden` - User is not an admin/owner
- `404 Not Found` - Club or template not found
- `500 Internal Server Error` - Database error

---

### Delete Fine Template
**Endpoint:** `DELETE /api/v1/clubs/{clubid}/fine-templates/{templateid}`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Delete a fine template. Only admins/owners can delete templates.

**Path Parameters:**
- `clubid` (UUID) - Club identifier
- `templateid` (UUID) - Fine template identifier

**Responses:**
- `204 No Content` - Fine template deleted successfully
- `401 Unauthorized` - Invalid or missing token
- `403 Forbidden` - User is not an admin/owner
- `404 Not Found` - Club or template not found
- `500 Internal Server Error` - Database error

---

## Join Request Management Endpoints

### Create Join Request
**Endpoint:** `POST /api/v1/clubs/{clubid}/joinRequests`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Create a join request for a user to join a club. Only club owners can create join requests.

**Path Parameters:**
- `clubid` (UUID) - Club identifier

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Responses:**
- `201 Created` - Join request created successfully
- `400 Bad Request` - Invalid request payload or missing email
- `401 Unauthorized` - Invalid or missing token
- `403 Forbidden` - User is not the club owner
- `404 Not Found` - Club not found
- `500 Internal Server Error` - Database error

---

### Get Club Join Requests
**Endpoint:** `GET /api/v1/clubs/{clubid}/joinRequests`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Get all join requests for a club. Only club owners can view join requests.

**Path Parameters:**
- `clubid` (UUID) - Club identifier

**Response:**
```json
[
  {
    "id": "request-uuid",
    "email": "user@example.com",
    "created_at": "2024-01-01T10:00:00Z"
  }
]
```

**Responses:**
- `200 OK` - List of join requests
- `400 Bad Request` - Invalid club ID format
- `401 Unauthorized` - Invalid or missing token
- `403 Forbidden` - User is not the club owner
- `500 Internal Server Error` - Database error

---

### Get User Join Requests
**Endpoint:** `GET /api/v1/joinRequests`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Get all join requests for the authenticated user.

**Response:**
```json
[
  {
    "id": "request-uuid",
    "clubName": "Club Name"
  }
]
```

**Responses:**
- `200 OK` - List of user's join requests
- `401 Unauthorized` - Invalid or missing token
- `500 Internal Server Error` - Database error

---

### Accept Join Request
**Endpoint:** `POST /api/v1/joinRequests/{requestid}/accept`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Accept a join request. Only authorized users can accept.

**Path Parameters:**
- `requestid` (UUID) - Join request identifier

**Responses:**
- `204 No Content` - Join request accepted
- `400 Bad Request` - Invalid request ID format
- `401 Unauthorized` - Invalid or missing token, or user cannot edit this request
- `500 Internal Server Error` - Database error

---

### Reject Join Request
**Endpoint:** `POST /api/v1/joinRequests/{requestid}/reject`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Reject a join request. Only authorized users can reject.

**Path Parameters:**
- `requestid` (UUID) - Join request identifier

**Responses:**
- `204 No Content` - Join request rejected
- `400 Bad Request` - Invalid request ID format
- `401 Unauthorized` - Invalid or missing token, or user cannot edit this request
- `500 Internal Server Error` - Database error

---

## User Profile Endpoints

### Get Current User Profile
**Endpoint:** `GET /api/v1/me`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Get the authenticated user's profile information.

**Response:**
```json
{
  "id": "user-uuid",
  "name": "User Name",
  "email": "user@example.com"
}
```

**Responses:**
- `200 OK` - User profile
- `401 Unauthorized` - Invalid or missing token
- `500 Internal Server Error` - Database error

---

### Update Current User Profile
**Endpoint:** `PUT /api/v1/me`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Update the authenticated user's profile information.

**Request Body:**
```json
{
  "name": "Updated User Name"
}
```

**Responses:**
- `204 No Content` - Profile updated successfully
- `400 Bad Request` - Invalid request body or name required
- `401 Unauthorized` - Invalid or missing token
- `500 Internal Server Error` - Failed to update user

---

### Get Current User's Fines
**Endpoint:** `GET /api/v1/me/fines`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Get all unpaid fines for the authenticated user across all clubs.

**Response:**
```json
[
  {
    "id": "fine-uuid",
    "clubId": "club-uuid",
    "clubName": "Club Name",
    "reason": "Fine reason",
    "amount": 25.50,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z",
    "paid": false,
    "createdByName": "Admin Name"
  }
]
```

**Responses:**
- `200 OK` - List of user's unpaid fines
- `401 Unauthorized` - Invalid or missing token
- `500 Internal Server Error` - Database error

**Note:** This endpoint returns only unpaid fines for the user. Administrators can see all fines (paid and unpaid) using the club fines endpoint.

---

## Event Management Endpoints

### Get Club Events
**Endpoint:** `GET /api/v1/clubs/{clubid}/events`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Get all events for a club. User must be a member.

**Path Parameters:**
- `clubid` (UUID) - Club identifier

**Response:**
```json
[
  {
    "id": "event-uuid",
    "club_id": "club-uuid",
    "name": "Event Name",
    "start_date": "2024-06-01",
    "start_time": "10:00",
    "end_date": "2024-06-01",
    "end_time": "12:00",
    "created_at": "2024-01-01T10:00:00Z",
    "created_by": "user-uuid",
    "updated_at": "2024-01-01T10:00:00Z",
    "updated_by": "user-uuid"
  }
]
```

**Responses:**
- `200 OK` - List of events
- `400 Bad Request` - Invalid club ID format
- `401 Unauthorized` - Invalid or missing token
- `403 Forbidden` - User is not a member of the club
- `404 Not Found` - Club not found
- `500 Internal Server Error` - Database error

---

### Create Event
**Endpoint:** `POST /api/v1/clubs/{clubid}/events`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Create a new event for a club. Only club admins/owners can create events.

**Path Parameters:**
- `clubid` (UUID) - Club identifier

**Request Body:**
```json
{
  "name": "Event Name",
  "start_date": "2024-06-01",
  "start_time": "10:00",
  "end_date": "2024-06-01",
  "end_time": "12:00"
}
```

**Response:**
```json
{
  "id": "new-event-uuid",
  "club_id": "club-uuid",
  "name": "Event Name",
  "start_date": "2024-06-01",
  "start_time": "10:00",
  "end_date": "2024-06-01",
  "end_time": "12:00",
  "created_at": "2024-01-01T10:00:00Z",
  "created_by": "user-uuid",
  "updated_at": "2024-01-01T10:00:00Z",
  "updated_by": "user-uuid"
}
```

**Responses:**
- `201 Created` - Event created successfully
- `400 Bad Request` - Invalid club ID format, invalid request body, or invalid date/time format
- `401 Unauthorized` - Invalid or missing token
- `403 Forbidden` - User is not an admin/owner
- `404 Not Found` - Club not found
- `500 Internal Server Error` - Database error

---

### Update Event
**Endpoint:** `PUT /api/v1/clubs/{clubid}/events/{eventid}`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Update an existing event. Only club admins/owners can update events.

**Path Parameters:**
- `clubid` (UUID) - Club identifier
- `eventid` (UUID) - Event identifier

**Request Body:**
```json
{
  "name": "Updated Event Name",
  "start_date": "2024-06-02",
  "start_time": "11:00",
  "end_date": "2024-06-02",
  "end_time": "13:00"
}
```

**Response:**
```json
{
  "id": "event-uuid",
  "club_id": "club-uuid",
  "name": "Updated Event Name",
  "start_date": "2024-06-02",
  "start_time": "11:00",
  "end_date": "2024-06-02",
  "end_time": "13:00",
  "created_at": "2024-01-01T10:00:00Z",
  "created_by": "user-uuid",
  "updated_at": "2024-01-01T12:00:00Z",
  "updated_by": "user-uuid"
}
```

**Responses:**
- `200 OK` - Event updated successfully
- `400 Bad Request` - Invalid club/event ID format, invalid request body, or invalid date/time format
- `401 Unauthorized` - Invalid or missing token
- `403 Forbidden` - User is not an admin/owner
- `404 Not Found` - Club or event not found
- `500 Internal Server Error` - Database error

---

### Delete Event
**Endpoint:** `DELETE /api/v1/clubs/{clubid}/events/{eventid}`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Delete an event. Only club admins/owners can delete events. This will also delete all RSVPs for the event.

**Path Parameters:**
- `clubid` (UUID) - Club identifier
- `eventid` (UUID) - Event identifier

**Responses:**
- `204 No Content` - Event deleted successfully
- `400 Bad Request` - Invalid club/event ID format
- `401 Unauthorized` - Invalid or missing token
- `403 Forbidden` - User is not an admin/owner
- `404 Not Found` - Club or event not found
- `500 Internal Server Error` - Database error

---

### Get Upcoming Events
**Endpoint:** `GET /api/v1/clubs/{clubid}/events/upcoming`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Get upcoming events for a club with user's RSVP status. User must be a member.

**Path Parameters:**
- `clubid` (UUID) - Club identifier

**Response:**
```json
[
  {
    "id": "event-uuid",
    "club_id": "club-uuid",
    "name": "Event Name",
    "start_date": "2024-06-01",
    "start_time": "10:00",
    "end_date": "2024-06-01",
    "end_time": "12:00",
    "created_at": "2024-01-01T10:00:00Z",
    "created_by": "user-uuid",
    "updated_at": "2024-01-01T10:00:00Z",
    "updated_by": "user-uuid",
    "user_rsvp": {
      "id": "rsvp-uuid",
      "event_id": "event-uuid",
      "user_id": "user-uuid",
      "response": "yes",
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z"
    }
  }
]
```

**Responses:**
- `200 OK` - List of upcoming events with RSVP status
- `400 Bad Request` - Invalid club ID format
- `401 Unauthorized` - Invalid or missing token
- `403 Forbidden` - User is not a member of the club
- `404 Not Found` - Club not found
- `500 Internal Server Error` - Database error

---

### RSVP to Event
**Endpoint:** `POST /api/v1/clubs/{clubid}/events/{eventid}/rsvp`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Create or update an RSVP response for an event. User must be a club member.

**Path Parameters:**
- `clubid` (UUID) - Club identifier
- `eventid` (UUID) - Event identifier

**Request Body:**
```json
{
  "response": "yes"
}
```

**Valid responses:** `"yes"` or `"no"`

**Response:**
```json
{
  "status": "success"
}
```

**Responses:**
- `200 OK` - RSVP updated successfully
- `400 Bad Request` - Invalid club/event ID format, invalid request body, or invalid response value
- `401 Unauthorized` - Invalid or missing token
- `403 Forbidden` - User is not a member of the club
- `404 Not Found` - Club or event not found
- `500 Internal Server Error` - Database error

---

### Get Event RSVPs
**Endpoint:** `GET /api/v1/clubs/{clubid}/events/{eventid}/rsvps`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Get RSVP counts and details for an event. Only club admins/owners can view RSVPs.

**Path Parameters:**
- `clubid` (UUID) - Club identifier
- `eventid` (UUID) - Event identifier

**Response:**
```json
{
  "counts": {
    "yes": 5,
    "no": 2
  },
  "rsvps": [
    {
      "id": "rsvp-uuid",
      "event_id": "event-uuid",
      "user_id": "user-uuid",
      "response": "yes",
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z",
      "user": {
        "id": "user-uuid",
        "name": "User Name",
        "email": "user@example.com"
      }
    }
  ]
}
```

**Responses:**
- `200 OK` - RSVP counts and details
- `400 Bad Request` - Invalid club/event ID format
- `401 Unauthorized` - Invalid or missing token
- `403 Forbidden` - User is not an admin/owner
- `404 Not Found` - Club or event not found
- `500 Internal Server Error` - Database error

---

## Club Settings Endpoints

### Get Club Settings
**Endpoint:** `GET /api/v1/{clubid}/settings`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Get settings for a specific club.

**Path Parameters:**
- `clubid` (UUID) - Club identifier

**Response:**
```json
{
  "setting1": "value1",
  "setting2": "value2"
}
```

**Responses:**
- `200 OK` - Club settings
- `400 Bad Request` - Club ID is required
- `401 Unauthorized` - Invalid or missing token
- `404 Not Found` - Club not found
- `500 Internal Server Error` - Database error

---

### Update Club Settings
**Endpoint:** `POST /api/v1/{clubid}/settings`  
**Authentication:** Bearer token required  
**Rate Limit:** API limiter (30/5s)

**Description:** Update settings for a specific club.

**Path Parameters:**
- `clubid` (UUID) - Club identifier

**Request Body:**
```json
{
  "setting1": "new_value1",
  "setting2": "new_value2"
}
```

**Responses:**
- `204 No Content` - Settings updated successfully
- `400 Bad Request` - Club ID is required or invalid request body
- `401 Unauthorized` - Invalid or missing token
- `404 Not Found` - Club not found
- `500 Internal Server Error` - Database error

---

## Error Responses

All endpoints may return the following error responses:

### Common HTTP Status Codes
- `400 Bad Request` - Invalid request format, missing required fields, or invalid data
- `401 Unauthorized` - Missing, invalid, or expired authentication token
- `403 Forbidden` - Valid authentication but insufficient permissions
- `404 Not Found` - Requested resource does not exist
- `405 Method Not Allowed` - HTTP method not supported for this endpoint
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server-side error, database issues, or unexpected errors

### Error Response Format
Most error responses include a plain text error message:
```
Error message describing the issue
```

Some endpoints may return JSON error objects for structured error handling.

---

## Notes

1. All UUIDs must be valid UUID format
2. All timestamps are in ISO 8601 format (RFC 3339) and use camelCase field names (e.g., `createdAt`, `updatedAt`)
3. **Date and Time Format Conventions:**
   - **Date fields**: Use format `YYYY-MM-DD` (e.g., `"2024-06-01"`)
   - **Time fields**: Use format `HH:MM` in 24-hour format (e.g., `"10:00"`, `"15:30"`)
   - **Full timestamps**: Use ISO 8601 format (RFC 3339) with timezone (e.g., `"2024-01-01T10:00:00Z"`)
   - Event endpoints use separate date and time fields for better UX, while other endpoints use full timestamps
4. All monetary amounts are represented as floating-point numbers
5. Authentication is required for all endpoints except magic link request and verification
6. Rate limiting is enforced per IP address
7. CORS is enabled with permissive settings for development
8. JSON field names follow camelCase convention to ensure frontend compatibility
