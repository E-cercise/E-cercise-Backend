package model

import (
	"github.com/E-cercise/E-cercise/src/enum"
	"github.com/google/uuid"
)

type Order struct {
	ID              uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID          uuid.UUID        `gorm:"type:uuid;not null" json:"user_id"`
	User            User             `gorm:"foreignKey:UserID" json:"user"`
	LineEquipments  []LineEquipment  `gorm:"foreignKey:OrderID" json:"line_equipments"`
	DeliveryAddress string           `gorm:"type:text;not null" json:"delivery_address"`
	PaymentType     string           `gorm:"type:varchar(50);not null" json:"payment_type"`
	TotalPrice      float64          `gorm:"type:decimal(10,2);not null" json:"total_price"`
	OrderStatus     enum.OrderStatus `gorm:"type:order_status;default:'Placed';not null" json:"order_status"`
}
