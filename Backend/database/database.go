package database

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

var Db *gorm.DB

func NewConnection(config *Config) error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	var err error
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	log.Printf("Successfully connected to database at %s:%d (database: %s)", config.Host, config.Port, config.DBName)

	Db.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\"")

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

	dbName := os.Getenv("DATABASE_NAME")
	if dbName == "" {
		dbName = "civo"
	}

	sslMode := os.Getenv("DATABASE_SSL_MODE")
	if sslMode == "" {
		sslMode = "disable"
	}

	config := &Config{
		Host:     dbUrl,
		Port:     dbPortInt,
		User:     dbUser,
		Password: dbUserPassword,
		DBName:   dbName,
		SSLMode:  sslMode,
	}

	return NewConnection(config)
}
