package router

import (
	"github.com/E-cercise/E-cercise/src/controller"
	"github.com/gofiber/fiber/v2"
)

func TagRouter(app fiber.Router, ctrl controller.TagController) {
	tag := app.Group("/tags")
	tag.Get("/", ctrl.GetTags)
	tag.Get("/me", ctrl.GetUserPreferences)
	tag.Post("/me", ctrl.SetUserPreferences)
}
