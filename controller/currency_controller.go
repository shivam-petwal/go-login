package controller

import (
	"encoding/json"
	"go-login/dto"
	"go-login/models"
	"go-login/service"
	"go-login/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CurrencyController struct {
	service service.CurrencyService
}

func NewCurrencyController(service service.CurrencyService) *CurrencyController {
	return &CurrencyController{service: service}
}

func toCurrencyResponse(c models.Currency) dto.CurrencyResponse {
	return dto.CurrencyResponse{
		ID:        c.ID,
		Code:      c.Code,
		Name:      c.Name,
		Symbol:    c.Symbol,
		IsActive:  c.IsActive,
		CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: c.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func (ctrl *CurrencyController) CreateCurrency(ctx *gin.Context) {
	var req dto.CreateCurrencyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	currency, err := ctrl.service.CreateCurrency(req.Code, req.Name, req.Symbol)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "currency code already exists" {
			status = http.StatusConflict
		}
		ctx.JSON(status, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.SuccessResponse{
		Message: "Currency created successfully",
		Data:    toCurrencyResponse(*currency),
	})
}

func (ctrl *CurrencyController) GetCurrency(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid currency ID"})
		return
	}

	currency, err := ctrl.service.GetCurrency(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, toCurrencyResponse(*currency))
}

func (ctrl *CurrencyController) GetAllCurrencies(ctx *gin.Context) {
	currencies, err := ctrl.service.GetAllCurrencies()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	result := make([]dto.CurrencyResponse, len(currencies))
	for i, c := range currencies {
		result[i] = toCurrencyResponse(c)
	}
	ctx.JSON(http.StatusOK, result)
}

func (ctrl *CurrencyController) UpdateCurrency(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid currency ID"})
		return
	}

	var raw map[string]interface{}
	if err := ctx.ShouldBindJSON(&raw); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	if err := utils.ValidateNoUnknownFields(raw, dto.UpdateCurrencyRequest{}); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	var req dto.UpdateCurrencyRequest
	jsonBytes, _ := json.Marshal(raw)
	json.Unmarshal(jsonBytes, &req)

	currency, err := ctrl.service.UpdateCurrency(uint(id), req)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "currency not found" {
			status = http.StatusNotFound
		}
		ctx.JSON(status, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Currency updated successfully",
		Data:    toCurrencyResponse(*currency),
	})
}

func (ctrl *CurrencyController) DeleteCurrency(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid currency ID"})
		return
	}

	if err := ctrl.service.DeleteCurrency(uint(id)); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "currency not found" {
			status = http.StatusNotFound
		}
		ctx.JSON(status, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{Message: "Currency deleted successfully"})
}
