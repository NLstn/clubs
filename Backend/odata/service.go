package odata

import (
	"fmt"
	"log/slog"

	"github.com/nlstn/go-odata"
	"gorm.io/gorm"
)

// Service wraps the go-odata service with our configuration
type Service struct {
	*odata.Service
	db     *gorm.DB
	logger *slog.Logger
}

// NewService creates a new OData service instance with the ClubsService namespace
func NewService(db *gorm.DB) (*Service, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is required")
	}

	// Create OData service with default configuration
	odataService, err := odata.NewServiceWithConfig(db, odata.ServiceConfig{
		PersistentChangeTracking: false, // Can enable later if needed for delta queries
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create OData service: %w", err)
	}

	// Set custom namespace for our service
	if err := odataService.SetNamespace("ClubsService"); err != nil {
		return nil, fmt.Errorf("failed to set namespace: %w", err)
	}

	// Create logger
	logger := slog.Default()
	odataService.SetLogger(logger)

	service := &Service{
		Service: odataService,
		db:      db,
		logger:  logger,
	}

	// Register all entities
	if err := service.registerEntities(); err != nil {
		return nil, fmt.Errorf("failed to register entities: %w", err)
	}

	// Register custom actions (Phase 4)
	if err := service.registerActions(); err != nil {
		return nil, fmt.Errorf("failed to register actions: %w", err)
	}

	// Register custom functions (Phase 4)
	if err := service.registerFunctions(); err != nil {
		return nil, fmt.Errorf("failed to register functions: %w", err)
	}

	// Register virtual entity handlers
	if err := service.registerTimelineHandlers(); err != nil {
		return nil, fmt.Errorf("failed to register timeline handlers: %w", err)
	}

	return service, nil
}
