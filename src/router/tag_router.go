package router

import (
	"github.com/E-cercise/E-cercise/src/controller"
	"github.com/E-cercise/E-cercise/src/middleware"
	"github.com/E-cercise/E-cercise/src/repository"
	"github.com/gofiber/fiber/v2"
)

func TagRouter(app fiber.Router, tagController controller.TagController, userRepo repository.UserRepository) {
	tag := app.Group("/tags")
	tag.Get("/", middleware.Authentication(userRepo), tagController.GetTags)
	tag.Get("/me", middleware.Authentication(userRepo), tagController.GetUserPreferences)
	tag.Post("/me", middleware.Authentication(userRepo), tagController.SetUserPreferences)
}
