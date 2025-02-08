package router

import (
	"github.com/E-cercise/E-cercise/src/controller"
	"github.com/E-cercise/E-cercise/src/enum"
	"github.com/E-cercise/E-cercise/src/middleware"
	"github.com/E-cercise/E-cercise/src/repository"
	"github.com/gofiber/fiber/v2"
)

func CartRouter(router fiber.Router, cartController *controller.CartController, userRepo repository.UserRepository) {
	cartGroup := router.Group("/cart")

	cartGroup.Post("/cart/item", middleware.Authentication(userRepo), middleware.RoleAuthorization(enum.RoleUser, enum.RoleAdmin), cartController.AddEquipmentToCart)

	//imageGroup.Post("/upload", middleware.Authentication(userRepo), middleware.RoleAuthorization(enum.RoleAdmin))

}
