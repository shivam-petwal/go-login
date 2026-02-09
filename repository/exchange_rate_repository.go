package repository

import (
	"go-login/models"

	"gorm.io/gorm"
)

type ExchangeRateRepository interface {
	Create(rate *models.ExchangeRate) error
	GetByID(id uint) (*models.ExchangeRate, error)
	GetAllActive() ([]models.ExchangeRate, error)
	GetByCurrencyPair(fromID, toID uint) (*models.ExchangeRate, error)
	GetActiveByCurrencyCodes(fromCode, toCode string) (*models.ExchangeRate, error)
	Update(rate *models.ExchangeRate) error
}

type exchangeRateRepository struct {
	db *gorm.DB
}

func NewExchangeRateRepository(db *gorm.DB) ExchangeRateRepository {
	return &exchangeRateRepository{db: db}
}

func (r *exchangeRateRepository) Create(rate *models.ExchangeRate) error {
	return r.db.Create(rate).Error
}

func (r *exchangeRateRepository) GetByID(id uint) (*models.ExchangeRate, error) {
	var rate models.ExchangeRate
	if err := r.db.Preload("FromCurrency").Preload("ToCurrency").
		Where("deleted_at IS NULL").First(&rate, id).Error; err != nil {
		return nil, err
	}
	return &rate, nil
}

func (r *exchangeRateRepository) GetAllActive() ([]models.ExchangeRate, error) {
	var rates []models.ExchangeRate
	err := r.db.Preload("FromCurrency").Preload("ToCurrency").
		Where("is_active = ? AND deleted_at IS NULL", true).Order("created_at DESC").Find(&rates).Error
	return rates, err
}

func (r *exchangeRateRepository) GetByCurrencyPair(fromID, toID uint) (*models.ExchangeRate, error) {
	var rate models.ExchangeRate
	if err := r.db.Where("from_currency_id = ? AND to_currency_id = ? AND deleted_at IS NULL", fromID, toID).First(&rate).Error; err != nil {
		return nil, err
	}
	return &rate, nil
}

func (r *exchangeRateRepository) GetActiveByCurrencyCodes(fromCode, toCode string) (*models.ExchangeRate, error) {
	var rate models.ExchangeRate
	err := r.db.Preload("FromCurrency").Preload("ToCurrency").
		Joins("JOIN currencies AS fc ON fc.id = exchange_rates.from_currency_id").
		Joins("JOIN currencies AS tc ON tc.id = exchange_rates.to_currency_id").
		Where("UPPER(fc.code) = UPPER(?) AND UPPER(tc.code) = UPPER(?) AND exchange_rates.is_active = ? AND exchange_rates.deleted_at IS NULL",
			fromCode, toCode, true).
		First(&rate).Error
	if err != nil {
		return nil, err
	}
	return &rate, nil
}

func (r *exchangeRateRepository) Update(rate *models.ExchangeRate) error {
	return r.db.Save(rate).Error
}
