package controller

import (
	"fmt"
	"github.com/E-cercise/E-cercise/src/data/request"
	"github.com/E-cercise/E-cercise/src/helper"
	"github.com/E-cercise/E-cercise/src/service"
	"github.com/gofiber/fiber/v2"
)

type CartController struct {
	CartService service.CartService
}

func NewCartControllerImpl(cartService service.CartService) *CartController {
	return &CartController{
		CartService: cartService,
	}
}

func (c *CartController) AddEquipmentToCart(ctx *fiber.Ctx) error {
	req, ok := ctx.Locals("req").(request.CartItemPostRequest)

	if !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to parse request body (Controller)",
		})
	}

	user, err := helper.GetCurrentUser(ctx)

	if err != nil {
		return err
	}

	err = c.CartService.AddEquipmentToCart(req, user.ID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": fmt.Sprintf("equipmentID: %v option: %v add to cart successfully", req.EquipmentID, req.EquipmentOptionID),
	})
}
