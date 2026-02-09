package controller

import (
	"go-login/dto"
	"go-login/models"
	"go-login/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ExchangeRateController struct {
	service service.ExchangeRateService
}

func NewExchangeRateController(service service.ExchangeRateService) *ExchangeRateController {
	return &ExchangeRateController{service: service}
}


func toExchangeRateResponse(r models.ExchangeRate) dto.ExchangeRateResponse {
	resp := dto.ExchangeRateResponse{
		ID:             r.ID,
		FromCurrencyID: r.FromCurrencyID,
		ToCurrencyID:   r.ToCurrencyID,
		Rate:           r.Rate,
		IsActive:       r.IsActive,
		CreatedAt:      r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:      r.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if r.FromCurrency.ID != 0 {
		fc := toCurrencyResponse(r.FromCurrency)
		resp.FromCurrency = &fc
	}
	if r.ToCurrency.ID != 0 {
		tc := toCurrencyResponse(r.ToCurrency)
		resp.ToCurrency = &tc
	}
	return resp
}


func (ctrl *ExchangeRateController) CreateExchangeRate(ctx *gin.Context) {
	var req dto.CreateExchangeRateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	rate, err := ctrl.service.CreateExchangeRate(req.FromCurrencyID, req.ToCurrencyID, req.Rate)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "exchange rate already exists for this currency pair" {
			status = http.StatusConflict
		}
		ctx.JSON(status, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.SuccessResponse{
		Message: "Exchange rate created successfully",
		Data:    toExchangeRateResponse(*rate),
	})
}


func (ctrl *ExchangeRateController) GetExchangeRate(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid exchange rate ID"})
		return
	}

	rate, err := ctrl.service.GetExchangeRate(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, toExchangeRateResponse(*rate))
}


func (ctrl *ExchangeRateController) GetAllExchangeRates(ctx *gin.Context) {
	rates, err := ctrl.service.GetAllActiveExchangeRates()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	result := make([]dto.ExchangeRateResponse, len(rates))
	for i, r := range rates {
		result[i] = toExchangeRateResponse(r)
	}
	ctx.JSON(http.StatusOK, result)
}


func (ctrl *ExchangeRateController) UpdateExchangeRate(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid exchange rate ID"})
		return
	}

	var req dto.UpdateExchangeRateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	rate, err := ctrl.service.UpdateExchangeRate(uint(id), req.Rate, req.IsActive)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "exchange rate not found" {
			status = http.StatusNotFound
		}
		ctx.JSON(status, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Exchange rate updated successfully",
		Data:    toExchangeRateResponse(*rate),
	})
}


func (ctrl *ExchangeRateController) DeleteExchangeRate(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid exchange rate ID"})
		return
	}

	if err := ctrl.service.DeleteExchangeRate(uint(id)); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "exchange rate not found" {
			status = http.StatusNotFound
		}
		ctx.JSON(status, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{Message: "Exchange rate deleted successfully"})
}
