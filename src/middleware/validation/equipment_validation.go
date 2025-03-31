package validation

import (
	"github.com/E-cercise/E-cercise/src/data/request"
	"github.com/E-cercise/E-cercise/src/logger"
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

		// Log raw body (optional, for debugging)
		raw := ctx.Body()
		logger.Log.Infof("Received UpdateEquipment payload: %s", string(raw))

		// Parse the request body into the struct
		if err := ctx.BodyParser(&req); err != nil {
			logger.Log.WithError(err).Error("Failed to parse equipment update request")
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		// Validate muscle group format
		if req.MuscleGroupUsed != nil && !request.ValidateMuscleGroup(req.MuscleGroupUsed) {
			logger.Log.Warn("Invalid muscle group format in update request")
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid muscle group format. Allowed formats are 'ft_{int}' and 'bk_{int}'",
			})
		}

		// Safe nil-checking for option updates
		if req.Option != nil && req.Option.Updated != nil {
			for _, opt := range req.Option.Updated {
				if opt.Images != nil {
					if err := request.ValidateImagePutReq(*opt.Images); err != nil {
						logger.Log.WithError(err).Error("Invalid image format in update request")
						return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
							"error": "Invalid image, " + err.Error(),
						})
					}
				}
			}
		}

		// Store validated request for controller
		ctx.Locals("req", req)
		return ctx.Next()
	}
}
