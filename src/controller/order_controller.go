package controller

import (
	"github.com/E-cercise/E-cercise/src/data/request"
	"github.com/E-cercise/E-cercise/src/helper"
	"github.com/E-cercise/E-cercise/src/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OrderController struct {
	OrderService service.OrderService
}

func NewOrderControllerImpl(orderService service.OrderService) *OrderController {
	return &OrderController{
		OrderService: orderService,
	}
}

func (c *OrderController) CreateOrder(ctx *fiber.Ctx) error {
	req, ok := ctx.Locals("req").(request.CheckoutCartRequest)

	if !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to parse request body (Controller)",
		})
	}

	user, err := helper.GetCurrentUser(ctx)

	if err != nil {
		return err
	}

	if err = c.OrderService.CreateOrder(req, user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Order created successfully",
	})
}

func (c *OrderController) GetOrderDetail(ctx *fiber.Ctx) error {
	orderID := uuid.MustParse(ctx.Params("id"))

	user, err := helper.GetCurrentUser(ctx)

	if err != nil {
		return err
	}

	resp, err := c.OrderService.GetOrderDetail(orderID, user)

	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func (c *OrderController) UpdateOrderStatus(ctx *fiber.Ctx) error {
	orderID := uuid.MustParse(ctx.Params("id"))

	err := c.OrderService.UpdateOrderStatus(orderID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Order status updated successfully",
	})
}

func (c *OrderController) GetMyOrders(ctx *fiber.Ctx) error {

	user, err := helper.GetCurrentUser(ctx)
	if err != nil {
		return err
	}

	var req request.OrderDetailRequest

	if err := ctx.QueryParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request data",
		})
	}

	resp, err := c.OrderService.GetMyOrders(user.ID, req.OrderStatus)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error(), "message": "error during get my orders"})
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}
