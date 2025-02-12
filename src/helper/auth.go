package helper

import (
	"fmt"
	"github.com/E-cercise/E-cercise/src/config"
	"github.com/E-cercise/E-cercise/src/enum"
	"github.com/E-cercise/E-cercise/src/logger"
	"github.com/E-cercise/E-cercise/src/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

var secretKey = []byte(config.JwtSecret)

func CreateToken(userId uuid.UUID, name string, lastName string) (string, error) {

	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		logger.Log.WithError(err).Error("Failed to load timezone")
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"name":    name + " " + lastName,
		"exp":     time.Now().Add(time.Hour * 3).In(location).Unix(), // TODO: Need to change
	})
	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetClaimFromToken(jwtToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid or expired token: " + err.Error())
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid Token")
	}

	return claims, nil
}

func ContainsRole(roles []enum.Role, role enum.Role) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

func GetCurrentUser(c *fiber.Ctx) (*model.User, error) {
	currentUser := c.Locals("currentUser")
	if currentUser == nil {
		slog.Error("Unauthorized User: User doesn't exist")
		return nil, c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	user, ok := currentUser.(*model.User)
	if !ok {
		return nil, c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	return user, nil
}
