package controller

import (
"go-login/dto"
"go-login/service"
"net/http"
"strconv"
"strings"

"github.com/gin-gonic/gin"
)

type ConversionController struct {
service service.ConversionService
}

func NewConversionController(service service.ConversionService) *ConversionController {
return &ConversionController{service: service}
}

func (ctrl *ConversionController) ConvertCurrency(ctx *gin.Context) {
from := ctx.Query("from")
to := ctx.Query("to")
amountStr := ctx.Query("amount")

if from == "" || to == "" || amountStr == "" {
ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "from, to, and amount query parameters are required"})
return
}

amount, err := strconv.ParseFloat(amountStr, 64)
if err != nil {
ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "amount must be a valid number"})
return
}

convertedAmount, rate, err := ctrl.service.ConvertCurrency(from, to, amount)
if err != nil {
status := http.StatusBadRequest
if strings.Contains(err.Error(), "not found") {
status = http.StatusNotFound
}
ctx.JSON(status, dto.ErrorResponse{Error: err.Error()})
return
}

ctx.JSON(http.StatusOK, dto.ConversionResponse{
From:            strings.ToUpper(from),
To:              strings.ToUpper(to),
Amount:          amount,
ExchangeRate:    rate,
ConvertedAmount: convertedAmount,
})
}
