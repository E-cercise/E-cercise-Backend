package response

import (
	"github.com/E-cercise/E-cercise/src/enum"
	"github.com/google/uuid"
)

type Address struct {
	FullName string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	AddressLine string `json:"address_line"`
}

type OrderDetailResponse struct {
	ID uuid.UUID `json:"id"`
	OrderStatus enum.OrderStatus `json:"order_status"`
	Address Address `json:"address"`
	Orders []LineEquipment `json:"orders"`
	NetPrice float64 `json:"net_price"`
}
