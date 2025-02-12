package router

import (
	"github.com/E-cercise/E-cercise/src/config"
	"github.com/E-cercise/E-cercise/src/controller"
	logger2 "github.com/E-cercise/E-cercise/src/logger"
	"github.com/E-cercise/E-cercise/src/repository"
	"github.com/E-cercise/E-cercise/src/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
	"net/http"
)

func InitRouter(db *gorm.DB) *fiber.App {

	userRepo := repository.NewUserRepository(db)
	equipmentRepo := repository.NewEquipmentRepository(db)
	imageRepo := repository.NewImageRepository(db)
	muscleGroupRepo := repository.NewMuscleGroupRepository(db)
	cartRepo := repository.NewCartRepository(db)

	cloudinaryService, err := service.NewCloudinaryService()

	if err != nil {
		panic(err)
	}

	userService := service.NewUserService(db, userRepo)
	imageService := service.NewImageService(db, imageRepo, cloudinaryService)
	equipmentService := service.NewEquipmentService(db, equipmentRepo, muscleGroupRepo, imageService)
	cartService := service.NewCartService(db, cartRepo, equipmentRepo)

	authController := controller.NewAuthControllerImpl(userService)
	equipmentController := controller.NewEquipmentControllerImpl(equipmentService)
	imageController := controller.NewImageControllerImpl(imageService)
	cartController := controller.NewCartControllerImpl(cartService)

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

	app.Use(logger.New())

	// Define API group
	apiGroup := app.Group("/api")

	// Root endpoint
	apiGroup.Get("", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Hello E-cercise"})
	})

	AuthRouter(apiGroup, authController)
	EquipmentRouter(apiGroup, equipmentController, userRepo)
	ImageRouter(apiGroup, imageController, userRepo)
	CartRouter(apiGroup, cartController, userRepo)

	logger2.Log.Info("Router initialized")
	for _, route := range app.GetRoutes() {
		// You can format the output however you like
		logger2.Log.Infof("Method: %s \t Path: %s\n", route.Method, route.Path)
	}

	return app
}
