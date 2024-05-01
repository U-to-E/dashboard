package database

import (
	"github.com/U-to-E/dashboard/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	con := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(con), &gorm.Config{})

	if err != nil {
		panic("could not connect to the database")
	}

	DB = db

	db.AutoMigrate(&models.Student{})

}
