package controller

import (
	"github.com/E-cercise/E-cercise/src/helper"
	"github.com/E-cercise/E-cercise/src/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TagController struct {
	tagService  service.TagService
	prefService service.UserPreferenceService
}

func NewTagControllerImpl(t service.TagService, p service.UserPreferenceService) TagController {
	return TagController{t, p}
}

func (ctrl TagController) GetTags(ctx *fiber.Ctx) error {
	tags, err := ctrl.tagService.GetAllTags()
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": "cannot fetch tags"})
	}
	return ctx.JSON(tags)
}

func (ctrl TagController) SetUserPreferences(ctx *fiber.Ctx) error {
	user, err := helper.GetCurrentUser(ctx)

	if err != nil {
		return err
	}

	var req struct {
		TagIDs []uuid.UUID `json:"tag_ids"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "invalid input"})
	}

	if err := ctrl.prefService.SetUserPreferences(user.ID, req.TagIDs); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot set preferences"})
	}
	return ctx.JSON(fiber.Map{"message": "preferences saved"})
}

func (ctrl TagController) GetUserPreferences(ctx *fiber.Ctx) error {
	user, err := helper.GetCurrentUser(ctx)
	if err != nil {
		return err
	}

	tags, err := ctrl.prefService.GetUserPreferences(user.ID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot get preferences"})
	}
	return ctx.JSON(tags)

}
