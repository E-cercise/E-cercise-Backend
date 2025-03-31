package controller

import (
	"github.com/E-cercise/E-cercise/src/data/request"
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
		return err
	}

	userDetail := c.UserService.GetUserProfile(user)

	return ctx.Status(fiber.StatusOK).JSON(userDetail)

}

func (c *UserController) UpdateProfile(ctx *fiber.Ctx) error {
	reqBody := ctx.Locals("reqBody").(request.UpdateUserProfileRequest)

	user, err := helper.GetCurrentUser(ctx)
	if err != nil {
		return err
	}

	err = c.UserService.UpdateUserProfile(user, reqBody)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to update user profile")
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "Failed to update user profile",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User profile updated",
	})
}

func (c *UserController) CheckEmailExists(ctx *fiber.Ctx) error {
	var emailReq request.CheckEmailExistsRequest

	if err := ctx.QueryParser(&emailReq); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "Invalid email address",
		})
	}
	isUserExists, err := c.UserService.CheckUserExist(emailReq.Email)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "Failed to check user existence",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"isUserExists": isUserExists,
	})

}
