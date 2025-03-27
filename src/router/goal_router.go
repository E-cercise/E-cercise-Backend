package router

import (
	"github.com/E-cercise/E-cercise/src/controller"
	"github.com/gofiber/fiber/v2"
)

func GoalRouter(app fiber.Router, goalController *controller.GoalController) {
	goalGroup := app.Group("/goals")
	goalGroup.Get("/", goalController.GetGoals)
}
