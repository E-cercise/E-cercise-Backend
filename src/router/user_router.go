package router

import (
	"github.com/E-cercise/E-cercise/src/controller"
	"github.com/E-cercise/E-cercise/src/middleware"
	"github.com/E-cercise/E-cercise/src/repository"
	"github.com/gofiber/fiber/v2"
)

func UserRouter(router fiber.Router, userController *controller.UserController, userRepo repository.UserRepository) {
	profileGroup := router.Group("/profile")
	profileGroup.Get("/me", middleware.Authentication(userRepo), userController.GetProfile)
}
