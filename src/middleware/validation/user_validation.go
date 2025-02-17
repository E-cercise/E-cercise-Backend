package validation

import (
	"fmt"
	"github.com/E-cercise/E-cercise/src/data/request"
	"github.com/E-cercise/E-cercise/src/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func ValidateUserRegister() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var reqBody request.RegisterRequest
		if err := ctx.BodyParser(&reqBody); err != nil {
			logger.Log.WithError(err).Error("Invalid JSON data")
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid JSON data: %v", err.Error()),
			})
		}

		if err := validate.Struct(&reqBody); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		ctx.Locals("reqBody", reqBody)
		return ctx.Next()
	}
}

func ValidateLoginRequest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var loginBody request.LoginRequest
		if err := c.BodyParser(&loginBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request data",
			})
		}

		if err := validate.Struct(&loginBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		c.Locals("loginBody", loginBody)
		return c.Next()
	}
}
