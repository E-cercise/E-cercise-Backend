package model

import "github.com/google/uuid"

type LineEquipment struct {
	ID                uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	OrderID           *uuid.UUID `gorm:"type:uuid;"`
	CartID            uuid.UUID  `gorm:"type:uuid;not null" json:"cart_id"`
	EquipmentID       uuid.UUID  `gorm:"type:uuid;not null" json:"equipment_id"`
	EquipmentOptionID uuid.UUID  `gorm:"type:uuid;not null" json:"equipment_option_id"`
	Quantity          int        `gorm:"type:int;not null;default:1" json:"quantity"`
}
