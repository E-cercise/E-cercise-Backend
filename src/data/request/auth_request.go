package request

import (
	"github.com/E-cercise/E-cercise/src/enum"
	"github.com/google/uuid"
)

type RegisterRequest struct {
	Email       string              `json:"email" validate:"required,email"`
	Password    string              `json:"password" validate:"required,min=8"`
	FirstName   string              `json:"first_name"`
	LastName    string              `json:"last_name"`
	Address     string              `json:"address"`
	PhoneNumber string              `json:"phone_number"`
	Weight      float64             `json:"weight"`
	Height      float64             `json:"height"`
	Experience  enum.UserExperience `json:"experience"`
	GoalID      uuid.UUID           `json:"goal_id"`
	Preferences []uuid.UUID         `json:"preferences"`
	Gender      enum.Gender         `json:"gender"`
	Age         int                 `json:"age"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
