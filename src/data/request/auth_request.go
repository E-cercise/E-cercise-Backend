package request

import (
	"github.com/E-cercise/E-cercise/src/enum"
	"github.com/google/uuid"
)

type RegisterRequest struct {
	Email       string               `json:"email" validate:"required,email"`
	Password    string               `json:"password" validate:"required,min=8"`
	FirstName   string               `json:"first_name"`
	LastName    string               `json:"last_name"`
	Address     string               `json:"address"`
	PhoneNumber string               `json:"phone_number"`
	Weight      *float64             `json:"weight,omitempty"`
	Height      *float64             `json:"height,omitempty"`
	Experience  *enum.UserExperience `json:"experience,omitempty"`
	GoalID      *uuid.UUID           `json:"goal_id,omitempty"`
	Preferences []uuid.UUID          `json:"preferences,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
