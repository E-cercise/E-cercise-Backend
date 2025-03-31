package response

import (
	"github.com/E-cercise/E-cercise/src/enum"
	"github.com/google/uuid"
)

type UserProfileResponse struct {
	Email       string              `json:"email"`
	FirstName   string              `json:"first_name"`
	LastName    string              `json:"last_name"`
	Address     string              `json:"address"`
	PhoneNumber string              `json:"phone_number"`
	Weight      float64             `json:"weight"`
	Height      float64             `json:"height"`
	Experience  enum.UserExperience `json:"experience"`
	Goal        GoalResponse        `json:"goal"`
	Preferences []PrefResponse      `json:"preferences"`
	Gender      enum.Gender         `json:"gender"`
	Age         int                 `json:"age"`
}

type GoalResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type PrefResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
