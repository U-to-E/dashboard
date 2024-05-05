package models

import "gorm.io/gorm"

type Material struct {
	gorm.Model
	Name     string `gorm:"column:name;not null"`
	FileName string `gorm:"column:filename;not null"`
}
