package response

import (
	"github.com/google/uuid"
)

type Address struct {
	FullName string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	AddressLine string `json:"address_line"`
}

type OrderResponse struct {
	ID uuid.UUID `json:"id"`
	OrderStatus string `json:"order_status"`
	Address Address `json:"address"`
	Orders []LineEquipment `json:"orders"`
	NetPrice float64 `json:"net_price"`
}