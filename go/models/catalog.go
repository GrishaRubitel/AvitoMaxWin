package models

type CatalogItem struct {
	ID    uint   `gorm:"primaryKey"`
	Item  string `gorm:"not null"`
	Price int    `gorm:"check:price >= 0"`
}
