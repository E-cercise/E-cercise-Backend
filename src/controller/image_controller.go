package controller

import (
	"github.com/E-cercise/E-cercise/src/logger"
	"github.com/E-cercise/E-cercise/src/service"
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

type ImageController struct {
	ImageService service.ImageService
}

func NewImageControllerImpl(imageService service.ImageService) *ImageController {
	return &ImageController{
		ImageService: imageService,
	}
}

func (c *ImageController) UplaodFile(ctx *fiber.Ctx) error {

	// Parse the uploaded file
	fileHeader, err := ctx.FormFile("img")
	if err != nil {
		slog.Error("Error during file upload", "err", err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "file not provided or invalid"})
	}

	// Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		logger.Log.WithError(err).Error("Error opening uploaded file", "err", err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to open the file"})
	}
	defer file.Close()

	fileID, err := c.ImageService.UploadImage(ctx.Context(), file, fileHeader, false)

	if err != nil {
		logger.Log.WithError(err).Error("error uploading image", "err", err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error uploading image: " + err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "file uploaded successfully", "fileID": fileID})
}
