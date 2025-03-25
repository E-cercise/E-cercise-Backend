package request

import (
	"github.com/E-cercise/E-cercise/src/enum"
	"github.com/google/uuid"
)

type CheckoutCartRequest struct {
	LineEquipments []uuid.UUID `json:"line_equipments"`
}

type OrderDetailRequest struct {
	OrderStatus enum.OrderStatus `query:"order_status"`
}
