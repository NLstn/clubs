package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/NLstn/clubs/auth"
	"github.com/NLstn/clubs/azure"
	"github.com/NLstn/clubs/database"
	"github.com/NLstn/clubs/handlers"
	"github.com/NLstn/clubs/models"
	"github.com/NLstn/clubs/odata"
	frontend "github.com/NLstn/clubs/tools"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	// ignore error if .env file is not found

	err := database.Init()
	if err != nil {
		log.Fatal("Could not initialize database:", err)
	}

	// FIXME: This should be in the database.go file, but importing the models there would result
	//        in a circular dependency.
	err = database.Db.AutoMigrate(&models.Club{},
		&models.Member{},
		&models.Team{},
		&models.TeamMember{},
		&models.MagicLink{},
		&models.User{},
		&models.JoinRequest{},
		&models.Invite{},
		&models.RefreshToken{},
		&models.Fine{},
		&models.FineTemplate{},
		&models.Shift{},
		&models.ShiftMember{},
		&models.Event{},
		&models.EventRSVP{},
		&models.News{},
		&models.ClubSettings{},
		&models.Notification{},
		&models.UserNotificationPreferences{},
		&models.UserPrivacySettings{},
		&models.MemberPrivacySettings{},
		&models.Activity{},
		&models.APIKey{},
	)
	if err != nil {
		log.Fatal("Could not migrate database:", err)
	}

	// Data migration: Move club-specific privacy settings to MemberPrivacySettings
	err = migratePrivacySettings()
	if err != nil {
		log.Printf("Warning: Could not migrate privacy settings: %v", err)
	}

	err = azure.Init()
	if err != nil {
		log.Fatal("Could not initialize Azure SDK:", err)
	}

	err = auth.Init()
	if err != nil {
		log.Fatal("Could not initialize auth:", err)
	}

	err = auth.InitKeycloak()
	if err != nil {
		log.Printf("Warning: Could not initialize Keycloak: %v", err)
	}

	err = frontend.Init()
	if err != nil {
		log.Fatal("Could not initialize frontend:", err)
	}

	mux := http.NewServeMux()

	// Add health check at root level for container monitoring
	mux.HandleFunc("/health", handlers.HealthCheck)

	mux.Handle("/api/v1/", handlers.Handler_v1())

	// Mount OData v2 API (Phase 2: With authentication and authorization)
	odataService, err := odata.NewService(database.Db)
	if err != nil {
		log.Fatal("Could not initialize OData service:", err)
	}

	// Get JWT secret for OData authentication middleware
	jwtSecret := []byte(auth.GetJWTSecret())

	// Create a submux for /api/v2/ to handle both OData and custom routes
	odataV2Mux := http.NewServeMux()

	// Register custom handlers (e.g., file uploads) that don't fit standard OData patterns
	odataService.RegisterCustomHandlers(odataV2Mux)

	// Register the OData service as the default handler
	odataV2Mux.Handle("/", odataService)

	// Wrap OData v2 service with authentication middleware
	// This enforces JWT token validation on all /api/v2/ endpoints
	// except for metadata and service document endpoints
	odataWithAuth := http.StripPrefix("/api/v2", odata.AuthMiddleware(jwtSecret)(odataV2Mux))
	mux.Handle("/api/v2/", odataWithAuth)

	handler := handlers.CorsMiddleware(mux)
	handlerWithLogging := handlers.LoggingMiddleware(handler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handlerWithLogging))
}

// migratePrivacySettings migrates club-specific privacy settings from UserPrivacySettings to MemberPrivacySettings
func migratePrivacySettings() error {
	// Check if ClubID column exists in user_privacy_settings table
	var columnExists bool
	err := database.Db.Raw(`
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns 
			WHERE table_name = 'user_privacy_settings' AND column_name = 'club_id'
		)
	`).Scan(&columnExists).Error
	if err != nil {
		return fmt.Errorf("failed to check for club_id column: %w", err)
	}

	if !columnExists {
		// Migration already completed or not needed
		log.Println("Privacy settings migration: club_id column not found, skipping migration")
		return nil
	}

	log.Println("Starting privacy settings migration...")

	// Find all club-specific settings (where club_id IS NOT NULL)
	var clubSettings []struct {
		ID             string
		UserID         string
		ClubID         string
		ShareBirthDate bool
		CreatedAt      time.Time
		UpdatedAt      time.Time
	}
	
	err = database.Db.Raw("SELECT id, user_id, club_id, share_birth_date, created_at, updated_at FROM user_privacy_settings WHERE club_id IS NOT NULL").Scan(&clubSettings).Error
	if err != nil {
		return fmt.Errorf("failed to fetch club-specific settings: %w", err)
	}

	log.Printf("Found %d club-specific privacy settings to migrate", len(clubSettings))

	// For each club-specific setting, find the corresponding member and create MemberPrivacySettings
	migratedCount := 0
	for _, setting := range clubSettings {
		var member models.Member
		err := database.Db.Where("user_id = ? AND club_id = ?", setting.UserID, setting.ClubID).First(&member).Error
		if err != nil {
			log.Printf("Warning: Could not find member for user_id=%s, club_id=%s: %v", setting.UserID, setting.ClubID, err)
			continue
		}

		// Check if member privacy setting already exists
		var existingCount int64
		database.Db.Model(&models.MemberPrivacySettings{}).Where("member_id = ?", member.ID).Count(&existingCount)
		if existingCount > 0 {
			log.Printf("Skipping member_id=%s, settings already exist", member.ID)
			continue
		}

		// Create new MemberPrivacySettings
		memberPrivacy := models.MemberPrivacySettings{
			MemberID:       member.ID,
			ShareBirthDate: setting.ShareBirthDate,
			CreatedAt:      setting.CreatedAt,
			UpdatedAt:      setting.UpdatedAt,
		}
		
		err = database.Db.Create(&memberPrivacy).Error
		if err != nil {
			log.Printf("Warning: Could not create MemberPrivacySettings for member_id=%s: %v", member.ID, err)
			continue
		}

		migratedCount++
	}

	log.Printf("Migrated %d club-specific privacy settings to MemberPrivacySettings", migratedCount)

	// Delete migrated club-specific settings from user_privacy_settings
	if migratedCount > 0 {
		result := database.Db.Exec("DELETE FROM user_privacy_settings WHERE club_id IS NOT NULL")
		if result.Error != nil {
			return fmt.Errorf("failed to delete migrated settings: %w", result.Error)
		}
		log.Printf("Deleted %d club-specific settings from UserPrivacySettings", result.RowsAffected)
	}

	// Drop the club_id column from user_privacy_settings table
	err = database.Db.Exec("ALTER TABLE user_privacy_settings DROP COLUMN IF EXISTS club_id").Error
	if err != nil {
		return fmt.Errorf("failed to drop club_id column: %w", err)
	}
	log.Println("Dropped club_id column from user_privacy_settings table")

	// Add unique constraint on user_id if it doesn't exist
	err = database.Db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint 
				WHERE conname = 'idx_user_privacy_settings_user_id'
			) THEN
				CREATE UNIQUE INDEX idx_user_privacy_settings_user_id ON user_privacy_settings(user_id);
			END IF;
		END $$;
	`).Error
	if err != nil {
		return fmt.Errorf("failed to add unique constraint: %w", err)
	}
	log.Println("Added unique constraint on user_id in user_privacy_settings table")

	log.Println("Privacy settings migration completed successfully")
	return nil
}
