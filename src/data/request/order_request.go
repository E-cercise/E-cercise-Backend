package request

import (
	"github.com/E-cercise/E-cercise/src/enum"
	"github.com/google/uuid"
)

type CheckoutCartRequest struct {
	LineEquipments []uuid.UUID `json:"line_equipments"`
}

type OrderMeRequest struct {
	OrderStatus enum.OrderStatus `query:"order_status"`
}

type OrderListRequest struct {
	OrderStatus *string `query:"order_status"`
	UserID      *string `query:"user_id"`
	OrderID     *string `query:"order_id"`
	PaymentType *string `query:"payment_type"`
}
