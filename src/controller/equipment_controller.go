package controller

import (
	"github.com/E-cercise/E-cercise/src/data/request"
	"github.com/E-cercise/E-cercise/src/helper"
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
	page, ok := ctx.Locals("page").(int)

	if !ok {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "page not found in context",
		})
	}
	limit, ok := ctx.Locals("limit").(int)
	if !ok {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "limit not found in context",
		})
	}

	paginator := helper.NewPaginator(page, limit)

	var req request.EquipmentListRequest

	if err := ctx.QueryParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request data",
		})
	}

	equipments, err := c.EquipmentService.GetEquipmentData(req, paginator)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"equipments":  &equipments,
		"page":        paginator.Page,
		"limit":       paginator.Limit,
		"total_pages": paginator.TotalPages,
		"total_rows":  paginator.TotalRows,
	})
}
