package router

import (
	"github.com/E-cercise/E-cercise/src/controller"
	"github.com/E-cercise/E-cercise/src/middleware/validation"
	"github.com/gofiber/fiber/v2"
)

func AuthRouter(router fiber.Router, authController *controller.AuthController) {
	authGroup := router.Group("/auth")

	authGroup.Post("/register", validation.ValidateUserRegister(), authController.UserRegister)
	authGroup.Post("/login", validation.ValidateLoginRequest(), authController.Login)

}
