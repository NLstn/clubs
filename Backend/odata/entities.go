package odata

import (
	"fmt"

	"github.com/NLstn/clubs/models"
)

// registerEntities registers all entity types with the OData service
func (s *Service) registerEntities() error {
	entities := []interface{}{
		// Core entities
		&models.User{},
		&models.Club{},
		&models.Member{},

		// Team entities
		&models.Team{},

		// Event entities
		&models.Event{},
		&models.EventRSVP{},

		// Shift entities
		&models.Shift{},
		&models.ShiftMember{},

		// Fine entities
		&models.Fine{},
		&models.FineTemplate{},

		// Invite and join request entities
		&models.Invite{},
		&models.JoinRequest{},

		// News and notification entities
		&models.News{},
		&models.Notification{},
		&models.UserNotificationPreferences{},

		// Settings and privacy entities
		&models.ClubSettings{},
		&models.UserPrivacySettings{},
	}

	for _, entity := range entities {
		if err := s.Service.RegisterEntity(entity); err != nil {
			return fmt.Errorf("failed to register entity %T: %w", entity, err)
		}
	}

	// Register virtual entities
	virtualEntities := []interface{}{
		&models.TimelineItem{},
	}

	for _, entity := range virtualEntities {
		if err := s.Service.RegisterVirtualEntity(entity); err != nil {
			return fmt.Errorf("failed to register virtual entity %T: %w", entity, err)
		}
	}

	return nil
}
