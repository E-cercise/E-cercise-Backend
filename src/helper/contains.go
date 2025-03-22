package helper

import "github.com/google/uuid"

func Contains(slice []uuid.UUID, value uuid.UUID) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}