package validation

import (
	"fmt"
	"github.com/E-cercise/E-cercise/src/data/request"
	"github.com/gofiber/fiber/v2"
)

func ValidateAddLineEquipment() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var req request.CartItemPostRequest

		// Parse the request body into the struct
		if err := ctx.BodyParser(&req); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid request format",
				"message": err.Error(),
			})
		}

		ctx.Locals("req", req)
		return ctx.Next()
	}
}

func ValidateModifyLineEquipmentRequest() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var req request.CartItemPutRequest

		// Parse the request body into the struct
		if err := ctx.BodyParser(&req); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid request format",
				"message": err.Error(),
			})
		}

		for _, item := range req.Items {
			if item.Quantity <= 0 {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": fmt.Sprintf("Quantity of lineEquipmentID: %v must more than 0", item.LineEquipmentID),
				})
			}
		}

		ctx.Locals("req", req)
		return ctx.Next()
	}
}
