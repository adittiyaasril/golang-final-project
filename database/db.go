package database

import (
	"final-project/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func StartDB() {
	dsn := "host=localhost user=postgres password=230798 dbname=my-gram-test port=5432 sslmode=disable TimeZone=Asia/Jakarta"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("error connecting to database : ", err)
	}

	err = DB.Debug().AutoMigrate(&models.User{}, &models.Photo{}, &models.Comment{}, &models.SocialMedia{})
	if err != nil {
		log.Fatal("error migrating database : ", err)
	}

	log.Println("Database connection successfully opened")
}
