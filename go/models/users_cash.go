package models

type UserCash struct {
	Login string `gorm:"primaryKey"`
	Cash  int    `gorm:"default:1000;check:cash >= 0"`
}
