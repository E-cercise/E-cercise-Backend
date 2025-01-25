package model

import "github.com/google/uuid"

type Cart struct {
	ID             uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID         uuid.UUID       `gorm:"type:uuid;not null;unique" json:"user_id"`
	LineEquipments []LineEquipment `gorm:"foreignKey:CartID" json:"line_equipments"`
}
