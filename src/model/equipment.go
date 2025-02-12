package model

import "github.com/google/uuid"

type Equipment struct {
	ID               uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name             string             `gorm:"type:text;not null" json:"name"`
	Description      string             `gorm:"type:text;not null" json:"description"`
	Brand            string             `gorm:"type:text" json:"brand"`
	Model            string             `gorm:"type:text" json:"model"`
	Color            string             `gorm:"type:text" json:"color"`
	Material         string             `gorm:"type:varchar(100)" json:"material"`
	MuscleGroups     []MuscleGroup      `gorm:"many2many:equipment_muscle_groups" json:"muscle_groups"`
	EquipmentFeature []EquipmentFeature `gorm:"foreignKey:EquipmentID" json:"equipment_feature"`
	EquipmentOptions []EquipmentOption  `gorm:"foreignKey:EquipmentID" json:"equipment_options"`
	Attribute        []Attribute        `gorm:"foreignKey:EquipmentID" json:"attributes"`
}
