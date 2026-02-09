package service

import (
	"errors"
	"go-login/models"
	"go-login/repository"
	"strings"
	"time"

	"gorm.io/gorm"
)

type CurrencyService interface {
	CreateCurrency(code, name, symbol string) (*models.Currency, error)
	GetCurrency(id uint) (*models.Currency, error)
	GetAllCurrencies() ([]models.Currency, error)
	UpdateCurrency(id uint, name, symbol string, isActive *bool) (*models.Currency, error)
	DeleteCurrency(id uint) error
}

type currencyService struct {
	repo repository.CurrencyRepository
}

func NewCurrencyService(repo repository.CurrencyRepository) CurrencyService {
	return &currencyService{repo: repo}
}

func (s *currencyService) CreateCurrency(code, name, symbol string) (*models.Currency, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	if code == "" {
		return nil, errors.New("currency code is required")
	}

	if _, err := s.repo.GetByCode(code); err == nil {
		return nil, errors.New("currency code already exists")
	}

	currency := &models.Currency{
		Code:      code,
		Name:      strings.TrimSpace(name),
		Symbol:    strings.TrimSpace(symbol),
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.repo.Create(currency); err != nil {
		return nil, err
	}
	return currency, nil
}

func (s *currencyService) GetCurrency(id uint) (*models.Currency, error) {
	currency, err := s.repo.GetByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("currency not found")
	}
	return currency, err
}

func (s *currencyService) GetAllCurrencies() ([]models.Currency, error) {
	return s.repo.GetAll()
}

func (s *currencyService) UpdateCurrency(id uint, name, symbol string, isActive *bool) (*models.Currency, error) {
	currency, err := s.GetCurrency(id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		currency.Name = strings.TrimSpace(name)
	}
	if symbol != "" {
		currency.Symbol = strings.TrimSpace(symbol)
	}
	if isActive != nil {
		currency.IsActive = *isActive
	}
	currency.UpdatedAt = time.Now()

	if err := s.repo.Update(currency); err != nil {
		return nil, err
	}
	return currency, nil
}

func (s *currencyService) DeleteCurrency(id uint) error {
	currency, err := s.GetCurrency(id)
	if err != nil {
		return err
	}
	now := time.Now()
	currency.IsActive = false
	currency.DeletedAt = &now
	currency.UpdatedAt = now
	return s.repo.Update(currency)
}
