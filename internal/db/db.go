package db

import (
	"fmt"
	"log"

	"github.com/SoliMark/gotasker-pro/internal/model"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := viper.GetString("DB_URL")
	if dsn == "" {
		log.Fatal("DB_URL is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatalf("Failed to migratet : %v", err)
	}

	fmt.Println("Connected to database successfully!")

	DB = db
}
