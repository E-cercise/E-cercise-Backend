package model

type MuscleGroup struct {
	ID         string      `gorm:"type:varchar(50);primaryKey" json:"id"`
	Name       string      `gorm:"type:varchar(100);not null" json:"name"`
	Category   string      `gorm:"type:varchar(50);not null" json:"category"`
	Equipments []Equipment `gorm:"many2many:equipment_muscle_groups" json:"equipments"`
}
