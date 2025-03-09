package model

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LineEquipment struct {
	ID                uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	OrderID           *uuid.UUID `gorm:"type:uuid;"`
	CartID            *uuid.UUID `gorm:"type:uuid" json:"cart_id"`
	EquipmentID       uuid.UUID  `gorm:"type:uuid;not null" json:"equipment_id"`
	EquipmentOptionID uuid.UUID  `gorm:"type:uuid;not null" json:"equipment_option_id"`
	Quantity          int        `gorm:"type:int;not null;default:1" json:"quantity"`
}

func (l *LineEquipment) BeforeUpdate(tx *gorm.DB) (err error) {
	if l.Quantity < 0 {
		return errors.New("quantity cannot be less than 0")
	}

	if l.Quantity == 0 {
		tx.Delete(l)
	}
	return nil
}
