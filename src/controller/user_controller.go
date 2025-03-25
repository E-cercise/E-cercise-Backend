package controller

import (
	"github.com/E-cercise/E-cercise/src/helper"
	"github.com/E-cercise/E-cercise/src/logger"
	"github.com/E-cercise/E-cercise/src/service"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	UserService service.UserService
}

func NewUserControllerImpl(userService service.UserService) *UserController {
	return &UserController{UserService: userService}
}

func (c *UserController) GetProfile(ctx *fiber.Ctx) error {
	user, err := helper.GetCurrentUser(ctx)

	if err != nil {
		logger.Log.WithError(err).Error("Failed to get current user")
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get current user")
	}

	userDetail := c.UserService.GetUserProfile(user)

	return ctx.Status(fiber.StatusOK).JSON(userDetail)

}
