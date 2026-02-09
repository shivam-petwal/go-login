package dto

// --- Auth ---

type RegisterRequest struct {
Username string `json:"username" binding:"required"`
Email    string `json:"email" binding:"required,email"`
Password string `json:"password" binding:"required,min=6"`
}

type RegisterResponse struct {
Message  string `json:"message"`
Username string `json:"username"`
Email    string `json:"email"`
}

type LoginRequest struct {
Email    string `json:"email" binding:"required,email"`
Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
Token string `json:"token"`
}

// --- Currency ---

type CreateCurrencyRequest struct {
Code   string `json:"code" binding:"required"`
Name   string `json:"name" binding:"required"`
Symbol string `json:"symbol" binding:"required"`
}

type UpdateCurrencyRequest struct {
Name     string `json:"name"`
Symbol   string `json:"symbol"`
IsActive *bool  `json:"is_active"`
}

type CurrencyResponse struct {
ID        uint   `json:"id"`
Code      string `json:"code"`
Name      string `json:"name"`
Symbol    string `json:"symbol"`
IsActive  bool   `json:"is_active"`
CreatedAt string `json:"created_at"`
UpdatedAt string `json:"updated_at"`
}

// --- Exchange Rate ---

type CreateExchangeRateRequest struct {
FromCurrencyID uint    `json:"from_currency_id" binding:"required"`
ToCurrencyID   uint    `json:"to_currency_id" binding:"required"`
Rate           float64 `json:"rate" binding:"required,gt=0"`
}

type UpdateExchangeRateRequest struct {
Rate     *float64 `json:"rate" binding:"omitempty,gt=0"`
IsActive *bool    `json:"is_active"`
}

type ExchangeRateResponse struct {
ID             uint              `json:"id"`
FromCurrencyID uint              `json:"from_currency_id"`
ToCurrencyID   uint              `json:"to_currency_id"`
Rate           float64           `json:"rate"`
IsActive       bool              `json:"is_active"`
FromCurrency   *CurrencyResponse `json:"from_currency,omitempty"`
ToCurrency     *CurrencyResponse `json:"to_currency,omitempty"`
CreatedAt      string            `json:"created_at"`
UpdatedAt      string            `json:"updated_at"`
}

// --- Conversion ---

type ConversionResponse struct {
From            string  `json:"from"`
To              string  `json:"to"`
Amount          float64 `json:"amount"`
ExchangeRate    float64 `json:"exchange_rate"`
ConvertedAmount float64 `json:"converted_amount"`
}

// --- Common ---

type ErrorResponse struct {
Error string `json:"error"`
}

type SuccessResponse struct {
Message string      `json:"message"`
Data    interface{} `json:"data,omitempty"`
}
