package helper

import (
	"errors"
	"gorm.io/gorm"
)

type CustomRecordNotFoundError struct {
	Msg string
}

func (e *CustomRecordNotFoundError) Error() string {
	return e.Msg
}

// Implement `Is` method to make it compatible with `errors.Is(err, gorm.ErrRecordNotFound)`
func (e *CustomRecordNotFoundError) Is(target error) bool {
	return errors.Is(target, gorm.ErrRecordNotFound)
}
