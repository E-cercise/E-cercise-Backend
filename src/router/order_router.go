package router

import (
	"github.com/E-cercise/E-cercise/src/controller"
	"github.com/E-cercise/E-cercise/src/enum"
	"github.com/E-cercise/E-cercise/src/middleware"
	"github.com/E-cercise/E-cercise/src/middleware/validation"
	"github.com/E-cercise/E-cercise/src/repository"
	"github.com/gofiber/fiber/v2"
)

func OrderRouter(router fiber.Router, cartController *controller.CartController, userRepo repository.UserRepository) {
	orderGroup := router.Group("/order")

	orderGroup.Post("", middleware.Authentication(userRepo), middleware.RoleAuthorization(enum.RoleUser, enum.RoleAdmin), validation.ValidateAddLineEquipment(), cartController.AddEquipmentToCart)

}
