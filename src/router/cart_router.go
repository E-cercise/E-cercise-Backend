package router

import (
	"github.com/E-cercise/E-cercise/src/controller"
	"github.com/E-cercise/E-cercise/src/enum"
	"github.com/E-cercise/E-cercise/src/middleware"
	"github.com/E-cercise/E-cercise/src/middleware/validation"
	"github.com/E-cercise/E-cercise/src/repository"
	"github.com/gofiber/fiber/v2"
)

func CartRouter(router fiber.Router, cartController *controller.CartController, userRepo repository.UserRepository) {
	cartGroup := router.Group("/cart")

	cartGroup.Post("/item", middleware.Authentication(userRepo), middleware.RoleAuthorization(enum.RoleUser, enum.RoleAdmin), validation.ValidateAddLineEquipment(), cartController.AddEquipmentToCart)
	cartGroup.Delete("/:line_equipment_id", middleware.Authentication(userRepo), middleware.RoleAuthorization(enum.RoleUser, enum.RoleAdmin), validation.ValidateParam("line_equipment_id", "uuid"),
		cartController.DeleteItemInCart)

	//imageGroup.Post("/upload", middleware.Authentication(userRepo), middleware.RoleAuthorization(enum.RoleAdmin))

}
