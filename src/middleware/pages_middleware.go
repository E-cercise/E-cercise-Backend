package middleware

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func PreparePagination(defaultPage, defaultLimit string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page, err := strconv.Atoi(c.Query("page", defaultPage))
		if err != nil || page < 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "page must be a positive integer",
				"error":   "invalid page parameter",
			})
		}

		limit, err := strconv.Atoi(c.Query("limit", defaultLimit))
		if err != nil || limit < 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "limit must be a positive integer",
				"error":   "invalid limit parameter",
			})
		}

		// Store page and limit in the request context
		c.Locals("page", page)
		c.Locals("limit", limit)

		return c.Next()
	}
}
