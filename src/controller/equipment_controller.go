package controller

import (
	"github.com/E-cercise/E-cercise/src/data/request"
	"github.com/E-cercise/E-cercise/src/helper"
	"github.com/E-cercise/E-cercise/src/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

func (c *EquipmentController) AddEquipment(ctx *fiber.Ctx) error {
	req, ok := ctx.Locals("req").(request.EquipmentPostRequest)

	if !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to parse request body (Controller)",
		})
	}

	if err := c.EquipmentService.AddEquipment(req, ctx.Context()); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":     "failed to add equipment",
			"error_msg": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "equipment add successfully",
	})

}

func (c *EquipmentController) GetEquipment(ctx *fiber.Ctx) error {
	equipmentID := uuid.MustParse(ctx.Params("id"))

	resp, err := c.EquipmentService.GetEquipmentDetail(equipmentID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func (c *EquipmentController) UpdateEquipment(ctx *fiber.Ctx) error {
	equipmentID := uuid.MustParse(ctx.Params("id"))
	req := ctx.Locals("req").(request.EquipmentPutRequest)

	err := c.EquipmentService.UpdateEquipment(equipmentID, ctx.Context(), req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "equipment update successfully"})
}
