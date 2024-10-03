package initializers

import (
	"golang/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitializeDB() {
	var err error
	dsn := os.Getenv("DB_URL")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to the Database")
	}

	err = DB.AutoMigrate(&models.User{}, &models.Note{}, &models.CreditCard{})
	if err != nil {
		log.Println("Error during AutoMigrate:", err)
	} else {
		log.Println("Database schema migrated successfully")
	}
}
