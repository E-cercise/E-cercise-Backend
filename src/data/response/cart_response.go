package response

type GetCartItemResponse struct {
	LineEquipments []LineEquipment `json:"line_equipments"`
	TotalPrice     float64         `json:"total_price"`
}

type LineEquipment struct {
	EquipmentName   string  `json:"equipment_name"`
	LineEquipmentID string  `json:"line_equipment_id"`
	Quantity        int     `json:"quantity"`
	Total           float64 `json:"total"`
}
