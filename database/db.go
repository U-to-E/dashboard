package database

import (
	"fmt"

	"github.com/U-to-E/dashboard/config"
	"github.com/U-to-E/dashboard/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	con := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", config.Config("DB_HOST"), config.Config("DB_USER"), config.Config("DB_PASSWORD"), config.Config("DB_NAME"), config.Config("DB_PORT"))
	db, err := gorm.Open(postgres.Open(con), &gorm.Config{})

	if err != nil {
		panic("could not connect to the database")
	}

	DB = db

	fmt.Println("Connection Opened to Database")
	db.AutoMigrate(&models.Student{}, &models.Login{}, &models.Mentor{}, &models.Mapping{}, &models.Quiz{}, &models.Marks{}, &models.PasswordResetToken{}, &models.QuizState{})
	fmt.Println("Database Migrated")

}
