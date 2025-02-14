package validation

import (
	"encoding/json"
	"fmt"
	"github.com/E-cercise/E-cercise/src/config"
	"github.com/E-cercise/E-cercise/src/data/request"
	"github.com/E-cercise/E-cercise/src/helper"
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

		type payloadStruct struct {
			IV        string `json:"iv"`
			LoginBody string `json:"login_body"`
		}

		var loginPayload payloadStruct
		if err := c.BodyParser(&loginPayload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request data",
			})
		}

		decryptedJSON, err := helper.DecryptPayload(loginPayload.LoginBody, loginPayload.IV, config.SecretKey)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Decryption failed",
				"message": err.Error(),
			})
		}

		var loginBody request.LoginRequest
		if err := json.Unmarshal([]byte(decryptedJSON), &loginBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"err":     "Invalid decrypted JSON",
				"message": err.Error(),
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
