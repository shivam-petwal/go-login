package models

import "time"

type ExchangeRate struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	FromCurrencyID uint       `gorm:"not null" json:"from_currency_id"`
	ToCurrencyID   uint       `gorm:"not null" json:"to_currency_id"`
	Rate           float64    `gorm:"type:decimal(18,6);not null" json:"rate"`
	IsActive       bool       `gorm:"default:true;not null" json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	FromCurrency Currency `gorm:"foreignKey:FromCurrencyID" json:"from_currency"`
	ToCurrency   Currency `gorm:"foreignKey:ToCurrencyID" json:"to_currency"`
}
