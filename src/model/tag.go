package model

import (
	"github.com/google/uuid"
)

type Tag struct {
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name string    `gorm:"type:varchar(50);unique;not null" json:"name"`
}
