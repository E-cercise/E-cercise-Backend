package response

import (
	"github.com/E-cercise/E-cercise/src/enum"
	"github.com/google/uuid"
)

type Address struct {
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	AddressLine string `json:"address_line"`
}

type OrderDetailResponse struct {
	ID          uuid.UUID        `json:"id"`
	OrderStatus enum.OrderStatus `json:"order_status"`
	Address     Address          `json:"address"`
	Orders      []LineEquipment  `json:"orders"`
	NetPrice    float64          `json:"net_price"`
}

type OrderListResponse struct {
	Orders []Order `json:"orders"`
}

type Order struct {
	CreatedAt          string             `json:"created_at"`
	FirstLineEquipment FirstLineEquipment `json:"first_line_equipment"`
	ID                 uuid.UUID          `json:"id"`
	OrderStatus        enum.OrderStatus   `json:"order_status"`
	PaymentType        enum.PaymentType   `json:"payment_type"`
	TotalPrice         float64            `json:"total_price"`
	UpdatedAt          string             `json:"updated_at"`
}

type FirstLineEquipment struct {
	ImgURL string `json:"img_url"`
	Name   string `json:"name"`
}
