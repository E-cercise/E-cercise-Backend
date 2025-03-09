package request

import "github.com/google/uuid"

type CheckoutCartRequest struct {
	LineEquipments []uuid.UUID `json:"line_equipments"`
}
