package model

import (
	"github.com/E-cercise/E-cercise/src/enum"
	"github.com/google/uuid"
	"time"
)

type Order struct {
	ID              uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID          uuid.UUID        `gorm:"type:uuid;not null" json:"user_id"`
	User            User             `gorm:"foreignKey:UserID" json:"user"`
	LineEquipments  []LineEquipment  `gorm:"foreignKey:OrderID" json:"line_equipments"`
	DeliveryAddress string           `gorm:"type:text;not null" json:"delivery_address"`
	PaymentType     enum.PaymentType `gorm:"type:payment_type" json:"payment_type"`
	TotalPrice      float64          `gorm:"type:decimal(10,2);not null" json:"total_price"`
	OrderStatus     enum.OrderStatus `gorm:"type:order_status;default:'Placed';not null" json:"order_status"`
	CreatedAt       time.Time        `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time        `gorm:"default:CURRENT_TIMESTAMP"`
}

func (o *Order) BeforeUpdate() error {
	o.UpdatedAt = time.Now()
	return nil
}
