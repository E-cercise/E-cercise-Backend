package model

import "github.com/google/uuid"

type Image struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	EquipmentID uuid.UUID `gorm:"type:uuid;not null" json:"equipment_id"`
	IsPrimary   bool      `gorm:"type:boolean;default:false" json:"is_primary"`
	S3Path      string    `gorm:"type:varchar(255)" json:"s3_path"`
}
