package router

import (
	"go-login/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter(authController *controller.AuthController) *gin.Engine {
	router := gin.Default()

	router.POST("/login", authController.Login)

	return router
}
