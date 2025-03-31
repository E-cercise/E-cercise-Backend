package model

import (
	"github.com/google/uuid"
)

type UserPreference struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	TagID  uuid.UUID `gorm:"type:uuid;not null" json:"tag_id"`

	Tag Tag `gorm:"foreignKey:TagID" json:"tag"`
}
