package persistence

import (
	"strings"

	"github.com/google/uuid"
)

func safeUUID(value string) uuid.UUID {
	parsed, err := uuid.Parse(strings.TrimSpace(value))
	if err != nil {
		return uuid.Nil
	}

	return parsed
}
