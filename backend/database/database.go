package database

import (
	"fmt"

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

	// Auto Migrate the schema
	err = Db.AutoMigrate(&models.Club{}, &models.Member{})
	if err != nil {
		return err
	}

	return nil
}
