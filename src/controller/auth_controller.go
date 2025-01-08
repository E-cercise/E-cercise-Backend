package controller

import (
	"github.com/E-cercise/E-cercise/src/data/request"
	"github.com/E-cercise/E-cercise/src/service"
	"github.com/gofiber/fiber/v2"
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

	err := c.UserService.RegisterUser(reqBody)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User registered successfully",
	})
}
