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
	OrderStatus *enum.OrderStatus `query:"order_status"`
	UserID      *uuid.UUID        `query:"user_id"`
	OrderID     *uuid.UUID        `query:"order_id"`
	PaymentType *enum.PaymentType `query:"payment_type"`
}
