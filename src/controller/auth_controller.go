package controller

import (
	"fmt"
	"github.com/E-cercise/E-cercise/src/data/request"
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

type AuthController struct {
	UserService service.UserService
}

func NewAuthControllerImpl(userService service.UserService) *AuthController {
	return &AuthController{
		UserService: userService,
	}
}

func (c *AuthController) UserRegister(ctx *fiber.Ctx) error {
	reqBody := ctx.Locals("reqBody").(request.RegisterRequest)

	user, err := c.UserService.RegisterUser(reqBody)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    user,
	})
}
