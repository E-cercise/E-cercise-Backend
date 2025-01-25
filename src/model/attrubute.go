package model

import "github.com/google/uuid"

type Attribute struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	EquipmentID uuid.UUID `gorm:"type:uuid;null" json:"equipment_id"`
	Key         string    `gorm:"type:varchar(50);not null"`
	Value       string    `gorm:"type:varchar(255);not null"`
}
