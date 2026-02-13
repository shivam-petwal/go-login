package repository

import (
	"go-login/models"

	"gorm.io/gorm"
)

type CurrencyRepository interface {
	Create(currency *models.Currency) error
	GetByID(id uint) (*models.Currency, error)
	GetByCode(code string) (*models.Currency, error)
	GetAll() ([]models.Currency, error)
	Update(currency *models.Currency) error
	PartialUpdate(id uint, fields map[string]interface{}) error
}

type currencyRepository struct {
	db *gorm.DB
}

func NewCurrencyRepository(db *gorm.DB) CurrencyRepository {
	return &currencyRepository{db: db}
}

func (r *currencyRepository) Create(currency *models.Currency) error {
	return r.db.Create(currency).Error
}

func (r *currencyRepository) GetByID(id uint) (*models.Currency, error) {
	var currency models.Currency
	if err := r.db.Where("deleted_at IS NULL").First(&currency, id).Error; err != nil {
		return nil, err
	}
	return &currency, nil
}

func (r *currencyRepository) GetByCode(code string) (*models.Currency, error) {
	var currency models.Currency
	if err := r.db.Where("UPPER(code) = UPPER(?) AND deleted_at IS NULL", code).First(&currency).Error; err != nil {
		return nil, err
	}
	return &currency, nil
}

func (r *currencyRepository) GetAll() ([]models.Currency, error) {
	var currencies []models.Currency
	err := r.db.Where("deleted_at IS NULL").Order("code ASC").Find(&currencies).Error
	return currencies, err
}

func (r *currencyRepository) Update(currency *models.Currency) error {
	return r.db.Save(currency).Error
}

func (r *currencyRepository) PartialUpdate(id uint, fields map[string]interface{}) error {
	return r.db.Model(&models.Currency{}).Where("id = ? AND deleted_at IS NULL", id).Updates(fields).Error
}
