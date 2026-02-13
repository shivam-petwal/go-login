package main

import (
	"go-login/config"
	"go-login/controller"

	//"go-login/models"
	"go-login/repository"
	"go-login/router"
	"go-login/service"
	"log"
	"os"

	"github.com/robfig/cron/v3"
)

func main() {
	config.ConnectDB()
	//config.DBcon.AutoMigrate(&models.User{}, &models.Currency{}, &models.ExchangeRate{})

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret-key-change"
	}

	// Repositories
	userRepo := repository.NewUserRepository(config.DBcon)
	currencyRepo := repository.NewCurrencyRepository(config.DBcon)
	exchangeRateRepo := repository.NewExchangeRateRepository(config.DBcon)

	// Services
	authService := service.NewAuthService(userRepo, jwtSecret)
	currencyService := service.NewCurrencyService(currencyRepo)
	exchangeRateService := service.NewExchangeRateService(exchangeRateRepo, currencyRepo)
	conversionService := service.NewConversionService(exchangeRateRepo, currencyRepo)

	// Sync exchange rates
	// service.NewRateSyncService(exchangeRateRepo).SyncAll()
	c := cron.New()
	if _, err := c.AddFunc("0 0,12 * * *", service.NewRateSyncService(exchangeRateRepo).SyncAll); err != nil {
		log.Fatalf("failed to schedule rate sync cron: %v", err)
	}
	c.Start()

	// Controllers
	authCtrl := controller.NewAuthController(authService)
	currencyCtrl := controller.NewCurrencyController(currencyService)
	exchangeRateCtrl := controller.NewExchangeRateController(exchangeRateService)
	conversionCtrl := controller.NewConversionController(conversionService)

	// Router
	r := router.SetupRouter(authCtrl, currencyCtrl, exchangeRateCtrl, conversionCtrl, jwtSecret)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
