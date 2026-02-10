package router

import (
	"go-login/controller"
	"go-login/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authCtrl *controller.AuthController,
	currencyCtrl *controller.CurrencyController,
	exchangeRateCtrl *controller.ExchangeRateController,
	conversionCtrl *controller.ConversionController,
	jwtSecret string,
) *gin.Engine {
	r := gin.Default()


	r.POST("/register", authCtrl.Register)
	r.POST("/login", authCtrl.Login)

	// equire  JWT
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtSecret))
	{
		// Currency
		protected.POST("/currencies", currencyCtrl.CreateCurrency)
		protected.GET("/currencies", currencyCtrl.GetAllCurrencies)
		protected.GET("/currencies/:id", currencyCtrl.GetCurrency)
		protected.PUT("/currencies/:id", currencyCtrl.UpdateCurrency)
		protected.DELETE("/currencies/:id", currencyCtrl.DeleteCurrency)

		// Exchange Rate 
		protected.POST("/exchange-rates", exchangeRateCtrl.CreateExchangeRate)
		protected.GET("/exchange-rates", exchangeRateCtrl.GetAllExchangeRates)
		protected.GET("/exchange-rates/:id", exchangeRateCtrl.GetExchangeRate)
		protected.PUT("/exchange-rates/:id", exchangeRateCtrl.UpdateExchangeRate)
		protected.DELETE("/exchange-rates/:id", exchangeRateCtrl.DeleteExchangeRate)

		// Conversion
		protected.GET("/convert", conversionCtrl.ConvertCurrency)
	}

	return r
}
