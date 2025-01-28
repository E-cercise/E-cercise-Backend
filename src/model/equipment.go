package model

import "github.com/google/uuid"

type Equipment struct {
	ID               uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name             string             `gorm:"type:varchar(255);not null" json:"name"`
	Description      string             `gorm:"type:varchar(255);not null" json:"description"`
	Images           []Image            `gorm:"foreignKey:EquipmentID" json:"images"`
	Brand            string             `gorm:"type:varchar(100)" json:"brand"`
	Model            string             `gorm:"type:varchar(100)" json:"model"`
	Color            string             `gorm:"type:varchar(50)" json:"color"`
	Material         string             `gorm:"type:varchar(100)" json:"material"`
	MuscleGroups     []MuscleGroup      `gorm:"many2many:equipment_muscle_groups" json:"muscle_groups"`
	EquipmentFeature []EquipmentFeature `gorm:"foreignKey:EquipmentID" json:"equipment_feature"`
	EquipmentOptions []EquipmentOption  `gorm:"foreignKey:EquipmentID" json:"equipment_options"`
	Attribute        []Attribute        `gorm:"foreignKey:EquipmentID" json:"attributes"`
}
