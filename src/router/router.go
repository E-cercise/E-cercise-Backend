package router

import (
	"github.com/E-cercise/E-cercise/src/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
	"net/http"
)

func InitRouter(db *gorm.DB) *fiber.App {

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173, https://login.microsoftonline.com, " + config.FrontendBaseURL,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, ngrok-skip-browser-warning, Authorization, Access-Control-Allow-Origin",
		AllowCredentials: true,
		ExposeHeaders:    "content-disposition",
	}))

	app.Use(helmet.New())

	// Recovery middleware
	app.Use(recover.New())

	// Define API group
	apiGroup := app.Group("/api")

	// Root endpoint
	apiGroup.Get("", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Hello TimeSheet"})
	})

	return app
}
