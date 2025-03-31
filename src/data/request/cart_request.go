package request

type CartItemGetRequest struct {
	LineEquipmentIDs string `query:"line_equipment_ids"`
}

type CartItemPostRequest struct {
	EquipmentID       string `json:"equipment_id"`
	EquipmentOptionID string `json:"equipment_option_id"`
	Quantity          int    `json:"quantity"`
}

type CartItemPutRequest struct {
	Items []Item `json:"items"`
}

type Item struct {
	LineEquipmentID string `json:"line_equipment_id"`
	Quantity        int    `json:"quantity"`
}
