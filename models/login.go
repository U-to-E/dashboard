package models

import "gorm.io/gorm"

type Login struct {
	gorm.Model
	Name      string `gorm:"column:name;not null"`
	Email     string `gorm:"column:email;not null;uniqueIndex"`
	Password  []byte `gorm:"not null" json:"password"`
	CollageID string `json:"collageid"`
	Role      string `json:"role"`
}
