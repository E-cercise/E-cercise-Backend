package controller

import (
	"github.com/E-cercise/E-cercise/src/service"
	"github.com/gofiber/fiber/v2"
)

type GoalController struct {
	goalService service.GoalService
}

func NewGoalControllerImpl(s service.GoalService) *GoalController {
	return &GoalController{goalService: s}
}

func (c *GoalController) GetGoals(ctx *fiber.Ctx) error {
	tags, err := c.goalService.GetAllGoal()
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": "cannot fetch goals"})
	}
	return ctx.JSON(tags)
}
