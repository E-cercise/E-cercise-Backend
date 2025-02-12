package controller

import (
	"fmt"
	"github.com/E-cercise/E-cercise/src/data/request"
	"github.com/E-cercise/E-cercise/src/helper"
	"github.com/E-cercise/E-cercise/src/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

func (c *CartController) DeleteItemInCart(ctx *fiber.Ctx) error {
	//TODO more: validate that another user cant delete item of another user cart
	lineEquipmentID := uuid.MustParse(ctx.Params("line_equipment_id"))

	status, err := c.CartService.DeleteLineEquipmentInCart(lineEquipmentID)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	if status == "204" {
		return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": "this line equipment not found or have deleted"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": fmt.Sprintf("line equipment id %v has been deleted successfully", lineEquipmentID)})

}

func (c *CartController) GetCartItems(ctx *fiber.Ctx) error {
	user, err := helper.GetCurrentUser(ctx)

	if err != nil {
		return err
	}

	resp, err := c.CartService.GetAllLineEquipmentInCart(user.ID)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("can't get equipment in cart: %v", err.Error())})
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func (c *CartController) ModifyItemInCart(ctx *fiber.Ctx) error {
	req, ok := ctx.Locals("req").(request.CartItemPutRequest)
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to parse request body (Controller)",
		})
	}

	if err := c.CartService.ModifyLineEquipmentInCart(req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("error cant modify item in cart with error: %v", err.Error()),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "cart modified successfully"})
}
