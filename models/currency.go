package models

import "time"

type Currency struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Code      string     `gorm:"type:varchar(10);uniqueIndex;not null" json:"code"`
	Name      string     `gorm:"type:varchar(100);not null" json:"name"`
	Symbol    string     `gorm:"type:varchar(10);not null" json:"symbol"`
	IsActive  bool       `gorm:"default:true;not null" json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
