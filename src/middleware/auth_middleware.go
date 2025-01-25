package middleware

import (
	"github.com/E-cercise/E-cercise/src/enum"
	"github.com/E-cercise/E-cercise/src/helper"
	"github.com/E-cercise/E-cercise/src/model"
	"github.com/E-cercise/E-cercise/src/repository"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func Authentication(userRepo repository.UserRepository) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")

		if authHeader == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header is missing"})
		}

		authToken := strings.Split(authHeader, " ")

		if len(authToken) != 2 || authToken[0] != "Bearer" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
		}

		tokenString := authToken[1]
		claims, err := helper.GetClaimFromToken(tokenString)

		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token payload"})
		}

		user, err := userRepo.FindByID(userID)
		if err != nil || user == nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
		}

		ctx.Locals("currentUser", user)

		return ctx.Next()
	}
}

func RoleAuthorization(allowedRoles ...enum.Role) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Retrieve currentUser from the context
		currentUser := ctx.Locals("currentUser")
		if currentUser == nil {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Roles not found in context"})
		}

		// Type assert currentUser to *model.User
		user, ok := currentUser.(*model.User)
		if !ok {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Invalid user context"})
		}

		// Check if any user role matches allowed roles
		if helper.ContainsRole(allowedRoles, user.Role) {
			return ctx.Next()
		}

		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}
}
