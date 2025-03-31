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
	Role        enum.Role `gorm:"type:role_type;not null;default:'USER'" json:"role"`
	Address     string    `gorm:"type:text;not null" json:"address"`
	PhoneNumber string    `gorm:"type:varchar(20);not null" json:"phone_number"`

	Weight     float64             `gorm:"type:decimal(10,2);not null" json:"weight"`
	Height     float64             `gorm:"type:decimal(10,2);not null" json:"height"`
	Experience enum.UserExperience `gorm:"type:user_experience;not null" json:"experience"`
	Gender     enum.Gender         `gorm:"type:gender_type;not null" json:"gender"`
	Age        int                 `gorm:"type:int;not null" json:"age"`

	GoalID uuid.UUID `gorm:"type:uuid" json:"goal_id,omitempty"`
	Goal   Goal      `gorm:"foreignKey:GoalID" json:"goal,omitempty"`

	UserPreferences []UserPreference `gorm:"foreignKey:UserID" json:"user_preferences"`
	Orders          []Order          `gorm:"foreignKey:UserID" json:"orders"`
	Cart            Cart             `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"cart"`
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
