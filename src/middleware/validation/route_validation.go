package validation

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"regexp"
	"strconv"
)

func ValidateParam(paramName, validationType string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		value := c.Params(paramName)
		if value == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": paramName + " is required"})
		}

		switch validationType {
		case "uuid":
			if _, err := uuid.Parse(value); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": paramName + " must be a valid UUID"})
			}
		case "int":
			if _, err := strconv.Atoi(value); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": paramName + " must be a valid integer"})
			}
		case "default":
			// allow only alphanumeric, underscores, and hyphens (a-z, A-Z, 0-9, _, -)
			regex := regexp.MustCompile("^[a-zA-Z0-9_-]+$")
			if !regex.MatchString(value) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": paramName + " must be alphanumeric"})
			}
		default:
		}

		return c.Next()
	}
}
