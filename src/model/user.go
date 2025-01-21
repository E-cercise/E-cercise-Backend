package model

import (
	"github.com/E-cercise/E-cercise/src/enum"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	FirstName   string    `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName    string    `gorm:"type:varchar(100);not null" json:"last_name"`
	Email       string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	Password    string    `gorm:"type:varchar(255);not null" json:"-"`
	Role        enum.Role `gorm:"gorm:type:role_type;unique;not null;default:'USER'" json:"role"`
	Address     string    `gorm:"type:text" json:"address"`
	PhoneNumber string    `gorm:"type:varchar(20)" json:"phone_number"`
	Orders      []Order   `gorm:"foreignKey:UserID" json:"orders"`
	Cart        Cart      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"cart"` // One-to-one relationship
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	cart := Cart{
		UserID: u.ID,
	}
	if err := tx.Create(&cart).Error; err != nil {
		return err
	}
	return nil
}
