package models

import "gorm.io/gorm"

type Material struct {
	gorm.Model
	Name     string `gorm:"column:name;not null"`
	FilePath string `gorm:"column:filepath;not null"`
}
