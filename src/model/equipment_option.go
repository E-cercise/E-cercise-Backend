package model

import "github.com/google/uuid"

type EquipmentOption struct {
	ID                uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	EquipmentID       uuid.UUID `gorm:"type:uuid;not null" json:"equipment_id"`
	Weight            float64   `gorm:"type:decimal(10,2);not null" json:"weight"`
	Price             float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	RemainingProducts int       `gorm:"type:int;not null;default:0" json:"remaining_products"`
}
