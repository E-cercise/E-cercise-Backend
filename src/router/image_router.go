package router

import (
	"github.com/E-cercise/E-cercise/src/controller"
	"github.com/E-cercise/E-cercise/src/enum"
	"github.com/E-cercise/E-cercise/src/middleware"
	"github.com/E-cercise/E-cercise/src/repository"
	"github.com/gofiber/fiber/v2"
)

func ImageRouter(router fiber.Router, imageController *controller.ImageController, userRepo repository.UserRepository) {
	imageGroup := router.Group("/image")

	imageGroup.Post("/upload", middleware.Authentication(userRepo), middleware.RoleAuthorization(enum.RoleAdmin), imageController.UplaodFile)
}
