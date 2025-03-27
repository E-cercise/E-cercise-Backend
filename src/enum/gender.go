package enum

import (
	"database/sql/driver"
	"fmt"
)

type Gender string

const (
	GenderMale   Gender = "Male"
	GenderFemale Gender = "Female"
)

func (r *Gender) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid str")
	}

	*r = Gender(str)
	return nil
}

func (r *Gender) Value() (driver.Value, error) {
	return string(*r), nil
}
