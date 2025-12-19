# API Usage Improvement Findings

This document identifies opportunities to improve the frontend's usage of backend APIs by leveraging OData features, adding navigation properties, and refactoring frontend components.

## Table of Contents
1. [N+1 Query Pattern Issues](#n1-query-pattern-issues)
2. [Missing Navigation Properties](#missing-navigation-properties)
3. [Underutilized OData Features](#underutilized-odata-features)
4. [Frontend-Backend Alignment Issues](#frontend-backend-alignment-issues)
5. [Suggested Backend API Additions](#suggested-backend-api-additions)

---

## N+1 Query Pattern Issues

### 1. AdminClubEventList.tsx - Event RSVPs and Shifts

**Location:** [Frontend/src/pages/clubs/admin/events/AdminClubEventList.tsx](../Frontend/src/pages/clubs/admin/events/AdminClubEventList.tsx)

**Status:** ✅ **RESOLVED** (December 19, 2024)

**Previous Implementation:**
```typescript
// Fetches all events, then loops through each to fetch RSVPs and Shifts
const response = await api.get(`/api/v2/Events?$filter=ClubID eq '${id}'`);
for (const event of response.data || []) {
    const rsvpResponse = await api.get(`/api/v2/Events('${event.id}')/EventRSVPs`);
    const shiftsResponse = await api.get(`/api/v2/Shifts?$filter=EventID eq '${event.id}'`);
}
```

**Problem:** For N events, this made 2N+1 API calls.

**Current Implementation:**
```typescript
// Single query with nested expansions
const response = await api.get(`/api/v2/Events?$filter=ClubID eq '${id}'&$expand=EventRSVPs,Shifts`);
```

**Resolution Notes:**
- Navigation properties (EventRSVPs and Shifts) were already present in the Event model
- Updated component to use `$expand` parameter for single-call data fetching
- Migrated from snake_case to PascalCase field names for OData v2 consistency
- Eliminated N+1 query pattern, reducing API calls from 2N+1 to 1

---

### 2. AdminEventDetails.tsx - Similar N+1 Pattern

**Location:** [Frontend/src/pages/clubs/admin/events/AdminEventDetails.tsx](../Frontend/src/pages/clubs/admin/events/AdminEventDetails.tsx)

**Status:** ✅ **RESOLVED** (December 19, 2024)

**Previous Implementation:**
```typescript
const eventResponse = await api.get(`/api/v2/Events('${eventId}')`);
const rsvpResponse = await api.get(`/api/v2/Events('${eventId}')/EventRSVPs`);
const shiftsResponse = await api.get(`/api/v2/Shifts?$filter=EventID eq '${eventId}'`);
```

**Current Implementation:**
```typescript
// Single call with $expand
const response = await api.get(`/api/v2/Events('${eventId}')?$expand=EventRSVPs($expand=User),Shifts`);
```

**Resolution Notes:**
- Component already had `$expand` implemented for RSVPs
- Updated all field references from snake_case to PascalCase for OData v2 consistency
- Reduced API calls from 3 to 1 per event details page load

---

### 3. MyTeams.tsx - Two-Step Query

**Location:** [Frontend/src/pages/clubs/MyTeams.tsx](../Frontend/src/pages/clubs/MyTeams.tsx)

**Current Implementation:**
```typescript
// Step 1: Get TeamMembers for user
const teamMembersResponse = await api.get(`/api/v2/TeamMembers?$filter=UserID eq '${userId}'`);
// Step 2: Get Teams by IDs with complex OR filter
const teamIdFilter = teamIds.map(id => `ID eq '${id}'`).join(' or ');
const teamsResponse = await api.get(`/api/v2/Teams?$filter=ClubID eq '${clubId}' and (${teamIdFilter})`);
```

**Problem:** Two API calls when one should suffice. Also builds complex OR filter.

**Solution:**
```typescript
// Single call with $expand from TeamMembers
const response = await api.get(
    `/api/v2/TeamMembers?$filter=UserID eq '${userId}'&$expand=Team($filter=ClubID eq '${clubId}')`
);
```

**Backend Enhancement:**
- Add `Team` navigation property to `TeamMember` model (may already exist, needs verification)

---

## Missing Navigation Properties

### 1. Shift Model - Missing Navigation Properties

**Location:** [Backend/models/shift_schedules.go](../Backend/models/shift_schedules.go)

**Status:** ✅ **VERIFIED** (December 19, 2024)

**Current State:**
The `Shift` and `ShiftMember` models have all required OData navigation properties:

```go
// Shift has:
Event        *Event        `gorm:"foreignKey:EventID" json:"Event,omitempty" odata:"nav"`
Club         *Club         `gorm:"foreignKey:ClubID" json:"Club,omitempty" odata:"nav"`
ShiftMembers []ShiftMember `gorm:"foreignKey:ShiftID" json:"ShiftMembers,omitempty" odata:"nav"`

// ShiftMember has:
Shift *Shift `gorm:"foreignKey:ShiftID" json:"Shift,omitempty" odata:"nav"`
User  *User  `gorm:"foreignKey:UserID" json:"User,omitempty" odata:"nav"`
```

**Impact:** 
- Navigation properties already implemented and registered in OData service
- [ProfileShifts.tsx](../Frontend/src/pages/profile/ProfileShifts.tsx) can use `$expand` for nested data
- AdminClubEventList now successfully fetches shifts with events in single call

---

### 2. Event Model - Missing Shifts Navigation

**Location:** [Backend/models/events.go](../Backend/models/events.go)

**Status:** ✅ **VERIFIED** (December 19, 2024)

**Current State:** Event has both `EventRSVPs` and `Shifts` navigation properties:

```go
// Navigation properties
EventRSVPs []EventRSVP `gorm:"foreignKey:EventID" json:"EventRSVPs,omitempty" odata:"nav"`
Shifts     []Shift     `gorm:"foreignKey:EventID" json:"Shifts,omitempty" odata:"nav"`
```

**Impact:** Frontend can fetch events with their shifts in a single `$expand` call (already implemented in AdminClubEventList).

---

### 3. ClubSettings - Could Use Club Navigation

**Location:** [Backend/models/club_settings.go](../Backend/models/club_settings.go)

**Opportunity:**
Instead of:
```typescript
api.get(`/api/v2/ClubSettings?$filter=ClubID eq '${clubId}'`)
```

Could use:
```typescript
api.get(`/api/v2/Clubs('${clubId}')/Settings`)
```

**Backend Change:** Add `Settings ClubSettings` navigation to Club model.

---

## Underutilized OData Features

### 1. $select for Reduced Payload

**Locations with Opportunity:**

| File | Current Query | Suggested Improvement |
|------|--------------|----------------------|
| ClubDetails.tsx | `api.get(/api/v2/Clubs('${id}'))` | Add `$select=ID,Name,Description,LogoURL,Deleted` |
| AdminTeamDetails.tsx | `api.get(/api/v2/Teams('${teamId}')/Events)` | Add `$select=ID,Name,StartTime,EndTime` |
| AdminTeamDetails.tsx | `api.get(/api/v2/Teams('${teamId}')/Fines?$expand=User)` | Add `$select=ID,Reason,Amount,Paid,UserID` |

---

### 2. $count for Pagination

**Current Pattern (Multiple Locations):**
```typescript
// Fetch all items to count them
const response = await api.get('/api/v2/SomeEntity');
const total = response.data.length;
```

**Better Pattern:**
```typescript
// Use $count for total without fetching all records
const countResponse = await api.get('/api/v2/SomeEntity/$count');
const total = parseInt(countResponse.data);
```

**Files to Update:**
- AdminClubEventList.tsx (event counts)
- Team member counts

---

### 3. OData Batch Requests

**Opportunity:** Multiple independent API calls could be batched.

**Example in ClubDetails.tsx:**
```typescript
// Current: Two separate calls
const [clubResponse, adminResponse] = await Promise.all([
    api.get(`/api/v2/Clubs('${id}')`),
    api.get(`/api/v2/Clubs('${id}')/IsAdmin()`)
]);
```

**Could Become:** Single OData batch request (requires backend batch support).

---

## Frontend-Backend Alignment Issues

### 1. Case Transformation Overhead

**Issue:** Frontend often transforms PascalCase (OData) to snake_case (local state).

**Example from EventDetails.tsx:**
```typescript
setEventData({
    id: event.ID,
    name: event.Name,
    description: event.Description,
    // ... 10+ more transformations
});
```

**Recommendation:** 
- Define TypeScript interfaces matching OData PascalCase
- Avoid transformation layer
- Update component props to use PascalCase types

---

### 2. Inconsistent RSVP Count Calculation

**Current Pattern:**
```typescript
// Client-side: Fetch all RSVPs, count manually
const rsvpList = parseODataCollection(rsvpResponse.data);
const computedCounts = calculateRSVPCounts(rsvpList);
```

**Recommendation:** Add server-side function:
```go
// GET /api/v2/Events('{id}')/GetRSVPCounts()
// Returns: { "Yes": 10, "No": 3, "Maybe": 5 }
```

---

### 3. ClubSettings Default Handling

**Location:** [Frontend/src/hooks/useClubSettings.ts](../Frontend/src/hooks/useClubSettings.ts)

**Current Issue:**
```typescript
// Falls back to hardcoded defaults if settings don't exist
catch (err) {
    setSettings({
        FinesEnabled: true,
        ShiftsEnabled: true,
        // etc.
    });
}
```

**Recommendation:** 
- Backend should auto-create ClubSettings when Club is created
- Or add a `GetOrCreateSettings()` function on Club entity

---

## Suggested Backend API Additions

### 1. GetRSVPCounts Bound Function

**Purpose:** Avoid fetching all RSVPs just to count them.

```go
// Register in odata/functions.go
{
    Name:       "GetRSVPCounts",
    IsBound:    true,
    EntitySet:  "Events",
    ReturnType: reflect.TypeOf(map[string]int64{}),
    Handler:    s.getRSVPCountsFunction,
}

// Handler returns: {"Yes": 10, "No": 3, "Maybe": 5}
```

---

### 2. GetMyTeams Bound Function on Club

**Purpose:** Simplify MyTeams.tsx query pattern.

```go
// GET /api/v2/Clubs('{clubId}')/GetMyTeams()
// Returns teams where current user is a member
```

---

### 3. GetEventWithDetails Bound Function

**Purpose:** Return event with RSVP counts and shifts in single call.

```go
// GET /api/v2/Events('{id}')/GetDetails()
// Returns:
// {
//   "Event": {...},
//   "RSVPCounts": {"Yes": 10, "No": 3},
//   "Shifts": [...],
//   "UserRSVP": {...}
// }
```

---

### 4. Event Summary for Lists

**Purpose:** Efficient list view with aggregated data.

```go
// GET /api/v2/Clubs('{clubId}')/GetEventSummaries()
// Returns events with pre-calculated RSVP counts and shift counts
```

---

## Priority Matrix

| Issue | Impact | Effort | Priority |
|-------|--------|--------|----------|
| AdminClubEventList N+1 | High (performance) | Medium | **P1** |
| Missing Shift navigation props | High (enables fixes) | Low | **P1** |
| Event Shifts navigation | Medium | Low | **P2** |
| GetRSVPCounts function | Medium | Low | **P2** |
| MyTeams two-step query | Low | Low | **P3** |
| Case transformation | Low | High | **P4** |
| Batch request support | Medium | High | **P4** |

---

## Implementation Checklist

### Backend Changes
- [x] ~~Add navigation properties to Shift model~~ (Already implemented)
- [x] ~~Add navigation properties to ShiftMember model~~ (Already implemented)
- [x] ~~Add Shifts navigation to Event model~~ (Already implemented)
- [ ] Add Settings navigation to Club model
- [ ] Implement GetRSVPCounts bound function
- [ ] Implement GetMyTeams bound function
- [ ] Consider GetEventDetails bound function

### Frontend Changes
- [x] Update AdminClubEventList to use $expand (Completed December 19, 2024)
- [x] Update AdminEventDetails to use $expand (Completed December 19, 2024)
- [ ] Simplify MyTeams.tsx with better OData query
- [ ] Add $select to reduce payload sizes
- [ ] Consider adopting PascalCase interfaces

---

## Related Documentation

- [OData v2 API Documentation](Backend/API.md)
- [Frontend Design System](Frontend/README.md)
- [Adding New Tables Guide](Backend/AddNewTable.md)
