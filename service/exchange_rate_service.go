package service

import (
	"errors"
	"fmt"
	"go-login/models"
	"go-login/repository"
	"time"

	"gorm.io/gorm"
)

type ExchangeRateService interface {
	CreateExchangeRate(fromCurrencyID, toCurrencyID uint, rate float64) (*models.ExchangeRate, error)
	GetExchangeRate(id uint) (*models.ExchangeRate, error)
	GetAllActiveExchangeRates() ([]models.ExchangeRate, error)
	UpdateExchangeRate(id uint, rate *float64, isActive *bool) (*models.ExchangeRate, error)
	DeleteExchangeRate(id uint) error
}

type exchangeRateService struct {
	repo         repository.ExchangeRateRepository
	currencyRepo repository.CurrencyRepository
}

func NewExchangeRateService(repo repository.ExchangeRateRepository, currencyRepo repository.CurrencyRepository) ExchangeRateService {
	return &exchangeRateService{repo: repo, currencyRepo: currencyRepo}
}

func (s *exchangeRateService) CreateExchangeRate(fromCurrencyID, toCurrencyID uint, rate float64) (*models.ExchangeRate, error) {
	if fromCurrencyID == toCurrencyID {
		return nil, errors.New("from and to currency must be different")
	}
	if rate <= 0 {
		return nil, errors.New("exchange rate must be greater than 0")
	}

	// verification that both currencies exist and are active
	fromCurrency, err := s.currencyRepo.GetByID(fromCurrencyID)
	if err != nil {
		return nil, errors.New("from currency not found")
	}
	if !fromCurrency.IsActive {
		return nil, fmt.Errorf("from currency (%s) is inactive", fromCurrency.Code)
	}

	toCurrency, err := s.currencyRepo.GetByID(toCurrencyID)
	if err != nil {
		return nil, errors.New("to currency not found")
	}
	if !toCurrency.IsActive {
		return nil, fmt.Errorf("to currency (%s) is inactive", toCurrency.Code)
	}

	// Check if pair already exists
	if _, err := s.repo.GetByCurrencyPair(fromCurrencyID, toCurrencyID); err == nil {
		return nil, errors.New("exchange rate already exists for this currency pair")
	}

	exchangeRate := &models.ExchangeRate{
		FromCurrencyID: fromCurrencyID,
		ToCurrencyID:   toCurrencyID,
		Rate:           rate,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := s.repo.Create(exchangeRate); err != nil {
		return nil, err
	}


	exchangeRate.FromCurrency = *fromCurrency
	exchangeRate.ToCurrency = *toCurrency
	return exchangeRate, nil
}

func (s *exchangeRateService) GetExchangeRate(id uint) (*models.ExchangeRate, error) {
	rate, err := s.repo.GetByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("exchange rate not found")
	}
	return rate, err
}

func (s *exchangeRateService) GetAllActiveExchangeRates() ([]models.ExchangeRate, error) {
	return s.repo.GetAllActive()
}

func (s *exchangeRateService) UpdateExchangeRate(id uint, rate *float64, isActive *bool) (*models.ExchangeRate, error) {
	exchangeRate, err := s.GetExchangeRate(id)
	if err != nil {
		return nil, err
	}

	if rate != nil {
		if *rate <= 0 {
			return nil, errors.New("exchange rate must be greater than 0")
		}
		exchangeRate.Rate = *rate
	}
	if isActive != nil {
		exchangeRate.IsActive = *isActive
	}
	exchangeRate.UpdatedAt = time.Now()

	if err := s.repo.Update(exchangeRate); err != nil {
		return nil, err
	}
	return exchangeRate, nil
}

func (s *exchangeRateService) DeleteExchangeRate(id uint) error {
	exchangeRate, err := s.GetExchangeRate(id)
	if err != nil {
		return err
	}
	now := time.Now()
	exchangeRate.IsActive = false
	exchangeRate.DeletedAt = &now
	exchangeRate.UpdatedAt = now
	return s.repo.Update(exchangeRate)
}
