package models

import (
	"time"
)

type Transaction struct {
	ID              int       `gorm:"primaryKey;autoIncrement"`
	Sender          *string   `gorm:"null"`
	Recipient       *string   `gorm:"null"`
	Amount          int       `gorm:"not null;check:amount > 0"`
	TransactionType string    `gorm:"not null;size:50"`
	Item            *string   `gorm:"null"`
	CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
