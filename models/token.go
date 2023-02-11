package models

import "time"

type Token struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time
	Value     string `json:"value"`
	UserID    uint   `json:"userId"`
}
