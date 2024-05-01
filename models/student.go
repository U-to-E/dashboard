package models

type Student struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	Email    string
	Password []byte
}
