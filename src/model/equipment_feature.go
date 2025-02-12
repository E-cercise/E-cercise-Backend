package model

import "github.com/google/uuid"

type EquipmentFeature struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	EquipmentID uuid.UUID `gorm:"type:uuid;not null" json:"equipment_id"`
	Description string    `gorm:"type:text" json:"description"`
}
