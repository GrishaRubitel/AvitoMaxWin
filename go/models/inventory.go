package models

type Inventory struct {
	Login    string `gorm:"primaryKey"`
	Item     string `gorm:"primaryKey"`
	Quantity int    `gorm:"not null;default:1;check:quantity > 0"`
}
