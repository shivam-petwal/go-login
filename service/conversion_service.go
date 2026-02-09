package service

import (
"errors"
"fmt"
"go-login/repository"
)

type ConversionService interface {
ConvertCurrency(fromCode, toCode string, amount float64) (float64, float64, error)
}

type conversionService struct {
exchangeRateRepo repository.ExchangeRateRepository
currencyRepo     repository.CurrencyRepository
}

func NewConversionService(exchangeRateRepo repository.ExchangeRateRepository, currencyRepo repository.CurrencyRepository) ConversionService {
return &conversionService{
exchangeRateRepo: exchangeRateRepo,
currencyRepo:     currencyRepo,
}
}

func (s *conversionService) ConvertCurrency(fromCode, toCode string, amount float64) (float64, float64, error) {
if amount <= 0 {
return 0, 0, errors.New("amount must be greater than 0")
}


from, err := s.currencyRepo.GetByCode(fromCode)
if err != nil {
return 0, 0, fmt.Errorf("source currency '%s' not found", fromCode)
}
if !from.IsActive {
return 0, 0, fmt.Errorf("source currency '%s' is not active", fromCode)
}


to, err := s.currencyRepo.GetByCode(toCode)
if err != nil {
return 0, 0, fmt.Errorf("target currency '%s' not found", toCode)
}
if !to.IsActive {
return 0, 0, fmt.Errorf("target currency '%s' is not active", toCode)
}


rate, err := s.exchangeRateRepo.GetActiveByCurrencyCodes(fromCode, toCode)
if err != nil {
return 0, 0, fmt.Errorf("exchange rate from '%s' to '%s' not found", fromCode, toCode)
}

return amount * rate.Rate, rate.Rate, nil
}
