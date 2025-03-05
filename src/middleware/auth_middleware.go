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

func OptionalAuthentication(userRepo repository.UserRepository) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 1) Check the Authorization header
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			// No token is provided, skip the user check
			// Let the request continue without a user
			return ctx.Next()
		}

		// 2) Split out "Bearer" and the token
		authToken := strings.Split(authHeader, " ")
		if len(authToken) != 2 || authToken[0] != "Bearer" {
			// If token is malformed, still skip *failing*
			// so we do not block unauthenticated usage
			return ctx.Next()
		}

		tokenString := authToken[1]

		// 3) Attempt to parse the token
		claims, err := helper.GetClaimFromToken(tokenString)
		if err != nil {
			// Token is invalid, but we do not throw an error
			return ctx.Next()
		}

		// 4) Extract user info from the token
		userID, ok := claims["user_id"].(string)
		if !ok {
			// If userID is missing or incorrectly typed,
			// we skip setting user and proceed
			return ctx.Next()
		}

		// 5) Fetch the user from the DB
		user, err := userRepo.FindByID(userID)
		if err != nil || user == nil {
			// If user does not exist, skip
			return ctx.Next()
		}

		// 6) If everything is valid, store user in context
		ctx.Locals("currentUser", user)
		return ctx.Next()
	}
}
