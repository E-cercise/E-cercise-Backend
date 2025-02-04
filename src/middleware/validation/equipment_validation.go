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
				"error":   "Invalid request format",
				"message": err.Error(),
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

func ValidateUpdateEquipment() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var req request.EquipmentPutRequest

		// Parse the request body into the struct
		if err := ctx.BodyParser(&req); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		if req.MuscleGroupUsed != nil && !request.ValidateMuscleGroup(req.MuscleGroupUsed) {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid muscle group format. Allowed formats are 'ft_{int}' and 'bk_{int}'",
			})
		}

		for _, opt := range req.Option.Updated {
			if opt.Images != nil {
				if err := request.ValidateImagePutReq(*opt.Images); err != nil {
					return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"error": "Invalid image, " + err.Error(),
					})
				}
			}
		}

		ctx.Locals("req", req)
		return ctx.Next()
	}

}
