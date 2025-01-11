package router

import (
	"github.com/E-cercise/E-cercise/src/controller"
	"github.com/E-cercise/E-cercise/src/enum"
	"github.com/E-cercise/E-cercise/src/middleware"
	"github.com/E-cercise/E-cercise/src/repository"
	"github.com/gofiber/fiber/v2"
)

func EquipmentRouter(router fiber.Router, equipmentController *controller.EquipmentController, userRepo repository.UserRepository) {
	equipmentGroup := router.Group("/equipment")

	equipmentGroup.Get("/all", middleware.Authentication(userRepo),
		middleware.RoleAuthorization(enum.RoleUser, enum.RoleAdmin), equipmentController.GetAllEquipment) //group by collaborative filtering later

}
