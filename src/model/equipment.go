package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Equipment struct {
	ID               uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name             string             `gorm:"type:text;not null" json:"name"`
	Category         string             `gorm:"type:varchar(50);not null" json:"category"`
	Description      string             `gorm:"type:text;not null" json:"description"`
	Brand            string             `gorm:"type:text" json:"brand"`
	Model            string             `gorm:"type:text" json:"model"`
	Color            string             `gorm:"type:text" json:"color"`
	Material         string             `gorm:"type:varchar(100)" json:"material"`
	MuscleGroups     []MuscleGroup      `gorm:"many2many:equipment_muscle_groups" json:"muscle_groups"`
	EquipmentFeature []EquipmentFeature `gorm:"foreignKey:EquipmentID;OnDelete:CASCADE" json:"equipment_feature"`
	EquipmentOptions []EquipmentOption  `gorm:"foreignKey:EquipmentID;OnDelete:CASCADE" json:"equipment_options"`
	Attribute        []Attribute        `gorm:"foreignKey:EquipmentID;OnDelete:CASCADE" json:"attributes"`
}

func (e *Equipment) BeforeDelete(tx *gorm.DB) error {
	err := tx.Model(&e).Association("MuscleGroups").Clear()
	if err != nil {
		return err
	}

	return nil
}
