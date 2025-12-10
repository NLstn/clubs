# OData v2 API Documentation

## Base URL
All OData v2 API endpoints are prefixed with `/api/v2/`

## Authentication
All OData v2 endpoints require JWT-based authentication via Bearer token in the Authorization header.

## OData v4.01 Support
The API implements the OData v4.01 specification and supports standard query options:
- `$filter` - Filter results based on conditions
- `$select` - Select specific properties
- `$expand` - Expand related entities
- `$orderby` - Sort results
- `$top` - Limit number of results
- `$skip` - Skip results for pagination
- `$count` - Include total count of results

## Field Naming Convention
All OData v2 API responses use **PascalCase** for field names, following the OData protocol standard.

---

## Timeline Entity

### Overview
The Timeline entity provides a unified view of activities, events, and news from all clubs the user is a member of. This is a **virtual entity** that aggregates data from multiple sources server-side.

### Endpoint
`GET /api/v2/TimelineItems`

### Description
Retrieves a unified timeline of activities, events, and news sorted by timestamp (most recent first). The timeline includes items from all clubs where the authenticated user is a member.

### Response Structure

```json
{
  "value": [
    {
      "ID": "string",          // Format: "{type}-{id}" (e.g., "event-123", "news-456", "activity-789")
      "ClubID": "string",      // UUID of the club
      "ClubName": "string",    // Name of the club
      "Type": "string",        // One of: "activity", "event", "news"
      "Title": "string",       // Title/name of the item
      "Content": "string",     // Optional: Content/description
      "Timestamp": "datetime", // ISO 8601 timestamp for sorting
      "CreatedAt": "datetime", // ISO 8601 creation timestamp
      "UpdatedAt": "datetime", // ISO 8601 last update timestamp
      
      // Event-specific fields (only present when Type="event")
      "StartTime": "datetime", // Event start time
      "EndTime": "datetime",   // Event end time
      "Location": "string",    // Event location
      "UserRSVP": {            // User's RSVP status if available
        "ID": "string",
        "EventID": "string",
        "UserID": "string",
        "Response": "string",  // "yes" or "no"
        "CreatedAt": "datetime",
        "UpdatedAt": "datetime"
      },
      
      // Activity-specific fields (only present when Type="activity")
      "Actor": "string",       // User ID who initiated the activity
      "ActorName": "string",   // Name of the user who initiated
      
      // Additional metadata
      "Metadata": {}           // Type-specific additional information
    }
  ]
}
```

### Timeline Item Types

#### 1. Activity Items (Type="activity")
Represents user activities such as:
- Member joining a club
- Role changes (promotion/demotion)
- Member-specific activities

**Example:**
```json
{
  "ID": "activity-abc123",
  "ClubID": "club-uuid",
  "ClubName": "Soccer Club",
  "Type": "activity",
  "Title": "New member joined",
  "Content": "John Doe has joined the club",
  "Timestamp": "2024-12-10T10:00:00Z",
  "CreatedAt": "2024-12-10T10:00:00Z",
  "UpdatedAt": "2024-12-10T10:00:00Z",
  "Actor": "user-uuid",
  "Metadata": {
    "old_role": "member",
    "new_role": "admin"
  }
}
```

#### 2. Event Items (Type="event")
Represents upcoming events in clubs.

**Example:**
```json
{
  "ID": "event-xyz789",
  "ClubID": "club-uuid",
  "ClubName": "Soccer Club",
  "Type": "event",
  "Title": "Weekend Practice Session",
  "Timestamp": "2024-12-15T14:00:00Z",
  "CreatedAt": "2024-12-10T09:00:00Z",
  "UpdatedAt": "2024-12-10T09:00:00Z",
  "StartTime": "2024-12-15T14:00:00Z",
  "EndTime": "2024-12-15T16:00:00Z",
  "Location": "Main Field",
  "UserRSVP": {
    "ID": "rsvp-uuid",
    "EventID": "event-xyz789",
    "UserID": "user-uuid",
    "Response": "yes",
    "CreatedAt": "2024-12-10T11:00:00Z",
    "UpdatedAt": "2024-12-10T11:00:00Z"
  },
  "Metadata": {
    "description": "Regular weekend practice",
    "location": "Main Field"
  }
}
```

#### 3. News Items (Type="news")
Represents news posts from clubs.

**Example:**
```json
{
  "ID": "news-def456",
  "ClubID": "club-uuid",
  "ClubName": "Soccer Club",
  "Type": "news",
  "Title": "Tournament Registration Open",
  "Content": "Registration for the annual tournament is now open...",
  "Timestamp": "2024-12-10T08:00:00Z",
  "CreatedAt": "2024-12-10T08:00:00Z",
  "UpdatedAt": "2024-12-10T08:00:00Z",
  "Metadata": {}
}
```

### Query Examples

#### Get all timeline items
```
GET /api/v2/TimelineItems
```

#### Filter by club
```
GET /api/v2/TimelineItems?$filter=ClubID eq 'club-uuid'
```

#### Filter by type
```
GET /api/v2/TimelineItems?$filter=Type eq 'event'
```

#### Get recent items (top 20)
```
GET /api/v2/TimelineItems?$top=20
```

#### Get items with pagination
```
GET /api/v2/TimelineItems?$top=10&$skip=0
```

#### Select specific fields
```
GET /api/v2/TimelineItems?$select=ID,Title,Type,Timestamp
```

#### Filter by date range
```
GET /api/v2/TimelineItems?$filter=Timestamp ge 2024-12-01T00:00:00Z and Timestamp le 2024-12-31T23:59:59Z
```

### Authorization
- Users can only see timeline items from clubs they are members of
- Authorization is enforced server-side automatically
- No special permissions required beyond club membership

### Performance Considerations
- Timeline items are aggregated server-side for optimal performance
- Results are pre-sorted by timestamp (most recent first)
- Default limit of 50 items per source (activities, events, news)
- Use `$top` and `$skip` for efficient pagination

### Benefits Over Separate Endpoints
1. **Single Request**: One API call instead of three separate calls
2. **Server-Side Aggregation**: Better performance and consistency
3. **Unified Sorting**: All items sorted by timestamp seamlessly
4. **Simplified Client Code**: No need to merge and sort data client-side
5. **OData Query Support**: Full OData query options available

---

## Migration from v1 Dashboard Functions

The Timeline entity replaces the following v1 custom functions:
- `GET /api/v2/GetDashboardNews()`
- `GET /api/v2/GetDashboardEvents()`
- `GET /api/v2/GetDashboardActivities()`

These functions are still available for backward compatibility but are deprecated and will be removed in a future release.

### Migration Example

**Before (v1 - Multiple Calls):**
```typescript
const [newsResponse, eventsResponse, activitiesResponse] = await Promise.all([
  api.get('/api/v2/GetDashboardNews()'),
  api.get('/api/v2/GetDashboardEvents()'),
  api.get('/api/v2/GetDashboardActivities()')
]);

const news = newsResponse.data.value || [];
const events = eventsResponse.data.value || [];
const activities = activitiesResponse.data.value || [];

// Client-side merging and sorting required...
```

**After (v2 - Single Call):**
```typescript
const response = await api.get('/api/v2/TimelineItems');
const timeline = response.data.value || [];

// Already aggregated, sorted, and ready to use!
```

---

## See Also
- [OData v4.01 Specification](https://docs.oasis-open.org/odata/odata/v4.01/odata-v4.01-part1-protocol.html)
- [OData Migration Plan](./OData_Migration_Plan.md)
- [OData Migration Summary](./OData_Migration_Summary.md)
