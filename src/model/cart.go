package model

import "github.com/google/uuid"

type Cart struct {
	ID             uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	LineEquipments []LineEquipment `gorm:"foreignKey:CartID" json:"line_equipments"`
}
