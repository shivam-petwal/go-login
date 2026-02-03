package main

import (
	"go-login/config"
	"go-login/controller"
	"go-login/models"
	"go-login/repository"
	"go-login/router"
	"go-login/service"
	"log"
	"os"
)

func main() {
	config.ConnectDB()
	config.DBcon.AutoMigrate(&models.User{})

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production"
	}

	userRepo := repository.NewUserRepository(config.DBcon)
	authService := service.NewAuthService(userRepo, jwtSecret)
	authController := controller.NewAuthController(authService)

	r := router.SetupRouter(authController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
