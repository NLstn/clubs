package database

import (
	"fmt"
	"os"
	"strconv"

	"github.com/NLstn/clubs/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

var Db *gorm.DB

func NewConnection(config *Config) error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)

	var err error
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	Db.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\"")

	// Auto Migrate the schema
	err = Db.AutoMigrate(&models.Club{}, &models.Member{}, &models.Event{}, &models.MagicLink{}, &models.User{})
	if err != nil {
		return err
	}

	return nil
}

func Init() error {
	dbUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		return fmt.Errorf("DATABASE_URL environment variable is required")
	}

	dbPort, ok := os.LookupEnv("DATABASE_PORT")
	if !ok {
		return fmt.Errorf("DATABASE_PORT environment variable is required")
	}

	dbPortInt, err := strconv.Atoi(dbPort)
	if err != nil {
		return fmt.Errorf("DATABASE_PORT must be an integer")
	}

	dbUser := os.Getenv("DATABASE_USER")
	if dbUser == "" {
		return fmt.Errorf("DATABASE_USER environment variable is required")
	}

	dbUserPassword := os.Getenv("DATABASE_USER_PASSWORD")
	if dbUserPassword == "" {
		return fmt.Errorf("DATABASE_USER_PASSWORD environment variable is required")
	}

	config := &Config{
		Host:     dbUrl,
		Port:     dbPortInt,
		User:     dbUser,
		Password: dbUserPassword,
		DBName:   "clubs",
	}

	return NewConnection(config)
}
