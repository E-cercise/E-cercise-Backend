package helper

import (
	"github.com/google/uuid"
	"strings"
)

func ParseUUIDs(ids []string) ([]uuid.UUID, error) {
	var uuids []uuid.UUID
	for _, id := range ids {
		parsedID, err := uuid.Parse(strings.TrimSpace(id))
		if err != nil {
			return nil, err
		}
		uuids = append(uuids, parsedID)
	}
	return uuids, nil
}