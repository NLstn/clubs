package core

import (
	"time"
)

// TimelineItem represents a unified timeline entry that can be an activity, event, or news item
// This is a virtual entity that aggregates data from multiple sources
type TimelineItem struct {
	ID        string    `json:"ID" odata:"key"`
	ClubID    string    `json:"ClubID"`
	ClubName  string    `json:"ClubName"`
	Type      string    `json:"Type"` // "activity", "event", "news"
	Title     string    `json:"Title"`
	Content   string    `json:"Content,omitempty"`
	Timestamp time.Time `json:"Timestamp"` // Unified timestamp for sorting
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`

	// Event-specific fields (only populated for Type="event")
	StartTime *time.Time `json:"StartTime,omitempty"`
	EndTime   *time.Time `json:"EndTime,omitempty"`
	Location  *string    `json:"Location,omitempty"`

	// Activity-specific fields (only populated for Type="activity")
	Actor     *string `json:"Actor,omitempty"`
	ActorName *string `json:"ActorName,omitempty"`

	// Metadata for additional information
	Metadata map[string]interface{} `json:"Metadata,omitempty"`

	// RSVP information for events (only for Type="event")
	UserRSVP *EventRSVP `json:"UserRSVP,omitempty"`
}
