package response

import (
	"github.com/E-cercise/E-cercise/src/enum"
	"github.com/google/uuid"
)

type UserProfileResponse struct {
	Email       string               `json:"email"`
	FirstName   string               `json:"first_name"`
	LastName    string               `json:"last_name"`
	Address     string               `json:"address"`
	PhoneNumber string               `json:"phone_number"`
	Weight      *float64             `json:"weight,omitempty"`
	Height      *float64             `json:"height,omitempty"`
	Experience  *enum.UserExperience `json:"experience,omitempty"`
	Goal        *GoalResponse        `json:"goal,omitempty"`
	Preferences []PrefResponse       `json:"preferences,omitempty"`
}

type GoalResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type PrefResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
