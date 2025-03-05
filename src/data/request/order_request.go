package request

type CheckoutCartRequest struct {
	Items []CartItem `json:"items"`
}

type CartItem struct {
	LineEquipmentID string `json:"line_equipment_id"`
	Quantity        int    `json:"quantity"`
}
