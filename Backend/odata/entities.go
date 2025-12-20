package odata

import (
	"fmt"

	"github.com/NLstn/clubs/models/core"
	"github.com/NLstn/clubs/models/auth"
)

// registerEntities registers all entity types with the OData service
func (s *Service) registerEntities() error {
	entities := []interface{}{
		// Core entities
		&core.User{},
		&core.UserSession{},
		&core.Club{},
		&core.Member{},

		// Team entities
		&core.Team{},
		&core.TeamMember{},

		// Event entities
		&core.Event{},
		&core.EventRSVP{},

		// Shift entities
		&core.Shift{},
		&core.ShiftMember{},

		// Fine entities
		&core.Fine{},
		&core.FineTemplate{},

		// Invite and join request entities
		&core.Invite{},
		&core.JoinRequest{},

		// News and notification entities
		&core.News{},
		&core.Notification{},
		&core.UserNotificationPreferences{},

		// Settings and privacy entities
		&core.ClubSettings{},
		&core.UserPrivacySettings{},
		&core.MemberPrivacySettings{},

		// API Key entities
		&auth.APIKey{},
	}

	for _, entity := range entities {
		if err := s.Service.RegisterEntity(entity); err != nil {
			return fmt.Errorf("failed to register entity %T: %w", entity, err)
		}
	}

	// Register virtual entities
	virtualEntities := []interface{}{
		&core.TimelineItem{},
	}

	for _, entity := range virtualEntities {
		if err := s.Service.RegisterVirtualEntity(entity); err != nil {
			return fmt.Errorf("failed to register virtual entity %T: %w", entity, err)
		}
	}

	return nil
}
