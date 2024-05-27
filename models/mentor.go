package models

import "gorm.io/gorm"

type Mentor struct {
	gorm.Model
	Name            string `gorm:"column:name;not null"`
	PhoneNumber     string `gorm:"column:phnumber;not null"`
	Email           string `gorm:"column:email;not null"`
	AssignedCollage string `gorm:"column:collage_id;"`
}
