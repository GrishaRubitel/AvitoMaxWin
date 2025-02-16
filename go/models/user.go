package models

type User struct {
	Login    string `gorm:"primaryKey"`
	PassHash string `gorm:"not null"`
}
