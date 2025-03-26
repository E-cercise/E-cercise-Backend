package router

import (
	"github.com/E-cercise/E-cercise/src/controller"
	"github.com/E-cercise/E-cercise/src/enum"
	"github.com/E-cercise/E-cercise/src/middleware"
	"github.com/E-cercise/E-cercise/src/middleware/validation"
	"github.com/E-cercise/E-cercise/src/repository"
	"github.com/gofiber/fiber/v2"
)

func OrderRouter(router fiber.Router, orderController *controller.OrderController, userRepo repository.UserRepository) {
	orderGroup := router.Group("/order")

	orderGroup.Post("", middleware.Authentication(userRepo), middleware.RoleAuthorization(enum.RoleUser, enum.RoleAdmin), validation.ValidateCheckoutOrder(), orderController.CreateOrder)
	orderGroup.Get("/me", middleware.Authentication(userRepo), middleware.RoleAuthorization(enum.RoleUser), orderController.GetMyOrders)
	orderGroup.Get("/list", middleware.Authentication(userRepo), middleware.RoleAuthorization(enum.RoleAdmin), orderController.GetOrderList)
	orderGroup.Get("/:id", middleware.Authentication(userRepo), middleware.RoleAuthorization(enum.RoleUser, enum.RoleAdmin), validation.ValidateParam("id", "uuid"), orderController.GetOrderDetail)
	orderGroup.Put("/status/:id", middleware.Authentication(userRepo), middleware.RoleAuthorization(enum.RoleUser, enum.RoleAdmin), validation.ValidateParam("id", "uuid"), orderController.UpdateOrderStatus)
}
