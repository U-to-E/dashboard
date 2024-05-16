package models

type Student struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null" json:"username"`
	Email     string `gorm:"uniqueIndexnot null" json:"email"`
	CollageID string `gorm:"not null" json:"collageid"`
	MentorID  string `gorm:"not null" json:"mentorid"`
	Level     int    `gorm:"column:level"`
	Marks     int    `gorm:"column:level"`
}
