package controller

import (
	"github.com/E-cercise/E-cercise/src/data/request"
	"github.com/E-cercise/E-cercise/src/service"
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

func (c *OrderController) CreateOrder(req request.CartItemPostRequest, userID uuid.UUID) error {

}
