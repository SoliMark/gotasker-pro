package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/SoliMark/gotasker-pro/internal/model"
)

func NewDB(dsn string) (*gorm.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("DB_URL is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect DB: %w", err)
	}

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	log.Println("Connected to database successfully!")
	return db, nil
}
