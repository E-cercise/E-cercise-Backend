package validation

import (
	"github.com/E-cercise/E-cercise/src/data/request"
	"github.com/gofiber/fiber/v2"
)

func ValidateAddEquipment() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var req request.EquipmentPostRequest

		// Parse the request body into the struct
		if err := ctx.BodyParser(&req); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		if !request.ValidateMuscleGroup(req.MuscleGroupUsed) {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid muscle group format. Allowed formats are 'ft_{int}' and 'bk_{int}'",
			})
		}

		ctx.Locals("req", req)

		return ctx.Next()

	}
}
