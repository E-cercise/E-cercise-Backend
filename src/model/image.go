package model

import (
	"github.com/E-cercise/E-cercise/src/enum"
	"github.com/google/uuid"
)

type Image struct {
	ID             uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	EquipmentID    uuid.UUID       `gorm:"type:uuid;null" json:"equipment_id"`
	IsPrimary      bool            `gorm:"type:boolean;default:false" json:"is_primary"`
	CloudinaryPath string          `gorm:"type:varchar(255)" json:"cloudinary_path"`
	State          enum.ImageState `gorm:"type:varchar(50);default:'temp'" json:"state"`
}
