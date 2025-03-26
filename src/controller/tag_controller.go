package controller

import (
	"github.com/E-cercise/E-cercise/src/service"
	"github.com/gofiber/fiber/v2"
)

type TagController struct {
	tagService  service.TagService
	prefService service.UserPreferenceService
}

func NewTagController(t service.TagService, p service.UserPreferenceService) TagController {
	return TagController{t, p}
}

func (ctrl TagController) GetTags(c *fiber.Ctx) error {
	tags, err := ctrl.tagService.GetAllTags()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "cannot fetch tags"})
	}
	return c.JSON(tags)
}

func (ctrl TagController) SetUserPreferences(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string) // assuming middleware sets this
	var req struct {
		TagIDs []string `json:"tag_ids"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}

	if err := ctrl.prefService.SetUserPreferences(userID, req.TagIDs); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "cannot set preferences"})
	}
	return c.JSON(fiber.Map{"message": "preferences saved"})
}

func (ctrl TagController) GetUserPreferences(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	tags, err := ctrl.prefService.GetUserPreferences(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "cannot get preferences"})
	}
	return c.JSON(tags)
}
