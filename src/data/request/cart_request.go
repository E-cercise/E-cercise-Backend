package request

type CartItemPostRequest struct {
	EquipmentID       string `json:"equipment_id"`
	EquipmentOptionID string `json:"equipment_option_id"`
	Quantity          int    `json:"quantity"`
}
