# Recurring Events Feature

## Overview

This feature allows club admins to create recurring events (daily, weekly, or monthly) with customizable intervals and end dates.

## Backend Implementation

### Model Changes
- Added recurring event fields to `Event` struct:
  - `IsRecurring bool` - Whether this is a recurring event template
  - `RecurrencePattern string` - "daily", "weekly", or "monthly"
  - `RecurrenceInterval int` - Every N intervals (e.g., every 2 weeks)
  - `RecurrenceEnd *time.Time` - When to stop creating events
  - `ParentEventID *string` - Links event instances to parent template

### API Endpoint
- `POST /api/v1/clubs/{clubid}/events/recurring` - Creates recurring events
- Generates individual event instances for each occurrence
- Returns count of created events and event details

### Database Schema
Individual event instances are stored as separate records, linked via `ParentEventID` for easy querying and management.

## Frontend Implementation

### UI Changes
- Added recurring events checkbox to AddEvent modal
- Dropdown for recurrence pattern (daily/weekly/monthly)
- Number input for interval (every N days/weeks/months)
- Date picker for recurrence end date
- Conditional validation for recurring event fields

### API Integration
- Uses different endpoint when recurring option is selected
- Enhanced form validation for recurring event parameters
- Maintains backward compatibility with existing event creation

## Usage Examples

### Weekly Team Meeting
```json
{
  "name": "Weekly Team Meeting",
  "description": "Every Monday team meeting",
  "location": "Conference Room A",
  "start_time": "2025-01-06T10:00:00Z",
  "end_time": "2025-01-06T11:00:00Z",
  "recurrence_pattern": "weekly",
  "recurrence_interval": 1,
  "recurrence_end": "2025-03-31T10:00:00Z"
}
```

This creates:
- 1 parent event (marked as recurring)
- ~12 individual event instances (one per week until end date)

### Daily Standup
```json
{
  "name": "Daily Standup",
  "description": "Daily team standup",
  "location": "Online",
  "start_time": "2025-01-06T09:00:00Z",
  "end_time": "2025-01-06T09:30:00Z",
  "recurrence_pattern": "daily",
  "recurrence_interval": 1,
  "recurrence_end": "2025-01-12T09:30:00Z"
}
```

This creates:
- 1 parent event
- 6 individual event instances (one per day for 6 days)

## Test Coverage

- ✅ Unit tests for `CreateRecurringEvent` function
- ✅ Handler tests for API endpoint validation
- ✅ Integration tests for end-to-end functionality
- ✅ Frontend tests for existing functionality
- ✅ All existing tests continue to pass

## Benefits

1. **Flexible Scheduling**: Supports daily, weekly, and monthly patterns
2. **Customizable Intervals**: Every N days/weeks/months
3. **Individual Management**: Each occurrence is a separate event for easy editing/deletion
4. **Backward Compatibility**: Existing event functionality unchanged
5. **Admin Control**: Only club owners and admins can create recurring events
6. **User-Friendly**: Simple checkbox-based UI with clear validation

## Technical Details

- Events are generated at creation time, not dynamically
- Each occurrence can be individually modified or deleted
- Parent event serves as template and metadata store
- Database queries remain efficient with existing indexes
- No breaking changes to existing event API