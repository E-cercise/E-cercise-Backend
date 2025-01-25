package enum

import (
	"database/sql/driver"
	"fmt"
)

type Role string

const (
	RoleUser  Role = "USER"
	RoleAdmin Role = "ADMIN"
)

func (r *Role) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid str")
	}

	*r = Role(str)
	return nil
}

func (r *Role) Value() (driver.Value, error) {
	return string(*r), nil
}
