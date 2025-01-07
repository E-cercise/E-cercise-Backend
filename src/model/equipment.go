package model

import "github.com/google/uuid"

type Equipment struct {
	ID                uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name              string        `gorm:"type:varchar(255);not null" json:"name"`
	Images            []Image       `gorm:"foreignKey:EquipmentID" json:"images"`
	Price             float64       `gorm:"type:decimal(10,2);not null" json:"price"`
	Brand             string        `gorm:"type:varchar(100)" json:"brand"`
	Model             string        `gorm:"type:varchar(100)" json:"model"`
	Color             string        `gorm:"type:varchar(50)" json:"color"`
	Material          string        `gorm:"type:varchar(100)" json:"material"`
	Weight            float64       `gorm:"type:decimal(10,2)" json:"weight"`
	RemainingProducts int           `gorm:"type:int;not null;default:0" json:"remaining_products"`
	SpecialFeature    string        `gorm:"type:text" json:"special_feature"`
	MuscleGroups      []MuscleGroup `gorm:"many2many:equipment_muscle_groups" json:"muscle_groups"`
}
