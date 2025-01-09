package controller

import (
	"github.com/E-cercise/E-cercise/src/service"
	"github.com/gofiber/fiber/v2"
)

type EquipmentController struct {
	EquipmentService service.EquipmentService
}

func NewEquipmentControllerImpl(equipmentService service.EquipmentService) *EquipmentController {
	return &EquipmentController{
		EquipmentService: equipmentService,
	}
}

func (c *EquipmentController) GetAllEquipment(ctx *fiber.Ctx) error {

	equipments, err := c.EquipmentService.GetAllEquipmentData()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"equipments": &equipments,
	})
}
