package request

import (
	"github.com/E-cercise/E-cercise/src/enum"
	"github.com/google/uuid"
)

type UpdateUserProfileRequest struct {
	Address     *string              `json:"address,omitempty"`
	Email       *string              `json:"email,omitempty"`
	FirstName   *string              `json:"first_name,omitempty"`
	LastName    *string              `json:"last_name,omitempty"`
	PhoneNumber *string              `json:"phone_number,omitempty"`
	Weight      *float64             `json:"weight,omitempty"`
	Height      *float64             `json:"height,omitempty"`
	Experience  *enum.UserExperience `json:"experience,omitempty"`
	GoalID      *uuid.UUID           `json:"goal_id,omitempty"`
	Preferences []uuid.UUID          `json:"preferences,omitempty"`
	Gender      *enum.Gender         `json:"gender,omitempty"`
	Age         *int                 `json:"age,omitempty"`
}
