package enum

import (
	"database/sql/driver"
	"fmt"
)

type UserExperience string

const (
	Beginner     UserExperience = "Beginner"
	Intermediate UserExperience = "Intermediate"
	Advanced     UserExperience = "Advanced"
	Athlete      UserExperience = "Athlete"
	Elderly      UserExperience = "Elderly"
)

func (r *UserExperience) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid str")
	}
	*r = UserExperience(str)
	return nil
}

func (r *UserExperience) Value() (driver.Value, error) {
	return string(*r), nil
}
