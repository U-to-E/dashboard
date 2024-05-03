package models

import (
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Username  string `gorm:"not null" json:"username"`
	Email     string `gorm:"uniqueIndexnot null" json:"email"`
	Password  []byte `gorm:"not null" json:"password"`
	CollageID int    `gorm:"not null" json:"collageid"`
	Role      string `gorm:"not null" json:"role"`
}
