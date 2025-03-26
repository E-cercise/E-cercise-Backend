package request

import (
	"github.com/E-cercise/E-cercise/src/enum"
	"github.com/google/uuid"
)

type CheckoutCartRequest struct {
	LineEquipments []uuid.UUID      `json:"line_equipments"`
	PaymentType    enum.PaymentType `json:"payment_type"`
	Address        string           `json:"address"`
}
