package models

import "gorm.io/gorm"

type Mappings struct {
	gorm.Model
	MentorID   uint
	CollegeID  uint
	StudentsID []uint `gorm:"type:integer[]"`
}
